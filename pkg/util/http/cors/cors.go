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

package cors

import "net/http"

func Headers(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "X-CSRF-Auth, Authorization, Content-Type, x-lastbackend, Origin, X-Requested-With, Content-HostPath, Accept")
	w.Header().Add("Content-Type", "application/json")
}
