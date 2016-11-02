package model

import "time"

type ImageList []Image

type Image struct {
	// Project uuid, incremented automatically
	ID string `json:"id, omitempty" gorethink:"id,omitempty"`
	// Project user
	User string `json:"user, omitempty" gorethink:"user,omitempty"`
	// Project name
	Name string `json:"name, omitempty" gorethink:"name,omitempty"`

	Tags map[string]ImageTag `json:"tags, omitempty" gorethink:"tags,omitempty"`
	// Project created time
	Created time.Time `json:"created, omitempty" gorethink:"created,omitempty"`
	// Project updated time
	Updated time.Time `json:"updated, omitempty" gorethink:"updated,omitempty"`
}

type ImageTag struct {
}
