package model

import "time"

type PodList []Pod

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
	ID      string
	// Pod owner
	Owner   string
	// Pod project
	Project string
	// Pod service
	Service string
}

type PodPolicy struct {
	// Pull image flag
	PullImage bool
	// Restart containers flag
	Restart   bool
}

