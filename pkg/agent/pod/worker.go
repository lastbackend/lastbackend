package pod

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/cri"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

const ContainerRestartTimeout = 10 //seconds

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

	policy types.PodPolicy
	spec   types.PodSpec

	pod *types.Pod
}

func NewWorker() *Worker {
	return &Worker{
		done: make(chan bool),
	}
}

func NewTask(policy types.PodPolicy, spec types.PodSpec, p *types.Pod) *Task {
	log := context.Get().GetLogger()
	log.Debugf("Create new task for pod: %s", p.ID())
	return &Task{policy: policy, spec: spec, pod: p, done: make(chan bool), close: make(chan bool)}
}

func (w *Worker) Proceed(policy types.PodPolicy, spec types.PodSpec, p *types.Pod) {
	log := context.Get().GetLogger()
	log.Debugf("Proceed new task for pod: %s", p.ID())

	// Clean next task if exists
	if w.next != nil {
		w.next.clean()
		w.next = nil
	}

	t := NewTask(policy, spec, p)

	// Update next task for execution
	if w.current != nil {
		w.lock.Lock()
		w.next = t
		w.lock.Unlock()
		w.current.stop()
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
	log.Debugf("start task for pod: %s", t.pod.ID())

	var err error

	crii := context.Get().GetCri()
	log.Debug(t.pod.Containers)

	if t.policy.Restart {
		for _, c := range t.pod.Containers {
			timeout := time.Duration(ContainerRestartTimeout) * time.Second
			crii.ContainerRestart(c.CID, &timeout)
			t.pod.SetContainer(c)
		}
		t.pod.Policy.Restart = false
		return
	}

	// Remove old containers
	for _, c := range t.pod.Containers {
		if c.CID != "" {
			crii.ContainerRemove(c.CID, true, true)
		}
		t.pod.DelContainer(c.ID)
	}

	// Pull new images
	n_images := make(map[string]types.ImageSpec)
	o_images := make(map[string]types.ImageSpec)

	for _, spec := range t.spec.Containers {
		n_images[spec.Image] = t.spec.Images[spec.Image]
	}

	for _, spec := range t.pod.Spec.Containers {
		o_images[spec.Image] = t.spec.Images[spec.Image]
	}

	for image, spec := range n_images {
		if _, ok := o_images[image]; ok {
			// old image exists
			// check if we need to pull it
			if t.policy.PullImage {
				crii.ImagePull(spec)
				continue
			}

			delete(o_images, image)
			continue
		}

		crii.ImagePull(spec)
	}

	// Clean up unused images
	for _, spec := range o_images {
		crii.ImageRemove(spec.Image())
	}

	for _, spec := range t.spec.Containers {

		c := types.Container{
			ID:      types.ContainerID(uuid.NewV4()),
			State:   types.ContainerStatePending,
			Created: time.Now(),
		}

		c.CID, err = crii.ContainerCreate(spec)

		if err != nil {
			c.State = types.ContainerStateError
			c.Status = err.Error()
			t.pod.AddContainer(c)
			continue
		}

		t.pod.AddContainer(c)
	}

	log.Debugf("done task for pod: %s", t.pod.ID())
}

func (t *Task) stop() {
	t.close <- true
}

func (t *Task) clean() {
	close(t.close)
}
