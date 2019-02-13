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

package docker

import (
	"encoding/json"
	"fmt"
	"time"
)

type jsonTime struct {
	time.Time
}

type jsonLogLine struct {
	Message          string            `json:"message"`
	ContainerId      string            `json:"container_id"`
	ContainerName    string            `json:"container_name"`
	ContainerCreated jsonTime          `json:"container_created"`
	ImageId          string            `json:"image_id"`
	ImageName        string            `json:"image_name"`
	Command          string            `json:"command"`
	Tag              string            `json:"tag"`
	Extra            map[string]string `json:"extra"`
	Host             string            `json:"host"`
	Timestamp        jsonTime          `json:"timestamp"`
}

func logMessage(lp *logPair, message []byte) error {
	lp.logLine.Message = string(message[:])
	lp.logLine.Timestamp = jsonTime{time.Now()}

	_, err := json.Marshal(lp.logLine)
	if err != nil {
		return err
	}

	return nil
}

func (t jsonTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", t.Format(time.RFC3339Nano))
	return []byte(str), nil
}
