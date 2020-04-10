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

package views

import (
	"github.com/lastbackend/lastbackend/internal/util/table"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type SecretList []*Secret
type Secret views.Secret

func (sl *SecretList) Print() {

	t := table.New([]string{"NAME", "TYPE", "SIZE"})
	t.VisibleHeader = true

	for _, s := range *sl {
		var data = map[string]interface{}{}
		data["NAME"] = s.Meta.Name
		data["TYPE"] = s.Spec.Type
		size := 0
		for _, d := range s.Spec.Data {
			size += len(d)
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
	meta["TYPE"] = s.Spec.Type
	println()
	table.PrintHorizontal(meta)
	println()
	var data = map[string]interface{}{}
	for n, d := range s.Spec.Data {
		data[n] = d
	}
	table.PrintHorizontal(data)

}

func FromApiSecretView(secret *views.Secret) *Secret {

	if secret == nil {
		return nil
	}

	item := Secret(*secret)
	return &item
}

func FromApiSecretListView(secrets *views.SecretList) *SecretList {
	var items = make(SecretList, 0)
	for _, secret := range *secrets {
		items = append(items, FromApiSecretView(secret))
	}
	return &items
}
