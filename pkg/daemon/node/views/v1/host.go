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

package v1

type Meta struct {
	Hostname     string `json:"hostname"`
	OSName       string `json:"os_name"`
	OSType       string `json:"os_type"`
	Architecture string `json:"architecture"`

	CRI     CRIMeta     `json:"cri"`
	CPU     HostCPU     `json:"cpu"`
	Memory  HostMemory  `json:"memory"`
	Network HostNetwork `json:"network"`
	Storage HostStorage `json:"storage"`
}

type CRIMeta struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type HostCPU struct {
	Name  string `json:"name"`
	Cores int64  `json:"cores"`
}

type HostMemory struct {
	Total     int64 `json:"total"`
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
}

type HostNetwork struct {
	Interface string   `json:"interface,omitempty"`
	IP        []string `json:"ip,omitempty"`
}

type HostStorage struct {
	Available string `json:"available"`
	Used      string `json:"used"`
	Total     string `json:"total"`
}
