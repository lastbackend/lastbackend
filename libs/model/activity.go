package model

import (
	"encoding/json"
	"time"
)

type ActivityList []Activity

type Activity struct {
	// Activity uuid, incremented automatically
	ID      string `json:"id" gorethink:"id,omitempty"`
	// Activity user
	User    string `json:"user" gorethink:"user,omitempty"`
	// Activity project
	Project string `json:"project" gorethink:"project,omitempty"`
	// Activity service
	Service string `json:"service" gorethink:"service,omitempty"`
	// Activity name
	Name    string `json:"name" gorethink:"name,omitempty"`
	// Activity status
	Event   string `json:"event" gorethink:"event,omitempty"`
	// Activity created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Activity updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

func (s *Activity) ToJson() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *ActivityList) ToJson() ([]byte, error) {

	if s == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
