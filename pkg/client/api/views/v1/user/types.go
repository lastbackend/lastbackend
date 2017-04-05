//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package user

import (
	"github.com/lastbackend/lastbackend/pkg/util/table"
	"time"
)

type User struct {
	Gravatar string    `json:"gravatar"`
	Username string    `json:"username"`
	Emails   Emails    `json:"emails"`
	Profile  Profile   `json:"profile"`
	Vendors  Vendors   `json:"integrations"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type Emails map[string]bool

type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Vendors map[string]string
