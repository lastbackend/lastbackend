//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"time"
)

const EmptyString = ""
const EmptyStringSlice = "[]"

// swagger:ignore
// swagger:model types_meta
type Meta struct {
	// Meta name
	Name string `json:"name,omitempty",yaml:"name,omitempty"`
	// Meta description
	Description string `json:"description",yaml:"description"`
	// Meta self link
	SelfLink string `json:"self_link",yaml:"self_link"`
	// Meta labels
	Labels map[string]string `json:"labels",yaml:"labels"`
	// Meta created time
	Created time.Time `json:"created",yaml:"created"`
	// Meta updated time
	Updated time.Time `json:"updated",yaml:"updated"`
}

// swagger:ignore
// swagger:model types_meta_create
type MetaCreateOptions struct {
	Name        string            `json:"name",yaml:"name"`
	Description string            `json:"description",yaml:"description"`
	Labels      map[string]string `json:"labels",yaml:"labels"`
}

// swagger:ignore
// swagger:model types_meta_update
type MetaUpdateOptions struct {
	Description *string           `json:"description",yaml:"description"`
	Labels      map[string]string `json:"labels",yaml:"labels"`
}

func (m *Meta) SetDefault() {
	m.Labels = make(map[string]string, 0)
	m.Created = time.Now().UTC()
	m.Updated = time.Now().UTC()
}
