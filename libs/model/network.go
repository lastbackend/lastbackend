package model

import "time"

type NetworkList []Image

type Network struct {
	// Image uuid, incremented automatically
	ID string `json:"id"`
	// Image user
	User string `json:"user"`
	// Image name
	Name string `json:"name"`
	// Image created time
	Created time.Time `json:"created"`
	// Image updated time
	Updated time.Time `json:"updated"`
}
