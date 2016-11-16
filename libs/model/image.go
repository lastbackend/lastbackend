package model

import "time"

type ImageList []Image

type Image struct {
	// Image uuid, incremented automatically
	ID string `json:"id, omitempty" gorethink:"id,omitempty"`
	// Image user
	User string `json:"user, omitempty" gorethink:"user,omitempty"`
	// Image name
	Name string `json:"name, omitempty" gorethink:"name,omitempty"`
	// Image tag lists
	Tags map[string]ImageTag `json:"tags, omitempty" gorethink:"tags,omitempty"`
	// Image created time
	Created time.Time `json:"created, omitempty" gorethink:"created,omitempty"`
	// Image updated time
	Updated time.Time `json:"updated, omitempty" gorethink:"updated,omitempty"`
}

type ImageTag struct {
}
