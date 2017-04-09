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

package types

import (
	"encoding/json"
	"time"
)

type EventList []Event

type Event struct {
	// Activity uuid, incremented automatically
	ID string `json:"id"`
	// Activity project
	Project string `json:"project"`
	// Activity service
	Service string `json:"service"`
	// Activity name
	Name string `json:"name"`
	// Activity status
	Event string `json:"event"`
	// Activity created time
	Created time.Time `json:"created"`
	// Activity updated time
	Updated time.Time `json:"updated"`
}

func (s *Event) ToJson() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *EventList) ToJson() ([]byte, error) {

	if s == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
