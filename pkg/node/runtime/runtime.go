//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package runtime

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/util/decoder"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	logNodeRuntimePrefix = "node:runtime"
	logLevel             = 3
)

// System main node process
type Runtime struct {
	spec chan *types.NodeManifest
}

// Restore node runtime state
func (r *Runtime) Restore(ctx context.Context) error {
	log.V(logLevel).Debugf("%s:restore:> restore init", logNodeRuntimePrefix)

	var network = envs.Get().GetNet()

	if network != nil {
		if err := network.SubnetRestore(ctx); err != nil {
			log.Errorf("can not restore subnets: %s", err.Error())
			return err
		}

		if err := network.EndpointRestore(ctx); err != nil {
			log.Errorf("can not restore endpoint: %s", err.Error())
			return err
		}

		if err := network.ResolverManage(ctx); err != nil {
			log.Errorf("%s:> can not manage resolver:%s", logNodeRuntimePrefix, err.Error())
		}
	}

	if err := VolumeRestore(ctx); err != nil {
		log.Errorf("can not restore volumes: %s", err.Error())
		return err
	}

	if err := ImageRestore(ctx); err != nil {
		log.Errorf("Can not restore images: %s", err.Error())
		return err
	}

	if err := PodRestore(ctx); err != nil {
		log.Errorf("Can not restore pods: %s", err.Error())
		return err
	}

	return nil

}

// Provision node manifest
func (r *Runtime) Provision(ctx context.Context, dir string) error {

	log.V(logLevel).Debugf("%s:provision:> local init", logNodeRuntimePrefix)

	log.V(logLevel).Debugf("%s:provision:> read manifests from dir: %s", logNodeRuntimePrefix, dir)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var (
		mf = new(types.NodeManifest)
	)
	mf.Configs = make(map[string]*types.ConfigManifest)
	mf.Pods = make(map[string]*types.PodManifest)
	mf.Volumes = make(map[string]*types.VolumeManifest)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		c, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			_ = fmt.Errorf("failed read data from file: %s", f)
			continue
		}

		items := decoder.YamlSplit(c)
		log.Debugf("manifests: %d", len(items))

		for _, i := range items {

			var m = new(request.Runtime)

			if err := yaml.Unmarshal([]byte(i), m); err != nil {
				log.Errorf("can not parse manifest: %s: %s", f.Name(), err.Error())
				continue
			}

			switch strings.ToLower(m.Kind) {
			case types.KindConfig:
				m := new(request.ConfigManifest)
				err := m.FromYaml(i)
				if err != nil {
					log.Errorf("invalid specification: %s", err.Error())
					return err
				}
				if m.Meta.Name == nil {
					break
				}
				log.Debugf("Add config Manifest: %s", *m.Meta.Name)
				mf.Configs[*m.Meta.Name] = m.GetManifest()
				break
			case types.KindPod:

				m := new(request.PodManifest)
				err := m.FromYaml(i)
				if err != nil {
					log.Errorf("invalid specification: %s", err.Error())
					return err
				}
				if m.Meta.Name == nil {
					break
				}
				log.Debugf("Add Pod Manifest: %s", *m.Meta.Name)
				mf.Pods[*m.Meta.Name] = m.GetManifest()
				envs.Get().GetState().Pods().SetLocal(*m.Meta.Name)
				break
			case types.KindVolume:

				m := new(request.VolumeManifest)
				err := m.FromYaml(i)
				if err != nil {
					log.Errorf("invalid specification: %s", err.Error())
					return err
				}
				if m.Meta.Name == nil {
					break
				}
				log.Debugf("Add Volume Manifest: %s", *m.Meta.Name)
				mf.Volumes[*m.Meta.Name] = m.GetManifest()
			}
		}
	}

	r.Sync(ctx, mf)
	return nil
}

// Sync node runtime with new spec
func (r *Runtime) Sync(ctx context.Context, spec *types.NodeManifest) error {
	log.V(logLevel).Debugf("%s:sync:> sync runtime state", logNodeRuntimePrefix)
	r.spec <- spec
	return nil
}

// Loop runtime method defines single runtime loop
func (r *Runtime) Loop(ctx context.Context) {

	log.V(logLevel).Debugf("%s:loop:> start runtime loop", logNodeRuntimePrefix)

	go func(ctx context.Context) {

		var network = envs.Get().GetNet()

		for {
			select {
			case spec := <-r.spec:

				log.V(logLevel).Debugf("%s:loop:> provision new spec", logNodeRuntimePrefix)

				if spec.Meta.Initial {

					if network != nil {

						log.V(logLevel).Debugf("%s> clean up endpoints", logNodeRuntimePrefix)
						endpoints := network.Endpoints().GetEndpoints()
						for e := range endpoints {

							// skip resolver endpoint
							if e == network.GetResolverEndpointKey() {
								continue
							}

							if _, ok := spec.Endpoints[e]; !ok {
								network.EndpointDestroy(context.Background(), e, endpoints[e])
							}
						}
					}

					log.V(logLevel).Debugf("%s> clean up pods", logNodeRuntimePrefix)
					pods := envs.Get().GetState().Pods().GetPods()

					for k := range pods {
						if _, ok := spec.Pods[k]; !ok {
							if !envs.Get().GetState().Pods().IsLocal(k) {
								PodDestroy(context.Background(), k, pods[k])
							}
						}
					}

					log.V(logLevel).Debugf("%s> clean up volumes", logNodeRuntimePrefix)
					volumes := envs.Get().GetState().Volumes().GetVolumes()

					for k := range volumes {
						if _, ok := spec.Volumes[k]; !ok {
							if !envs.Get().GetState().Volumes().IsLocal(k) {
								VolumeDestroy(context.Background(), k)
							}
						}
					}

					if network != nil {
						log.V(logLevel).Debugf("%s> clean up subnets", logNodeRuntimePrefix)
						nets := network.Subnets().GetSubnets()

						for cidr := range nets {
							if _, ok := spec.Network[cidr]; !ok {
								network.SubnetDestroy(ctx, cidr)
							}
						}
					}
				}

				if network != nil {
					if len(spec.Resolvers) != 0 {
						log.V(logLevel).Debugf("%s>set cluster dns ips: %#v", logNodeRuntimePrefix, spec.Resolvers)
						for key, res := range spec.Resolvers {
							network.Resolvers().SetResolver(key, res)
							network.ResolverManage(ctx)
						}
					}
				}

				if spec.Exporter != nil {
					log.V(logLevel).Debugf("%s>set cluster exporter endpoint: %s", logNodeRuntimePrefix, spec.Exporter.Endpoint)
					envs.Get().GetExporter().Reconnect(spec.Exporter.Endpoint)
				}

				log.V(logLevel).Debugf("%s> provision init", logNodeRuntimePrefix)

				if network != nil {
					log.V(logLevel).Debugf("%s> provision networks", logNodeRuntimePrefix)
					for cidr, n := range spec.Network {
						log.V(logLevel).Debugf("network: %v", n)
						if err := network.SubnetManage(ctx, cidr, n); err != nil {
							log.Errorf("Subnet [%s] create err: %s", n.CIDR, err.Error())
						}
					}
				}

				log.V(logLevel).Debugf("%s> update secrets %d", logNodeRuntimePrefix, len(spec.Secrets))
				for s, spec := range spec.Secrets {
					log.V(logLevel).Debugf("secret: %s > %s", s, spec.State)
				}

				log.V(logLevel).Debugf("%s> provision configs %d", logNodeRuntimePrefix, len(spec.Configs))
				for s, spec := range spec.Configs {
					log.V(logLevel).Debugf("config: %s > %s", s, spec.State)
					if err := ConfigManage(ctx, s, spec); err != nil {
						log.Errorf("Config [%s] manage err: %s", s, err.Error())
					}
				}

				log.V(logLevel).Debugf("%s> provision pods", logNodeRuntimePrefix)
				for p, spec := range spec.Pods {
					log.V(logLevel).Debugf("pod: %v", p)
					if err := PodManage(ctx, p, spec); err != nil {
						log.Errorf("Pod [%s] manage err: %s", p, err.Error())
					}
				}

				if network != nil {
					log.V(logLevel).Debugf("%s> provision endpoints", logNodeRuntimePrefix)
					for e, spec := range spec.Endpoints {
						log.V(logLevel).Debugf("endpoint: %v", e)
						if err := network.EndpointManage(ctx, e, spec); err != nil {
							log.Errorf("Endpoint [%s] manage err: %s", e, err.Error())
						}
					}
				}

				log.V(logLevel).Debugf("%s> provision volumes", logNodeRuntimePrefix)
				for v, spec := range spec.Volumes {
					log.V(logLevel).Debugf("volume: %v", v)
					if err := VolumeManage(ctx, v, spec); err != nil {
						log.Errorf("Volume [%s] manage err: %s", v, err.Error())
					}
				}
			}
		}
	}(ctx)
}

// Subscribe runtime for container events
func (r *Runtime) Subscribe(ctx context.Context) {

	log.V(logLevel).Debugf("%s:subscribe:> subscribe init", logNodeRuntimePrefix)
	go func() {
		if err := containerSubscribe(ctx); err != nil {
			log.Errorf("container subscribe err: %v", err)
		}
	}()
}

// NewRuntime method return new runtime pointer
func NewRuntime() (*Runtime, error) {

	var (
		r = new(Runtime)
	)

	r.spec = make(chan *types.NodeManifest)
	return r, nil
}
