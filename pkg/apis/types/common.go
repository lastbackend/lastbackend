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

package types

import (
	"time"
)

type Meta struct {
	// Service unique name
	Name string `json:"name"`
	// Service description
	Description string `json:"description"`
	// Service labels
	Labels map[string]string `json:"lables"`
	// Service created time
	Created time.Time `json:"created"`
	// Service updated time
	Updated time.Time `json:"updated"`
}
