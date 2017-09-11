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

package request

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/api/node/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"io"
	"io/ioutil"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 3

type RequestNodeEventS struct {
	v1.Event
}

func (s *RequestNodeEventS) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Node: decode and validate event data")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Node: decode and validate event data err: %s", err.Error())
		return errors.New("node event").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Node: convert struct from json err: %s", err.Error())
		return errors.New("node").IncorrectJSON(err)
	}

	return nil
}
