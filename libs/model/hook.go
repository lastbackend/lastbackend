package model

import "time"

type HookList []Hook

type Hook struct {
	// Hook uuid, incremented automatically
	ID string `json:"id, omitempty" gorethink:"id,omitempty"`
	// Hook owner
	User string `json:"user, omitempty" gorethink:"user,omitempty"`
	// Hook token
	Token string `json:"token, omitempty" gorethink:"token,omitempty"`
	// Hook image to build
	Image string `json:"image, omitempty" gorethink:"name,omitempty"`
	// Hook service to build images
	Service string `json:"service, omitempty" gorethink:"name,omitempty"`
	// Hook created time
	Created time.Time `json:"created, omitempty" gorethink:"created,omitempty"`
	// Hook updated time
	Updated time.Time `json:"updated, omitempty" gorethink:"updated,omitempty"`
}
