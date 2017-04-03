package pod

import "github.com/lastbackend/lastbackend/pkg/apis/types"

type Queue struct {
	task  chan *Task
	done  chan bool
	tasks map[string]*Task
}

type Task struct {
	Operation string
	Pod       *types.Pod
}

func New() *Queue {
	q := new(Queue)

	q.task = make(chan *Task)
	q.done = make(chan bool)

	go q.loop(q.done)
	return q
}

func (q *Queue) Add(operation string, pod *types.Pod) {

}

func (q *Queue) Del() {}

func (q *Queue) Get() {}

func (q *Queue) loop(done chan bool) {

}
