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

package http

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"time"
)

const (
	logLevel = 3
)

type JobHttpHook struct {
	timeout time.Time
	config  JobHttpHookConfig
}

type config map[string]interface{}

type JobHttpHookConfig struct {
}

func (hw *JobHttpHook) Execute(task *types.Task) error {
	return nil
}

func New(cfg config) (*JobHttpHook, error) {

	log.V(logLevel).Debug("Use http task watcher")

	var (
		provider *JobHttpHook
	)

	return provider, nil
}
