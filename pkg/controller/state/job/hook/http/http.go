//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"bytes"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"net/http"
	"strings"
	"time"
)

const (
	logLevel = 3
)

type JobHttpHook struct {
	timeout time.Time
	config  *types.JobSpecHookHTTP
}

func (h *JobHttpHook) Execute(task *types.Task) (err error) {

	response, err := v1.View().Task().New(task).ToJson()

	fmt.Println("call hook:", h.config.Endpoint)

	client := http.Client{}
	req, err := http.NewRequest(strings.ToUpper(h.config.Method), h.config.Endpoint, bytes.NewBuffer(response))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if len(h.config.Headers) > 0 {
		for k, v := range h.config.Headers {
			req.Header.Add(k, v)
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.V(logLevel).Debugf("http:job:hook:> response Status: %s", resp.Status)

	return nil
}

func New(cfg *types.JobSpecHookHTTP) (hook *JobHttpHook, err error) {
	log.V(logLevel).Debug("Use http hook")
	hook = new(JobHttpHook)
	hook.config = cfg
	return hook, err
}
