//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package v3

import (
	"context"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/store"
	"regexp"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/types"
)

const (
	incomingBufSize = 100
	outgoingBufSize = 100
)

type watcher struct {
	client *clientv3.Client
}

type watchChan struct {
	watcher   *watcher
	key       string
	filter    string
	rev       int64
	recursive bool
	ctx       context.Context
	cancel    context.CancelFunc
	event     chan *event
	result    chan *store.Event
	error     chan error
}

type event struct {
	key       string
	value     []byte
	prevValue []byte
	rev       int64
	isDeleted bool
	isCreated bool
}

func newWatcher(client *clientv3.Client) *watcher {
	return &watcher{
		client: client,
	}
}

func (w *watcher) Watch(ctx context.Context, key, keyRegexFilter string) (store.Watcher, error) {

	wc := w.newWatchChan(ctx, key, keyRegexFilter)

	go wc.run()

	return wc, nil
}

func (wc *watchChan) Stop() {
	wc.cancel()
}

func (wc *watchChan) ResultChan() <-chan *store.Event {
	return wc.result
}

func (w *watcher) newWatchChan(ctx context.Context, key, keyRegexFilter string) *watchChan {
	wc := &watchChan{
		watcher: w,
		key:     key,
		filter:  keyRegexFilter,
		event:   make(chan *event, incomingBufSize),
		result:  make(chan *store.Event, outgoingBufSize),
		error:   make(chan error, 1),
	}
	wc.ctx, wc.cancel = context.WithCancel(ctx)
	return wc
}

func (wc *watchChan) run() {
	watchClosedCh := make(chan struct{})
	go wc.watching(watchClosedCh)

	var resultChanWG sync.WaitGroup
	resultChanWG.Add(1)
	go wc.handleEvent(&resultChanWG)

	select {
	case err := <-wc.error:
		if err == context.Canceled {
			break
		}
		errResult := transformError(err)
		if errResult != nil {
			// guarantee of error after closing
			select {
			case wc.result <- errResult:
			case <-wc.ctx.Done():
			}
		}
	case <-watchClosedCh:
	case <-wc.ctx.Done():
	}

	wc.cancel()

	// wait until the result is used
	resultChanWG.Wait()
	close(wc.result)
}

func (wc *watchChan) watching(watchClosedCh chan struct{}) {

	if wc.rev == 0 {
		if err := wc.getState(); err != nil {
			log.Errorf("%s:watching:> failed to getState with latest state: %v", logPrefix, err)
			wc.sendError(err)
			return
		}
	}

	opts := []clientv3.OpOption{
		clientv3.WithRev(wc.rev + 1),
		clientv3.WithPrevKV(),
		clientv3.WithPrefix(),
	}

	r, _ := regexp.Compile(wc.filter)

	wch := wc.watcher.client.Watch(wc.ctx, wc.key, opts...)
	for wres := range wch {
		if wres.Err() != nil {
			err := wres.Err()
			log.Errorf("%s:watching:> watch chan err: %v", logPrefix, err)
			wc.sendError(err)
			return
		}

		for _, we := range wres.Events {
			if r.MatchString(string(we.Kv.Key)) {
				e := &event{
					key:       string(we.Kv.Key),
					value:     we.Kv.Value,
					rev:       we.Kv.ModRevision,
					isDeleted: we.Type == clientv3.EventTypeDelete,
					isCreated: we.IsCreate(),
				}
				if we.PrevKv != nil {
					e.prevValue = we.PrevKv.Value
				}
				wc.sendEvent(e)
			}
		}
	}

	close(watchClosedCh)
}

func (wc *watchChan) getState() error {

	opts := []clientv3.OpOption{clientv3.WithPrefix()}

	getResp, err := wc.watcher.client.Get(wc.ctx, wc.key, opts...)
	if err != nil {
		return err
	}

	wc.rev = getResp.Header.Revision

	for _, kv := range getResp.Kvs {
		wc.sendEvent(&event{
			key:       string(kv.Key),
			value:     kv.Value,
			prevValue: nil,
			rev:       kv.ModRevision,
			isDeleted: false,
			isCreated: true,
		})
	}

	return nil
}

func (wc *watchChan) handleEvent(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case e := <-wc.event:
			res := transformEvent(e)
			if res == nil {
				continue
			}
			if len(wc.result) == outgoingBufSize {
				log.Warnf("%s:handleevent:> buffered events: %d. Processing takes a long time", logPrefix, outgoingBufSize)
			}
			select {
			case wc.result <- res:
			case <-wc.ctx.Done():
				return
			}
		case <-wc.ctx.Done():
			return
		}
	}
}

func (wc *watchChan) sendError(err error) {
	select {
	case wc.error <- err:
	case <-wc.ctx.Done():
	}
}

func (wc *watchChan) sendEvent(e *event) {
	if len(wc.event) == incomingBufSize {
		log.Warnf("%s:sendevent:> buffered events: %d. Processing takes a long time", logPrefix, incomingBufSize)
	}
	select {
	case wc.event <- e:
	case <-wc.ctx.Done():
	}
}

func transformEvent(e *event) *store.Event {

	action := store.STORAGEUPDATEEVENT

	if e.isCreated {
		action = store.STORAGECREATEEVENT
	}

	if e.isDeleted {
		action = store.STORAGEDELETEEVENT
	}

	event := &store.Event{
		Type:   action,
		Key:    string(e.key),
		Object: e.value,
	}

	return event
}

func transformError(err error) *store.Event {
	return &store.Event{
		Type:   store.STORAGEERROREVENT,
		Object: err,
	}
}
