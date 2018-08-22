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
	Data map[string]int
}

type SecretMeta struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func (sl *SecretList) Print() {

	t := table.New([]string{"NAME", "TYPE", "SIZE"})
	t.VisibleHeader = true

	for _, s := range *sl {
		var data = map[string]interface{}{}
		data["NAME"] = s.Meta.Name
		data["TYPE"] = s.Meta.Kind
		size := 0
		for _, d := range s.Data {
			size+=d
		}
		data["SIZE"] = size
		t.AddRow(data)
	}
	println()
	t.Print()
	println()
}

func (s *Secret) Print() {
	var meta = map[string]interface{}{}
	meta["NAME"] = s.Meta.Name
	meta["TYPE"] = s.Meta.Kind
	println()
	table.PrintHorizontal(meta)
	println()
	var data = map[string]interface{}{}
	for n, d := range s.Data {
		data[n] = d
	}
	table.PrintHorizontal(data)

}

func FromApiSecretView(secret *views.Secret) *Secret {
	var item = new(Secret)

	item.Meta.Name = secret.Meta.Name
	item.Meta.Kind = secret.Meta.Kind
	item.Data = make(map[string]int, 0)

	for n,d := range secret.Data {
		item.Data[n] = len(d)
	}

	return item
}

func FromApiSecretListView(secrets *views.SecretList) *SecretList {
	var items = make(SecretList, 0)
	for _, secret := range *secrets {
		items = append(items, FromApiSecretView(secret))
	}
	return &items
}
