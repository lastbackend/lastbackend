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
	//"github.com/satori/go.uuid"
	"time"
	"encoding/json"
)

const EmptyString = ""
const EmptyStringSlice = "[]"

type Meta struct {
	// Meta name
	Name string `json:"name,omitempty"`
	// Meta description
	Description string `json:"description,omitempty"`
	// Meta labels
	Labels map[string]string `json:"labels,omitempty"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

func (m *Meta) SetDefault() {
	m.Labels = make(map[string]string)
	m.Created = time.Now()
	m.Updated = time.Now()
}

type StringSlice []string

func (s *StringSlice) ToJson() string {
	if s == nil {
		return EmptyStringSlice
	}
	res, err := json.Marshal(s)
	if err != nil {
		return EmptyStringSlice
	}
	if string(res) == "null" {
		return EmptyStringSlice
	}
	return string(res)
}