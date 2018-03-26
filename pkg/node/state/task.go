//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package state

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

func (s *TaskState) AddTask(key string, task *types.NodeTask) {
	log.V(logLevel).Debugf("Cache: PodCache: add cancel func pod: %#v", key)
	s.tasks[key] = *task
}

func (s *TaskState) GetTask(key string) *types.NodeTask {
	log.V(logLevel).Debugf("Cache: PodCache: get cancel func pod: %#v", key)

	if _, ok := s.tasks[key]; ok {
		t := s.tasks[key]
		return &t
	}

	return nil
}

func (s *TaskState) DelTask(pod *types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: del cancel func pod: %#v", pod)
	delete(s.tasks, pod.Meta.Name)
}
