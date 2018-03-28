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

import "time"

type Config struct {
	// Server requires Bearer authentication.
	BearerToken string
	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout time.Duration
	// Disable security check for a client
	Insecure bool
}

func (c *Config) SetDefault() {
	c.Timeout = 10
}
