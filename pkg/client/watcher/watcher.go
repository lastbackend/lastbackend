//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package watcher

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/tools/log"
	"io"
	"sync"
)

type IWatcher interface {
	Stop()
	ResultChan() <-chan Event
}

// EventType defines the possible types of events.
type EventType string

const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Error    EventType = "ERROR"
)

type Event struct {
	Type EventType
	Data interface{}
}

type Watcher struct {
	sync.Mutex
	reader  io.ReadCloser
	result  chan Event
	stopped bool
}

func NewStreamWatcher(reqder io.ReadCloser) *Watcher {
	sw := &Watcher{
		reader: reqder,
		result: make(chan Event),
	}
	go sw.receive()
	return sw
}

func (w *Watcher) ResultChan() <-chan Event {
	return w.result
}

func (w *Watcher) Stop() {
	w.Lock()
	defer w.Unlock()
	if !w.stopped {
		w.stopped = true
		w.reader.Close()
	}
}

func (w *Watcher) stopping() bool {
	w.Lock()
	defer w.Unlock()
	return w.stopped
}

func (w *Watcher) receive() {
	defer close(w.result)
	defer w.Stop()
	for {
		result := new(Event)
		err := json.NewDecoder(w.reader).Decode(result)
		if err != nil {
			log.Error(err)
			if w.stopping() {
				return
			}
			return
		}

		w.result <- Event{
			Type: result.Type,
			Data: result.Data,
		}
	}
}
