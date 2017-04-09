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
	"sync"
	"time"
)

type ImageList []Image

type Image struct {
	lock sync.RWMutex

	// Image id
	ID string `json:"id"`
	// Image name
	Name string `json:"name"`
	// Image tag lists
	Tags []string `json:"tags"`
	// Image created time
	Created time.Time `json:"created"`
	// Image updated time
	Updated time.Time `json:"updated"`
}

type ImageSpec struct {
	// Image full name
	Name string `json:"name"`
	// Image pull provision flag
	Pull bool `json:"image-pull"`
	// Image Auth base64 encoded string
	Auth string `json:"auth"`
}

func NewImage() *Image {
	return &Image{}
}
