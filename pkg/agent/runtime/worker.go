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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime/cri"
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

	meta  types.PodMeta
	spec  types.PodSpec

	pod *types.Pod
}

func NewWorker() *Worker {
	return &Worker{
		done: make(chan bool),
	}
}

func NewTask(meta types.PodMeta, spec types.PodSpec, pod *types.Pod) *Task {
	log := context.Get().GetLogger()
	log.Debugf("Create new task for pod: %s", pod.Meta.ID)
	return &Task{
		meta:  meta,
		spec:  spec,
		pod:   pod,
		done:  make(chan bool),
		close: make(chan bool),
	}
}

func (w *Worker) Proceed(meta types.PodMeta, spec types.PodSpec, p *types.Pod) {
	log := context.Get().GetLogger()
	log.Debugf("Proceed new task for pod: %s", p.Meta.ID)

	// Clean next task if exists
	if w.next != nil {
		w.next.clean()
		w.next = nil
	}

	t := NewTask(meta, spec, p)

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
		w.current = nil

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

	// Check spec version
	log.Debugf("pod spec: %s, new spec: %s", t.pod.Spec.ID, t.spec.ID)
	if t.spec.ID != t.pod.Spec.ID {
		log.Debugf("spec is differrent, apply new one: %s", t.pod.Spec.ID)
		// Set current spec
		t.pod.Spec.ID = t.spec.ID
		t.imagesUpdate()
		t.containersUpdate()
	}

	// check container state
	t.containersState()
	log.Debugf("done task for pod: %s", t.pod.Meta.ID)
}

func (t *Task) imagesUpdate() {
	log := context.Get().GetLogger()
	crii := context.Get().GetCri()

	// Check images states
	images := make(map[string]struct{})

	// Get images currently used by this pod
	for _, container := range t.pod.Containers {
		log.Debugf("Add images as used: %s", container.Image)
		images[container.Image] = struct{}{}
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
		log.Debugf("Delete unused images: %s", name)
		crii.ImageRemove(name)
	}

}

func (t *Task) containersUpdate() {

	log := context.Get().GetLogger()
	crii := context.Get().GetCri()

	log.Debugf("Start containers update process for pod: %s", t.pod.Meta.ID)
	var err error

	var ids []string
	var ncs []*types.Container

	// Remove old containers
	for _, c := range t.pod.Containers {
		if c.ID != "" {
			ids = append(ids, c.ID)
		}
	}

	// Create new containers
	for _, spec := range t.spec.Containers {
		log.Debugf("Container create")

		c := &types.Container{
			Image:   spec.Image.Name,
			State:   types.ContainerStatePending,
			Created: time.Now(),
		}

		if spec.Labels == nil {
			spec.Labels = make(map[string]string)
		}

		spec.Labels["LB_META"] = fmt.Sprintf("%s/%s", t.pod.Meta.ID, t.pod.Spec.ID)
		c.ID, err = crii.ContainerCreate(spec)

		if err != nil {
			log.Errorf("Container create error %s", err.Error())
			c.State = types.ContainerStateError
			c.Status = err.Error()
			break
		}

		log.Debugf("New container created: %s", c.ID)
		ncs = append(ncs, c)
	}

	for _, c := range ncs {
		t.pod.AddContainer(c)
	}

	for _, id := range ids {

		log.Debugf("Container %s remove", id)
		err := crii.ContainerRemove(id, true, true)
		if err != nil {
			log.Errorf("Container remove error: %s", err.Error())
		}
		t.pod.DelContainer(id)
	}

}

func (t *Task) containersState() {
	// TODO: wait 5 seconds and recheck container state
	log := context.Get().GetLogger()
	log.Debugf("update container state from: %s to %s", t.pod.Meta.State.State, t.meta.State.State)

	t.pod.Meta.State.State = "provision"

	crii := context.Get().GetCri()
	// Update containers states
	if t.meta.State.State == types.PodStateStarted || t.meta.State.State == types.PodStateRunning {
		for _, c := range t.pod.Containers {
			log.Debugf("Container: %s try to start", c.ID)
			err := crii.ContainerStart(c.ID)
			c.State = "running"
			c.Status = ""
			if err != nil {
				log.Errorf("Container: start error: %s", err.Error())
				c.State = "error"
				c.Status = err.Error()
			}
			log.Debugf("Container: %s started", c.ID)
			t.pod.SetContainer(c)
			log.Debugf("Container: %s updated", c.ID)
		}
		t.pod.UpdateState()
		return
	}

	if t.meta.State.State == types.PodStateStopped {
		for _, c := range t.pod.Containers {
			timeout := time.Duration(ContainerStopTimeout) * time.Second
			log.Debugf("Container: %s try to stop", c.ID)
			err := crii.ContainerStop(c.ID, &timeout)
			c.State = "stopped"
			c.Status = ""
			if err != nil {
				log.Errorf("Container: stop error: %s", err.Error())
				c.State = "error"
				c.Status = err.Error()
			}
			log.Debugf("Container: %s stopped", c.ID)
			t.pod.SetContainer(c)
			log.Debugf("Container: %s updated", c.ID)
		}
		t.pod.UpdateState()
		return
	}

	if t.meta.State.State == types.PodStateRestarted {
		for _, c := range t.pod.Containers {
			timeout := time.Duration(ContainerRestartTimeout) * time.Second
			log.Debugf("Container: %s try to restart", c.ID)
			err := crii.ContainerRestart(c.ID, &timeout)
			c.State = "running"
			c.Status = ""
			if err != nil {
				log.Errorf("Container: restart error: %s", err.Error())
				c.State = "error"
				c.Status = err.Error()
			}
			log.Debugf("Container: %s restarted", c.ID)
			t.pod.SetContainer(c)
			log.Debugf("Container: %s updated", c.ID)
		}
		t.pod.UpdateState()
		return
	}
}

func (t *Task) finish() {
	t.close <- true
}

func (t *Task) clean() {
	close(t.close)
}
