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

import "time"

type PodList []*Pod

type Pod struct {
	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Pod provision flag
	Policy PodPolicy `json:"provision"`
	// Container spec
	Spec []ContainerSpec `json:"spec"`
	// Containers status info
	Containers []ContainerStatusInfo `json:"containers"`
	// Container created time
	Created time.Time `json:"created"`
	// Container updated time
	Updated time.Time `json:"updated"`
}

type PodMeta struct {
	// Pod ID
	ID string
	// Pod owner
	Owner string
	// Pod project
	Project string
	// Pod service
	Service string
}

type PodPolicy struct {
	// Pull image flag
	PullImage bool
	// Restart containers flag
	Restart bool
}
