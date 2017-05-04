//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package events

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/system"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"time"
)

func NewTickerEvent() {
	var event = new(types.Event)
	event.Ticker = true
	event.Meta = system.GetNodeMeta()
	event.State = system.GetNodeState()
	event.Pods = make([]*types.Pod, 0)
	event.Timestamp = time.Now()

	context.Get().GetEventListener().Send(event)
	return
}

func NewInitialEvent(pods []*types.Pod) {
	var event = new(types.Event)
	event.Initial = true
	event.Meta = system.GetNodeMeta()
	event.State = system.GetNodeState()
	event.Pods = pods
	event.Timestamp = time.Now()

	context.Get().GetEventListener().Send(event)
	return
}

func NewEvent(pods []*types.Pod) {
	var event = new(types.Event)
	event.Meta = system.GetNodeMeta()
	event.State = system.GetNodeState()
	event.Pods = pods
	event.Timestamp = time.Now()
	context.Get().GetEventListener().Send(event)
	return
}

func SendPodState(pod *types.Pod) {
	p := types.Pod{
		Meta:       pod.Meta,
		State:      pod.State,
		Containers: pod.Containers,
	}
	NewEvent(append([]*types.Pod{}, &p))
}


