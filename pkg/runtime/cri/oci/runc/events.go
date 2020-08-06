// +build linux
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

package runc

import (
	"fmt"
	"sync"
	"time"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/types"
)

const (
	EventStats = "stats"
	EventOOM   = "oom"
)

type EventsOptions struct {
	// Set the stats collection interval
	Interval time.Duration
	// Display the container's stats then exit
	Stats bool
}

// EventsContainer - displays information about the container. By default the information is displayed once every 5 seconds.
// Where 'containerID' is your name for the instance of the container.
// EventsOptions:
//   Interval <duration>: set the stats collection interval
//   Stats <bool>: display the container's stats then exit
func (r *runc) EventsContainer(containerID string, opts EventsOptions) (*EventWatcher, error) {
	container, err := r.getContainer(containerID)
	if err != nil {
		return nil, err
	}

	if opts.Interval <= 0 {
		return nil, fmt.Errorf("duration interval must be greater than 0")
	}

	status, err := container.Status()
	if err != nil {
		return nil, err
	}

	if status == libcontainer.Stopped {
		return nil, fmt.Errorf("container with id %s is not running", container.ID())
	}

	eventWatcher := newEventWatcher()

	go func() {

		var (
			stats  = make(chan *libcontainer.Stats, 1)
			events = make(chan *types.Event, 1024)
			group  = &sync.WaitGroup{}
		)

		group.Add(1)

		go func() {
			defer group.Done()
			for e := range events {
				eventWatcher.send(e)
			}
		}()

		if opts.Stats {
			s, err := container.Stats()
			if err != nil {
				return
			}
			events <- &types.Event{Type: EventStats, ID: container.ID(), Data: convertLibcontainerStats(s)}
			close(events)
			group.Wait()
			return
		}

		go func() {
			for range time.Tick(opts.Interval) {
				s, err := container.Stats()
				if err != nil {
					fmt.Println(err)
					continue
				}
				stats <- s
			}
		}()

		n, err := container.NotifyOOM()
		if err != nil {
			return
		}

		for {
			select {
			case _, ok := <-n:
				if ok {
					// this means an oom event was received, if it is !ok then
					// the channel was closed because the container stopped and
					// the cgroups no longer exist.
					events <- &types.Event{Type: EventOOM, ID: container.ID()}
				} else {
					n = nil
				}
			case s := <-stats:
				events <- &types.Event{Type: EventStats, ID: container.ID(), Data: convertLibcontainerStats(s)}
			}
			if n == nil {
				close(events)
				break
			}
		}

		group.Wait()
	}()

	return eventWatcher, nil
}

type Event *types.Event

type EventWatcher struct {
	eventsChannel chan Event
}

func newEventWatcher() *EventWatcher {
	ew := new(EventWatcher)
	return ew
}

func (ew *EventWatcher) send(event *types.Event) {
	ew.eventsChannel <- event
}

func (ew *EventWatcher) Events() <-chan Event {
	return ew.eventsChannel
}
