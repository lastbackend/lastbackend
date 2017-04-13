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

package service

import (
	"fmt"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	n "github.com/lastbackend/lastbackend/pkg/daemon/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

type createS struct {
	Namespace string  `json:"namespace,omitempty"`
	Name      string  `json:"name,omitempty"`
	Template  string  `json:"template,omitempty"`
	Image     string  `json:"image,omitempty"`
	Url       string  `json:"url,omitempty"`
	Config    *Config `json:"config,omitempty"`
}

type Config struct {
	Scale int `json:"scale,omitempty"`
	//Ports   []string `json:"ports,omitempty"`
	//EnvVars     []string `json:"env,omitempty"`
	//Volumes []string `json:"volumes,omitempty"`
}

func CreateCmd(name, image, template, url string, scale int) {

	var (
		config *Config
	)

	if scale != 0 /* || len(env) != 0 || len(ports) != 0 || len(volumes) != 0 */ {
		config = new(Config)
		config.Scale = scale
		//config.EnvVars = env
		//config.Ports = ports
		//config.Volumes = volumes
	}

	err := Create(name, image, template, url, config)
	if err != nil {
		fmt.Print(err)
		return
	}

	// TODO: Waiting for start service
	// TODO: Show spinner
}

func Create(name, image, template, url string, config *Config) error {

	var (
		err       error
		http      = c.Get().GetHttpClient()
		storage   = c.Get().GetStorage()
		namespace = new(n.Namespace)
		er        = new(errors.Http)
		res       = new(struct{})
	)

	err = storage.Get("namespace", namespace)
	if err != nil {
		return errors.New(err.Error())
	}

	if namespace.Name == "" {
		return errors.New("Namespace didn't select")
	}

	var cfg = createS{}
	cfg.Namespace = namespace.Name

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
		POST("/namespace/"+namespace.Name+"/service").
		AddHeader("Content-Type", "application/json").
		BodyJSON(cfg).
		Request(res, er)
	if err != nil {
		return err
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	return nil
}
