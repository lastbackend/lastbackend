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

type VolumeList []Volume

type Volume struct {
	// Volume uuid, incremented automatically
	ID string `json:"id"`
	// Volume uuid, incremented automatically
	Project string `json:"project"`
	// Volume user
	User string `json:"user"`
	// Volume name
	Name string `json:"name"`
	// Volume tag lists
	Created time.Time `json:"created"`
	// Volume updated time
	Updated time.Time `json:"updated"`
}

type VolumesSpec struct {
	Volume    string
	MountPath string
}
