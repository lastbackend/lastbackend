package model

import "time"

type ContainerList []Image

type Container struct {
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
