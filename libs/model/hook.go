package model

import "time"

type HookList []Hook

type Hook struct {
	// Hook uuid, incremented automatically
	ID string `json:"id" gorethink:"id"`
	// Hook owner
	User string `json:"user" gorethink:"user"`
	// Hook token
	Token string `json:"token" gorethink:"token"`
	// Hook image to build
	Image string `json:"image" gorethink:"name"`
	// Hook service to build images
	Service string `json:"service" gorethink:"name"`
	// Hook created time
	Created time.Time `json:"created" gorethink:"created"`
	// Hook updated time
	Updated time.Time `json:"updated" gorethink:"updated"`
}
