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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/types"
)

type watcher struct {}

func newWatcher() *watcher {
	return &watcher{}
}

func (w *watcher) Watch(ctx context.Context, key, keyRegexFilter string) (types.Watcher, error) {
	return w, nil
}

func (wc *watcher) Stop() {
	return
}

func (wc *watcher) ResultChan() <-chan *types.Event {
	return nil
}
