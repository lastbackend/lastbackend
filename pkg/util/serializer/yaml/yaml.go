/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package yaml

import (
	"gopkg.in/yaml.v2"
	"io"
)

type Encoder struct{}
type Decoder struct{}

func (Encoder) Encode(objPtr interface{}, w io.Writer) error {
	buf, err := yaml.Marshal(objPtr)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func (Decoder) Decode(data []byte, objPtr interface{}) error {
	return yaml.Unmarshal(data, objPtr)
}
