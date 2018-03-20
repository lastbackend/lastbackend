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

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

type Secret struct {
	Meta  SecretMeta  `json:"meta" yaml:"meta"`
	State SecretState `json:"meta" yaml:"state"`
	Spec  SecretSpec  `json:"meta" yaml:"spec"`
}

type SecretMeta struct {
	Meta             `yaml:",inline"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type SecretSpec struct {
}

type SecretState struct {
}

func (s *Secret) GetHash() string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s:%s", s.Meta.Namespace, s.Meta.Name)))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
