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
	"github.com/lastbackend/lastbackend/pkg/agent/events"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/satori/go.uuid"
	"time"
)

const ContainerRestartTimeout = 10 // seconds
const ContainerStopTimeout = 10    // seconds

type Task struct {
	id string

	close chan bool
	done  chan bool

	meta  types.PodMeta
	state types.PodState
	spec  types.PodSpec

	pod *types.Pod
}

func (t *Task) exec() {

	var (
		log = context.Get().GetLogger()
		pod = context.Get().GetStorage().Pods()
	)

	t.pod.Spec.State = types.StateProvision
	events.SendPodState(t.pod)

	defer func() {
		t.pod.Spec.State = types.StateReady
		pod.SetPod(t.pod)
		events.SendPodState(t.pod)
		log.Debugf("Task [%s]: done task for pod: %s", t.id, t.pod.Meta.ID)
	}()

	log.Debugf("Task [%s]: start task for pod: %s", t.id, t.pod.Meta.ID)

	// Check spec version
	log.Debugf("Task [%s]: pod spec: %s, new spec: %s", t.id, t.pod.Spec.ID, t.spec.ID)

	if t.state.State == types.StateDestroy {
		log.Debugf("Task [%s]: pod is marked for deletion: %s", t.id, t.pod.Meta.ID)
		t.containersStateManage()
		return
	}

	if t.spec.ID != t.pod.Spec.ID {
		log.Debugf("Task [%s]: spec is differrent, apply new one: %s", t.id, t.pod.Spec.ID)
		t.pod.Spec.ID = t.spec.ID
		t.imagesUpdate()
		t.containersCreate()
	}

	if len(t.pod.Containers) != len(t.spec.Containers) {
		log.Debugf("Task [%s]: containers count mismatch: %d (%d)", t.id, len(t.pod.Containers), len(t.spec.Containers))
		t.containersCreate()
	}

	// check container state
	t.containersStateManage()
	t.containersRemove()
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
	}

	// Clean up unused images
	for name := range images {
		log.Debugf("Delete unused images: %s", name)
		crii.ImageRemove(name)
	}

}

func (t *Task) containersCreate() {

	var (
		err error
		log = context.Get().GetLogger()
		crii = context.Get().GetCri()
	)

	log.Debugf("Start containers creation process for pod: %s", t.pod.Meta.ID)

	// Create new containers
	for id, spec := range t.spec.Containers {
		log.Debugf("Container create")

		c := types.Container{
			Pod:     t.pod.Meta.ID,
			Spec:    id,
			Image:   spec.Image.Name,
			State:   types.ContainerStatePending,
			Created: time.Now(),
		}

		if spec.Labels == nil {
			spec.Labels = make(map[string]string)
		}

		spec.Labels["LB_META"] = fmt.Sprintf("%s/%s/s", t.pod.Meta.ID, t.pod.Spec.ID, spec.Meta.ID)
		c.ID, err = crii.ContainerCreate(spec)

		if err != nil {
			log.Errorf("Container create error %s", err.Error())
			c.State = types.ContainerStateError
			c.Status = err.Error()
			break
		}

		log.Debugf("New container created: %s", c.ID)
		t.pod.AddContainer(&c)
	}
}

func (t *Task) containersRemove() {

	var (
		log = context.Get().GetLogger()
		crii = context.Get().GetCri()
		specs = make(map[string]bool)
	)

	log.Debugf("Start containers removable process for pod: %s", t.pod.Meta.ID)

	for id := range t.spec.Containers {
		log.Debugf("Add spec %s to valid", id)
		specs[id] = false
	}

	// Remove old containers
	for _, c := range t.pod.Containers {
		log.Debugf("Container %s has spec %s", c.ID, c.Spec)
		if _, ok := specs[c.Spec]; !ok || specs[c.Spec] == true {
			log.Debugf("Container %s needs to be removed", c.ID)
			err := crii.ContainerRemove(c.ID, true, true)
			if err != nil {
				log.Errorf("Container remove error: %s", err.Error())
			}

			t.pod.DelContainer(c.ID)
			continue
		}

		specs[c.Spec] = true
	}

}

func (t *Task) containersStateManage() {

	log := context.Get().GetLogger()
	crii := context.Get().GetCri()

	defer t.pod.UpdateState()

	log.Debugf("update container state from: %s to %s", t.pod.State.State, t.state.State)

	if t.state.State == types.StateDestroy {
		log.Debugf("Pod %s delete %d containers", t.pod.Meta.ID, len(t.pod.Containers))

		for _, c := range t.pod.Containers {

			err := crii.ContainerRemove(c.ID, true, true)
			c.State = types.StateDestroyed
			c.Status = ""

			if err != nil {
				c.State = types.StateError
				c.Status = err.Error()
			}

			t.pod.DelContainer(c.ID)
		}

		return
	}

	// Update containers states
	if t.state.State == types.StateStart || t.state.State == types.StateRunning {
		for _, c := range t.pod.Containers {
			log.Debugf("Container: %s try to start", c.ID)
			err := crii.ContainerStart(c.ID)
			c.State = types.StateRunning
			c.Status = ""

			if err != nil {
				c.State = types.StateError
				c.Status = err.Error()
			}

			t.pod.SetContainer(c)
		}
		return
	}

	if t.state.State == types.StateStop {
		timeout := time.Duration(ContainerStopTimeout) * time.Second

		for _, c := range t.pod.Containers {

			err := crii.ContainerStop(c.ID, &timeout)

			c.State = types.StateStopped
			c.Status = ""

			if err != nil {
				log.Errorf("Container: stop error: %s", err.Error())
				c.State = "error"
				c.Status = err.Error()
			}
			t.pod.SetContainer(c)
		}

		return
	}

	if t.state.State == types.StateRestart {
		timeout := time.Duration(ContainerRestartTimeout) * time.Second

		for _, c := range t.pod.Containers {

			err := crii.ContainerRestart(c.ID, &timeout)
			c.State = types.StateRunning
			c.Status = ""

			if err != nil {
				c.State = "error"
				c.Status = err.Error()
			}

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

func NewTask(meta types.PodMeta, state types.PodState, spec types.PodSpec, pod *types.Pod) *Task {
	log := context.Get().GetLogger()
	uuid := uuid.NewV4().String()
	log.Debugf("Task [%s]: Create new task for pod: %s", uuid, pod.Meta.ID)
	log.Debugf("Task [%s]: Container spec count: %d", uuid, len(spec.Containers))

	return &Task{
		id:    uuid,
		meta:  meta,
		state: state,
		spec:  spec,
		pod:   pod,
		done:  make(chan bool),
		close: make(chan bool),
	}
}
