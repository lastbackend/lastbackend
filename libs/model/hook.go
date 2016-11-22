package model

import "time"

type HookList []Hook

type Hook struct {
	// Hook uuid, incremented automatically
	ID string `json:"id" gorethink:"id,omitempty"`
	// Hook owner
	User string `json:"user" gorethink:"user,omitempty"`
	// Hook token
	Token string `json:"token" gorethink:"token,omitempty"`
	// Hook image to build
	Image string `json:"image" gorethink:"name,omitempty"`
	// Hook service to build images
	Service string `json:"service" gorethink:"name,omitempty"`
	// Hook created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Hook updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}
