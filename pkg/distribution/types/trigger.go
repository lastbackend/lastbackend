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

import "fmt"

type TriggerList []Trigger

type Trigger struct {
	Meta   TriggerMeta   `json:"meta"`
	Spec   TriggerSpec   `json:"spec"`
	Status TriggerStatus `json:"status"`
}

type TriggerMeta struct {
	Meta
	Namespace string `json:"namespace"`
	Service   string `json:"service"`
}

type TriggerStatus struct {
	Stage   string `json:"stage"`
	Message string `json:"message"`
}

type TriggerSpec struct {
}

func (t *Trigger) SelfLink() string {
	if t.Meta.SelfLink == "" {
		t.Meta.SelfLink = fmt.Sprintf("%s:%s:%s", t.Meta.Namespace, t.Meta.Service, t.Meta.Name)
	}
	return t.Meta.SelfLink
}
