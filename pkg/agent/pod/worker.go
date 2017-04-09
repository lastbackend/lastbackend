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

package pod

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/cri"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"sync"
	"time"
)

const ContainerRestartTimeout = 10 // seconds
const ContainerStopTimeout = 10    // seconds

type Worker struct {
	lock sync.RWMutex

	cri     cri.CRI
	current *Task
	next    *Task

	done chan bool
}

type Task struct {
	close chan bool
	done  chan bool

	state types.PodState
	meta  types.PodMeta
	spec  types.PodSpec

	pod *types.Pod
}

func NewWorker() *Worker {
	return &Worker{
		done: make(chan bool),
	}
}

func NewTask(state types.PodState, meta types.PodMeta, spec types.PodSpec, pod *types.Pod) *Task {
	log := context.Get().GetLogger()
	log.Debugf("Create new task for pod: %s", pod.ID())
	return &Task{
		state: state,
		meta:  meta,
		spec:  spec,
		pod:   pod,
		done:  make(chan bool),
		close: make(chan bool),
	}
}

func (w *Worker) Proceed(state types.PodState, meta types.PodMeta, spec types.PodSpec, p *types.Pod) {
	log := context.Get().GetLogger()
	log.Debugf("Proceed new task for pod: %s", p.ID())

	// Clean next task if exists
	if w.next != nil {
		w.next.clean()
		w.next = nil
	}

	t := NewTask(state, meta, spec, p)

	// Update next task for execution
	if w.current != nil {
		w.lock.Lock()
		w.next = t
		w.lock.Unlock()
		w.current.finish()
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
	for {
		if w.current == nil {
			w.done <- true
			return
		}

		w.current.exec()
		w.lock.Lock()
		if w.next != nil {
			w.current = w.next
			w.next = nil
		}
		w.lock.Unlock()
	}
}

func (t *Task) exec() {

	log := context.Get().GetLogger()
	log.Debugf("start task for pod: %s", t.pod.Meta.ID)

	// Set current spec
	t.pod.Meta.Spec = t.spec.ID
	// Check spec version
	if t.meta.Spec != t.pod.Meta.Spec {
		t.imagesUpdate()
		t.containersUpdate()
	}

	// check container state
	t.containersState()
	log.Debugf("done task for pod: %s", t.pod.ID())
}

func (t *Task) imagesUpdate() {
	log := context.Get().GetLogger()
	crii := context.Get().GetCri()

	// Check images states
	images := make(map[string]struct{})

	// Get images currently used by this pod
	for _, container := range t.pod.Containers {
		log.Debugf("Add images as used: %s", container.Image.Name)
		images[container.Image.Name] = struct{}{}
	}

	// Check imaged we need to pull
	for _, spec := range t.spec.Containers {

		// Check image exists and not need to be pulled
		if _, ok := images[spec.Image.Name]; ok {

			log.Debugf("Image exists in prev spec: %s", spec.Image.Name)
			// Check if image need to be updated
			if !spec.Image.Pull {
				log.Debugf("Image not needed to pull: %s", spec.Image.Name)
				delete(images, spec.Image.Name)
				continue
			}

			log.Debugf("Delete images from unused: %s", spec.Image.Name)
			delete(images, spec.Image.Name)
		}

		log.Debugf("Image update needed: %s", spec.Image.Name)
		crii.ImagePull(&spec.Image)
		// add image to storage
	}

	// Clean up unused images
	for name := range images {
		log.Debug("Delete unused images: %s", name)
		crii.ImageRemove(name)
	}

}

func (t *Task) containersUpdate() {

	log := context.Get().GetLogger()
	crii := context.Get().GetCri()

	log.Debugf("Start containers update process for pod: %s", t.pod.Meta.ID)
	var err error

	// Remove old containers
	for _, c := range t.pod.Containers {
		if c.ID != "" {
			crii.ContainerRemove(c.ID, true, true)
		}
		t.pod.DelContainer(c.ID)
	}

	// Create new containers
	for _, spec := range t.spec.Containers {

		c := &types.Container{
			State:   types.ContainerStatePending,
			Created: time.Now(),
		}

		c.ID, err = crii.ContainerCreate(spec)

		if err != nil {
			c.State = types.ContainerStateError
			c.Status = err.Error()
			t.pod.AddContainer(c)
			continue
		}

		t.pod.AddContainer(c)
	}

}

func (t *Task) containersState() {

	crii := context.Get().GetCri()
	// Update containers states
	if t.pod.State.State == types.PodStateStarted {
		for _, c := range t.pod.Containers {
			crii.ContainerStart(c.ID)
			t.pod.SetContainer(c)
		}
		return
	}

	if t.pod.State.State == types.PodStateStopped {
		for _, c := range t.pod.Containers {
			timeout := time.Duration(ContainerStopTimeout) * time.Second
			crii.ContainerStop(c.ID, &timeout)
			t.pod.SetContainer(c)
		}
		return
	}

	if t.pod.State.State == types.PodStateRestarted {
		for _, c := range t.pod.Containers {
			timeout := time.Duration(ContainerRestartTimeout) * time.Second
			crii.ContainerRestart(c.ID, &timeout)
			t.pod.SetContainer(c)
		}
		return
	}
}

func (t *Task) finish() {
	t.close <- true
}

func (t *Task) clean() {
	close(t.close)
}
