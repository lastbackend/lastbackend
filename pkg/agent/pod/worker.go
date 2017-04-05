package pod

import "github.com/lastbackend/lastbackend/pkg/apis/types"

type Worker struct {
	spec []types.ContainerSpec
	pod  *types.Pod
}

func (w *Worker) Create(spec []types.ContainerSpec, wait chan bool) {

	if wait != nil {
		<-wait
	}

}
