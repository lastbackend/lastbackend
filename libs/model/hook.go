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
	ID string `json:"id"`
	// Project name
	Project string `json:"project"`
	// User username
	User string `json:"user"`
	// Hook token
	Token string `json:"token"`
	// Hook image to build
	Image string `json:"image"`
	// Hook service to build images
	Service string `json:"service"`
	// Hook created time
	Created time.Time `json:"created"`
	// Hook updated time
	Updated time.Time `json:"updated"`
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
