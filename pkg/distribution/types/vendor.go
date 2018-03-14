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

package types

import (
	"time"

	"golang.org/x/oauth2"
)

const (
	VENDOR_AUTH     = "auth"
	VENDOR_PLATFORM = "platform"
)

type Vendor struct {
	ServiceID string        `json:"service_id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	Host      string        `json:"host"`
	Token     *oauth2.Token `json:"token"`
	Created   time.Time     `json:"created"`
	Updated   time.Time     `json:"updated"`
}
