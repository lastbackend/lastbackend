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

package builder

import (
	"github.com/lastbackend/lastbackend/pkg/agent/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/builder/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"sync"
)

type Worker struct {
	lock sync.RWMutex

	id    string
	image string
	queue []*types.Build
	cri   cri.CRI

	done chan bool
}

func (w *Worker) NewBuild(build *types.Build) {
	log := context.Get().GetLogger()
	log.Debugf("Create new build: %s", build.Meta.ID)
	// Add new build to build queue

	// Run goroutine with current task
	go w.loop()
}

func (w *Worker) provision() error {

	var (
		err error
		log = context.Get().GetLogger()
	)

	spec := types.ContainerSpec{
		Image: types.ImageSpec{
			Name: "docker:dind",
		},
	}

	w.id, err = w.cri.ContainerCreate(context.Get(), &spec)
	if err != nil {
		log.Warnf("Can not create container for docker builds:%s", err.Error())
		return err
	}

	return nil
}

func (w *Worker) destroy() error {
	// remove docker daemon
	return w.cri.ContainerRemove(context.Get(), w.id, true, true)
}

func (w *Worker) loop() {

}

func NewWorker(cri cri.CRI) *Worker {

	var (
		w = new(Worker)
	)

	w.done = make(chan bool)
	w.cri = cri

	return w
}
