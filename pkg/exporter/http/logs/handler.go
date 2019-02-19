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

package logs

import (
	"github.com/lastbackend/lastbackend/pkg/exporter/envs"
	"github.com/lastbackend/lastbackend/pkg/log"
	"net/http"
)

const (
	logPrefix   = "exporter:http:logs"
	logLevel    = 3
	BUFFER_SIZE = 512
)

func LogsH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/service/{service}/logs service serviceLogs
	//
	// Shows logs of the service
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: service
	//     in: path
	//     description: service id
	//     required: true
	//     type: string
	//   - name: deployment
	//     in: query
	//     description: deployment id
	//     required: true
	//     type: string
	//   - name: pod
	//     in: query
	//     description: pod id
	//     required: true
	//     type: string
	//   - name: container
	//     in: query
	//     description: container id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Service logs received
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	kind := r.URL.Query().Get("kind")
	selflink := r.URL.Query().Get("selflink")

	log.V(logLevel).Debugf("%s:logs:> get by selflink `%s`", logPrefix, selflink)

	var (
		logger = envs.Get().GetLogger()
	)

	if logger == nil {
		return
	}

	if err := logger.Stream(r.Context(), kind, selflink, w); err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	<-r.Context().Done()
}
