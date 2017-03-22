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

package model

import "time"

type ImageList []Image

type Image struct {
	// Image uuid, incremented automatically
	ID string `json:"id" gorethink:"id,omitempty"`
	// Image user
	User string `json:"user" gorethink:"user,omitempty"`
	// Image name
	Name string `json:"name" gorethink:"name,omitempty"`
	// Image tag lists
	Tags map[string]ImageTag `json:"tags" gorethink:"tags,omitempty"`
	// Image created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Image updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

type ImageTag struct {
}
