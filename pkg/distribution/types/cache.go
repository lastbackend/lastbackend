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

package types

const STORAGEPUTEVENT = "PUT"
const STORAGEDELEVENT = "DELETE"

type NodeStatusEvent struct {
	Event  string `json:"event"`
	Node   string `json:"node"`
	Online bool   `json:"online"`
}

type NetworkSpecEvent struct {
	Event string      `json:"event"`
	Node  string      `json:"node"`
	Spec  NetworkSpec `json:"spec"`
}

type PodSpecEvent struct {
	Event string  `json:"event"`
	Node  string  `json:"node"`
	Name  string  `json:"name"`
	Spec  PodSpec `json:"spec"`
}

type RouteSpecEvent struct {
	Event string    `json:"event"`
	Name  string    `json:"name"`
	Spec  RouteSpec `json:"spec"`
}

type VolumeSpecEvent struct {
	Event string     `json:"event"`
	Node  string     `json:"node"`
	Name  string     `json:"name"`
	Spec  VolumeSpec `json:"spec"`
}

type IngressStatusEvent struct {
	Event  string        `json:"event"`
	Name   string        `json:"name"`
	Status IngressStatus `json:"status"`
}
