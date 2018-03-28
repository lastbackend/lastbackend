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

package trigger

import (
	"github.com/lastbackend/lastbackend/pkg/log"
	"net/http"
)

const (
	logLevel = 2
	logPrefix = "api:handler:trigger"
)

func HookExecuteH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:execute execute hook", logPrefix)

	var (
		err error
	)

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.Errorf("%s:execute write response err: %s", logPrefix, err.Error())
		return
	}
}
