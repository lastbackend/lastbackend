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

package haproxy

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/ingress"
	"github.com/lastbackend/lastbackend/pkg/node/runtime"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"text/template"
)

const (
	name         = "ingress-haproxy"
	file         = "haproxy.cfg"
	volume       = "config"
	defaultImage = "index.lastbackend.com/lastbackend/ingress-haproxy"
)

type Haproxy struct {
	ingress.Ingress
	Spec struct {
		Image  struct{
			Name string
			Secret string
		}
		Volume string
		Ports  []uint16
	}
	Status *types.IngressStatus
	Pod *types.PodStatus
	tmpl      *template.Template
}

func (h *Haproxy) Info(ctx context.Context) *types.IngressStatus {
	return nil
}

func (h *Haproxy) Update(ctx context.Context) error {

	routes := envs.Get().GetState().Routes().GetRoutes()

	log.Debugf("Update routes: %d", len(routes))

	var cfg = struct {
		Routes map[string]*types.RouteManifest
	}{
		Routes: routes,
	}

	buf := &bytes.Buffer{}

	h.tmpl.Execute(buf, cfg)
	vol := envs.Get().GetState().Volumes().GetVolume(fmt.Sprintf("%s-%s", name, volume))
	if vol != nil {
		path := filepath.Join(vol.Path, file)
		log.Debugf("config path: %s", path)
		return ioutil.WriteFile(path, buf.Bytes(), 0644)
	}

	log.Debugf("volume not ready: %s:%s", name, volume)
	return errors.New("volume not ready")
}

func (h *Haproxy) Destroy(ctx context.Context) error {
	return nil
}

func (h *Haproxy) Restore(ctx context.Context) error {
	return h.Provision(ctx)
}

func (h *Haproxy) Provision(ctx context.Context) (error) {

	envs.Get().GetState().Pods().SetLocal(name)

	var im = new(types.PodManifest)
	im.Local = true
	container := types.SpecTemplateContainer{
		Name: "ingress",
		Image:types.SpecTemplateContainerImage{
			Name: h.Spec.Image.Name,
			Secret: h.Spec.Image.Secret,
		},
		Labels: make(map[string]string, 0),
		Ports: make([]*types.SpecTemplateContainerPort, 0),
		Volumes: make([]*types.SpecTemplateContainerVolume, 0),
		Security: types.SpecTemplateContainerSecurity{
			Privileged: true,
		},
	}

	for _, p := range h.Spec.Ports {
		container.Ports = append(container.Ports, &types.SpecTemplateContainerPort {
			HostPort: p,
			ContainerPort: p,
		})
	}

	container.Volumes = append(container.Volumes, &types.SpecTemplateContainerVolume{
		Name: volume,
		Path: "/etc/haproxy",
	})

	volume := types.SpecTemplateVolume{
		Name: volume,
	}

	container.Labels[types.ContainerTypeLBC] = name
	im.Template.Containers = append(im.Template.Containers, &container)
	im.Template.Volumes    = append(im.Template.Volumes, &volume)

	if err := runtime.PodManage(ctx, name, im); err != nil {
		log.Errorf("create new pod for ingress err: %s", err.Error())
		return err
	}

	return nil
}

func New() (*Haproxy, error) {

	log.Debug("Enable ingress Haproxy server")
	var proxy = new(Haproxy)

	proxy.tmpl = template.Must(template.New("").Parse(HaproxyTemplate))
	proxy.Spec.Image.Name = defaultImage
	proxy.Spec.Ports = []uint16{80, 443}

	// Set custom ingress image
	if viper.IsSet("ingress.opts.image.name") {
		log.Debugf("ingress image use: %s", viper.GetString("ingress.opts.image.name") )
		proxy.Spec.Image.Name = viper.GetString("ingress.opts.image.name")
	}

	if viper.IsSet("ingress.opts.image.secret") {
		log.Debugf("ingress image secret: %s", viper.GetString("ingress.opts.image.secret") )
		proxy.Spec.Image.Secret = viper.GetString("ingress.opts.image.secret")
	}

	// Set custom ingress ports
	if viper.IsSet("ingress.opts.ports") {
		ports := viper.GetStringSlice("ingress.opts.ports")
		proxy.Spec.Ports = make([]uint16, 0)
		for _, p := range ports {
			var base = 10
			var size = 16
			port, err := strconv.ParseUint(p, base, size)
			if err != nil {
				continue
			}
			proxy.Spec.Ports = append(proxy.Spec.Ports, uint16(port))
		}
	}

	return proxy, nil
}
