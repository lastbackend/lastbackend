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

package service

import (
	"fmt"
	a "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

type createS struct {
	App      string  `json:"app,omitempty"`
	Name     string  `json:"name,omitempty"`
	Template string  `json:"template,omitempty"`
	Image    string  `json:"image,omitempty"`
	Url      string  `json:"url,omitempty"`
	Config   *Config `json:"config,omitempty"`
}

type Config struct {
	Replicas int `json:"replicas,omitempty"`
	//Ports   []string `json:"ports,omitempty"`
	//EnvVars     []string `json:"env,omitempty"`
	//Volumes []string `json:"volumes,omitempty"`
}

func CreateCmd(name, image, template, url string, replicas int) {

	var (
		config *Config
	)

	if replicas != 0 /* || len(env) != 0 || len(ports) != 0 || len(volumes) != 0 */ {
		config = new(Config)
		config.Replicas = replicas
		//config.EnvVars = env
		//config.Ports = ports
		//config.Volumes = volumes
	}

	err := Create(name, image, template, url, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Waiting for start service
	// TODO: Show spinner

	fmt.Println("Service `" + name + "` is succesfully created")
}

func Create(name, image, template, url string, config *Config) error {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		storage = c.Get().GetStorage()
		app     = new(a.App)
		er      = new(errors.Http)
		res     = new(struct{})
	)

	app, err = storage.App().Load()
	if err != nil {
		return err
	}

	if app.Meta.Name == "" {
		return errors.New("App didn't select")
	}

	var cfg = createS{}
	cfg.App = app.Meta.Name

	if name != "" {
		cfg.Name = name
	}

	if template != "" {
		cfg.Template = template
	}

	if image != "" {
		cfg.Image = image
	}

	if url != "" {
		cfg.Url = url
	}

	if config != nil {
		cfg.Config = config
	}

	_, _, err = http.
		POST(fmt.Sprintf("/app/%s/service", app.Meta.Name)).
		AddHeader("Content-Type", "application/json").
		BodyJSON(cfg).
		Request(res, er)
	if err != nil {
		return errors.New(er.Message)
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	return nil
}
