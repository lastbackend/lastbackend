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
	"net/http"
	"time"
)

const (
	logLevel = 3
)

type JobHttpProvider struct {
	timeout time.Time
	config  JobHttpProviderConfig
	client  http.Client
}

type config map[string]interface{}

type JobHttpProviderConfig struct {
}

func (hw *JobHttpProvider) Fetch() (*types.Task, error) {
	return nil, nil
}

func New(cfg config) (*JobHttpProvider, error) {

	log.V(logLevel).Debug("Use http task watcher")

	var (
		provider *JobHttpProvider
	)

	return provider, nil
}
