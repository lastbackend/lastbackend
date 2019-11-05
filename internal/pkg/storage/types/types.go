//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package types

import "github.com/lastbackend/lastbackend/internal/pkg/types"

const (
	STORAGEDELETEEVENT = types.EventActionDelete
	STORAGECREATEEVENT = types.EventActionCreate
	STORAGEUPDATEEVENT = types.EventActionUpdate
	STORAGEERROREVENT  = types.EventActionError
)

type Watcher interface {
	Stop()
	ResultChan() <-chan *Event
}

type WatcherEvent struct {
	System
	Action   string
	Name     string
	SelfLink string
	Data     interface{}
}

type Event struct {
	Type   string
	Key    string
	Rev    int64
	Object interface{}
}

type Kind string

func (k Kind) String() string {
	return string(k)
}

type QueryFilter string

func (qf QueryFilter) String() string {
	return string(qf)
}

type Opts struct {
	Ttl   uint64
	Force bool
	Rev   *int64
}

type System struct {
	types.System
}
