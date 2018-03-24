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

package view

import (
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

type SecretList []*Secret
type Secret struct {
	Meta SecretMeta `json:"meta"`
}

type SecretMeta struct {
	Name string `json:"name"`
}

func (sl *SecretList) Print() {

	t := table.New([]string{"NAME"})
	t.VisibleHeader = true

	for _, s := range *sl {
		var data = map[string]interface{}{}
		data["NAME"] = s.Meta.Name
		t.AddRow(data)
	}
	println()
	t.Print()
	println()
}

func (s *Secret) Print() {
	var data = map[string]interface{}{}
	data["NAME"] = s.Meta.Name
	println()
	table.PrintHorizontal(data)
	println()
}

func FromApiSecretView(secret *views.Secret) *Secret {
	var item = new(Secret)
	item.Meta.Name = secret.Meta.Name
	return item
}

func FromApiSecretListView(secrets *views.SecretList) *SecretList {
	var items = make(SecretList, 0)
	for _, secret := range *secrets {
		items = append(items, FromApiSecretView(secret))
	}
	return &items
}
