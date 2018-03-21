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

package storage

import (
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/util/filesystem"
	"github.com/lastbackend/lastbackend/pkg/util/homedir"
)

var path = fmt.Sprintf("%s/.lastbackend/token", homedir.HomeDir())

func SetToken(token string) error {
	return filesystem.WriteStrToFile(path, token, 0644)
}

func GetToken() (string, error) {

	buf, err := filesystem.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
