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

package user

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

func New(obj *types.User) *User {
	u := new(User)

	u.Username = obj.Username
	u.Gravatar = obj.Gravatar
	u.Profile.FirstName = obj.Profile.FirstName
	u.Profile.LastName = obj.Profile.LastName
	u.Updated = obj.Updated
	u.Created = obj.Created
	u.Emails = make(Emails, len(obj.Emails))

	for k, v := range obj.Emails {
		u.Emails[k] = v
	}

	return u
}

func (obj *User) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
