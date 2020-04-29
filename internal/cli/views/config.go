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

type ConfigList []*Config
type Config views.Config

func (sl *ConfigList) Print() {

	t := table.New([]string{"NAME", "TYPE", "SIZE"})
	t.VisibleHeader = true

	for _, s := range *sl {
		var data = map[string]interface{}{}
		data["NAME"] = s.Meta.Name
		data["TYPE"] = s.Meta.Kind
		t.AddRow(data)
	}
	println()
	t.Print()
	println()
}

func (s *Config) Print() {
	var meta = map[string]interface{}{}
	meta["NAME"] = s.Meta.Name
	meta["TYPE"] = s.Meta.Kind
	println()
	table.PrintHorizontal(meta)
	println()
}

func FromApiConfigView(secret *views.Config) *Config {

	if secret == nil {
		return nil
	}

	item := Config(*secret)
	return &item
}

func FromApiConfigListView(secrets *views.ConfigList) *ConfigList {
	var items = make(ConfigList, 0)
	for _, secret := range *secrets {
		items = append(items, FromApiConfigView(secret))
	}
	return &items
}
