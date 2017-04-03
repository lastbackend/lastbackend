package node

import "time"

type PodList struct {
}

type Pod struct {
	Gravatar string    `json:"gravatar"`
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type Container struct {
}
