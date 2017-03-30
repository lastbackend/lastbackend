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

import "time"

type ClusterList []Cluster

type Cluster struct {
	// Cluster uuid, generated automatically
	ID string `json:"id"`
	// Cluster owner username
	Owner string `json:"owner"`
	// Cluster name
	Name string `json:"name"`
	// Cluster region
	Region string `json:"name"`
	// Cluster labels lists
	Labels map[string]string `json:"labels"`
	// Cluster created time
	Created time.Time `json:"created"`
	// Cluster updated time
	Updated time.Time `json:"updated"`
}
