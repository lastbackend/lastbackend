//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

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
