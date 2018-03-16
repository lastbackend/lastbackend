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

func (s *TaskState) AddTask(pod *types.Pod, task *types.NodeTask) {
	log.V(logLevel).Debugf("Cache: PodCache: add cancel func pod: %#v", pod)
	s.tasks[pod.Meta.Name] = *task
}

func (s *TaskState) GetTask(pod *types.Pod) *types.NodeTask {
	log.V(logLevel).Debugf("Cache: PodCache: get cancel func pod: %#v", pod)

	if _, ok := s.tasks[pod.Meta.Name]; ok {
		t := s.tasks[pod.Meta.Name]
		return &t
	}

	return nil
}

func (s *TaskState) DelTask(pod *types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: del cancel func pod: %#v", pod)
	delete(s.tasks, pod.Meta.Name)
}
