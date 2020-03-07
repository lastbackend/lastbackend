//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package v3

const (
	logLevel  = 6
	logPrefix = "storage:etcd:v3"
)

type Config struct {
	Endpoints []string `json:"endpoints",yaml:"endpoints"`
	TLS       struct {
		Key  string `json:"key",yaml:"key"`
		Cert string `json:"cert",yaml:"cert"`
		CA   string `json:"ca",yaml:"ca"`
	} `json:"tls",yaml:"tls"`
	Prefix string `json:"prefix",yaml:"prefix"`
}
