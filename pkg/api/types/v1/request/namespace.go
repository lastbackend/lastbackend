//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package request

// swagger:model request_namespace_create
type NamespaceCreateOptions struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Domain      *string                 `json:"domain"`
	Quotas      *NamespaceQuotasOptions `json:"quotas"`
}

// swagger:model request_namespace_update
type NamespaceUpdateOptions struct {
	Description *string                 `json:"description"`
	Domain      *string                 `json:"domain"`
	Quotas      *NamespaceQuotasOptions `json:"quotas"`
}

// swagger:model request_namespace_remove
type NamespaceRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:model request_namespace_quotas
type NamespaceQuotasOptions struct {
	Disabled bool  `json:"disabled"`
	RAM      int64 `json:"ram"`
	Routes   int   `json:"routes"`
}
