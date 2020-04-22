//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/tools/log"
	"net/http"
	"strings"
	"time"
)

const (
	logLevel = 3
)

type JobHttpProvider struct {
	timeout time.Time
	config  *models.JobSpecProviderHTTP
	client  http.Client
}

func (h *JobHttpProvider) Fetch() (*models.TaskManifest, error) {

	var (
		err      error
		manifest = new(request.TaskManifest)
	)

	client := http.Client{}

	req, err := http.NewRequest(strings.ToUpper(h.config.Method), h.config.Endpoint, nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if len(h.config.Headers) > 0 {
		for k, v := range h.config.Headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if err := manifest.DecodeAndValidate(resp.Body); err != nil {
		log.Error(err.Err().Error())
		return nil, err.Err()
	}

	defer resp.Body.Close()

	mf := new(models.TaskManifest)
	manifest.SetTaskManifestMeta(mf)
	if err := manifest.SetTaskManifestSpec(mf); err != nil {
		return nil, err
	}

	return mf, nil
}

func New(cfg *models.JobSpecProviderHTTP) (*JobHttpProvider, error) {

	log.Debug("Use http task watcher")

	var (
		provider *JobHttpProvider
	)

	provider = new(JobHttpProvider)
	provider.config = cfg

	return provider, nil
}
