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

package json

import (
	"encoding/json"
	"io"
)

type Encoder struct{}
type Decoder struct{}

func (Encoder) Encode(objPtr interface{}, w io.Writer) error {
	buf, err := json.Marshal(objPtr)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func (Decoder) Decode(data []byte, objPtr interface{}) error {
	return json.Unmarshal(data, objPtr)
}
