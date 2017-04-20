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

package runtime

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"sync"
)

type Worker struct {
	lock sync.RWMutex

	pod     string
	cri     cri.CRI
	current *Task
	next    *Task

	done chan bool
}

func (w *Worker) Proceed(meta types.PodMeta, state types.PodState, spec types.PodSpec, p *types.Pod) {
	log := context.Get().GetLogger()
	log.Debugf("Proceed new task for pod: %s", p.Meta.ID)

	// Clean next task if exists

	t := NewTask(meta, state, spec, p)

	if w.next != nil {
		w.lock.Lock()
		log.Debugf("Worker: Replace tasks: %s > %s", w.next.id, t.id)
		w.next.clean()
		w.next = nil
		w.lock.Unlock()
	}

	// Update next task for execution
	if w.current != nil {
		w.lock.Lock()
		w.next = t
		w.lock.Unlock()
		return
	}

	// Create current task
	w.lock.Lock()
	w.current = t
	w.lock.Unlock()

	// Run goroutine with current task
	go w.loop()
}

func (w *Worker) loop() {
	log := context.Get().GetLogger()
	for {
		if w.current == nil {
			w.done <- true
			return
		}

		log.Debugf("Worker: Task [%s]: prepared for execution", w.current.id)
		w.current.exec()
		log.Debugf("Worker: Task [%s]: executed", w.current.id)
		w.current = nil

		w.lock.Lock()
		if w.next != nil {
			w.current = w.next
			w.next = nil
		}
		w.lock.Unlock()
	}
}

func NewWorker() *Worker {
	return &Worker{
		done: make(chan bool),
	}
}
