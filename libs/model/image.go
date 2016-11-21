package model

import "time"

type ImageList []Image

type Image struct {
	// Image uuid, incremented automatically
	ID string `json:"id" gorethink:"id"`
	// Image user
	User string `json:"user" gorethink:"user"`
	// Image name
	Name string `json:"name" gorethink:"name"`
	// Image tag lists
	Tags map[string]ImageTag `json:"tags" gorethink:"tags"`
	// Image created time
	Created time.Time `json:"created" gorethink:"created"`
	// Image updated time
	Updated time.Time `json:"updated" gorethink:"updated"`
}

type ImageTag struct {
}
