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

type HookList []Hook

type Hook struct {
	// Hook uuid, incremented automatically
	ID string `json:"id" gorethink:"id,omitempty"`
	// Hook owner
	User string `json:"user" gorethink:"user,omitempty"`
	// Hook token
	Token string `json:"token" gorethink:"token,omitempty"`
	// Hook image to build
	Image string `json:"image" gorethink:"image,omitempty"`
	// Hook service to build images
	Service string `json:"service" gorethink:"service,omitempty"`
	// Hook created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Hook updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

func (h *Hook) ToJson() ([]byte, error) {
	buf, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (h *HookList) ToJson() ([]byte, error) {

	if h == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
