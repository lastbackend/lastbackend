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

package router

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
	"github.com/lastbackend/lastbackend/pkg/util/nginx"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func Restore(ctx context.Context) error {

	log.Debug("Runtime restore state")

	dir := viper.GetString("node.volume")

	if !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir += string(os.PathSeparator)
	}

	dir += "routes"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		fmt.Println(info)
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".conf" {
			name := strings.Split(info.Name(), ".")[0]
			match := strings.Split(name, ":")
			envs.Get().GetState().Router().Set(match[0], match[1])
		}
		return nil
	})
	if err != nil {
		log.Errorf("Route: Restore: restore configs info err: %s", err)
		return err
	}

	return nil
}

func Create(ctx context.Context, config *types.RouterConfig) error {

	dir := viper.GetString("node.volume")

	if !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir += string(os.PathSeparator)
	}

	dir += "routes"

	filename := fmt.Sprintf("%s:%s.conf", config.ID, config.Hash)

	config.RootPath = viper.GetString("node.volume")

	err := nginx.Nginx{}.GenerateConfig(strings.Join([]string{dir, filename}, string(os.PathSeparator)), config)
	if err != nil {
		return err
	}

	envs.Get().GetState().Router().Set(config.ID, config.Hash)

	return nil
}

func Manage(ctx context.Context, route *types.Route) error {

	log.Debugf("Route: Manage: manage route %s", route.Meta.Name)

	config := route.GetRouteConfig()

	if err := Destroy(ctx, config.ID); err != nil {
		log.Errorf("Route: Manage: remove route %s config err: %s", config.ID, err)
		return err
	}

	if config.State.Destroy {
		events.NewRouteStateEvent(ctx, route)
		return nil
	}

	if err := Create(ctx, config); err != nil {
		log.Errorf("Route: Manage: create route %s config err: %s", config.ID, err)
		return err
	}

	events.NewRouteStateEvent(ctx, route)

	return nil
}

func Destroy(ctx context.Context, id string) error {

	dir := viper.GetString("node.volume")

	if !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir += string(os.PathSeparator)
	}

	dir += "routes"

	filename := fmt.Sprintf("%s:%s.conf", id, envs.Get().GetState().Router().Get(id))
	err := nginx.Nginx{}.RemoveConfig(strings.Join([]string{dir, filename}, string(os.PathSeparator)))
	if err != nil {
		return err
	}

	envs.Get().GetState().Router().Del(id)

	return nil
}
