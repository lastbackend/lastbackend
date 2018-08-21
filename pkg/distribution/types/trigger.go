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

type TriggerMap struct {
	Runtime
	Items map[string]*Trigger
}

type TriggerList struct {
	Runtime
	Items []*Trigger
}


type Trigger struct {
	Runtime
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
	State   string `json:"state"`
	Message string `json:"message"`
}

type TriggerSpec struct {
}

func (t *Trigger) SelfLink() string {
	if t.Meta.SelfLink == "" {
		t.Meta.SelfLink = t.CreateSelfLink(t.Meta.Namespace, t.Meta.Service, t.Meta.Name)
	}
	return t.Meta.SelfLink
}

func (t *Trigger) CreateSelfLink(namespace, service, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, service, name)
}

func NewTriggerList () *TriggerList {
	dm := new(TriggerList)
	dm.Items = make([]*Trigger, 0)
	return dm
}

func NewTriggerMap () *TriggerMap {
	dm := new(TriggerMap)
	dm.Items = make(map[string]*Trigger)
	return dm
}