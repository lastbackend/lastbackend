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
	"github.com/lastbackend/lastbackend/pkg/builder/context"
	"github.com/lastbackend/lastbackend/pkg/common/config"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/satori/go.uuid"
	"github.com/lastbackend/lastbackend/pkg/log"
)

type Task struct {
	id string

	close chan bool
	done  chan bool

	build  types.Build
	daemon config.Docker
}

func (t *Task) start() {

	// send build start event

	defer func() {
		// send build end event
		log.Debugf("Task [%s]: done task for build: %s", t.id, t.build.Meta.Name)
	}()

	log.Debugf("Task [%s]: start task for build: %s", t.id, t.build.Meta.Name)

}

func (t *Task) finish() {
	t.close <- true
}

func (t *Task) clean() {
	close(t.close)
}

func NewTask(build types.Build) *Task {
	uuid := uuid.NewV4().String()
	log.Debugf("Task [%s]: Container spec count: %d", uuid)

	return &Task{
		id:    uuid,
		build: build,
		done:  make(chan bool),
		close: make(chan bool),
	}
}
