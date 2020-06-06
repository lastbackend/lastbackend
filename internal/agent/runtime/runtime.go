//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/lastbackend/lastbackend/internal/agent/config"
	"github.com/lastbackend/lastbackend/internal/agent/exporter"
	"github.com/lastbackend/lastbackend/internal/agent/state"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/decoder"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/client/cluster"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/runtime/csi"
	"github.com/lastbackend/lastbackend/tools/logger"
	"gopkg.in/yaml.v3"
)

const (
	logNodeRuntimePrefix = "node:runtime"
	logLevel             = 3
)

// System main node process
type Runtime struct {
	ctx    context.Context
	cancel context.CancelFunc

	csi       map[string]csi.CSI
	cri       cri.CRI
	cii       cii.CII
	network   *network.Network
	state     *state.State
	exporter  *exporter.Exporter
	retClient cluster.IClient

	config config.Config

	spec chan *models.NodeManifest
}

// NewRuntime method return new runtime pointer
func New(cfg config.Config) (*Runtime, error) {

	_net, err := network.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("Can not initialize network: %v", err)
	}

	r := new(Runtime)
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.config = cfg
	r.network = _net
	r.state = state.New()

	r.spec = make(chan *models.NodeManifest, 0)
	return r, nil
}

// NewRuntime run daemon
func (r *Runtime) Run() error {

	_cii, err := cii.New(cii.DockerDriver, cii.DockerConfig{
		Root:    "/var/lib/lastbackend/storage",
		RunRoot: "/var/run/lastbackend/storage",
		StorageDriver: "overlay",
	})
	if err != nil {
		return fmt.Errorf("Cannot initialize iri: %v", err)
	}

	_cri, err := cri.New(cri.RuncDriver, cri.RuncConfig{})
	if err != nil {
		return fmt.Errorf("Cannot initialize cri: %v", err)
	}

	_csi := make(map[string]csi.CSI, 0)

	// TODO: Implement csi initialization logic
	//csis := app.config.GetStringMap("container.csi")
	//if csis != nil {
	//	for kind := range csis {
	//		si, err := csi.New(kind, dir.Config{RootDir: filepath.Join(app.config.WorkDir, "csi")})
	//		if err != nil {
	//			log.Errorf("Cannot initialize [%s] csi: %v", kind, err)
	//			return err
	//		}
	//		csi[kind] = si
	//	}
	//}

	// TODO: Implement cluster state logic
	//_state.Node().Info = runtime.NodeInfo()
	//_state.Node().Status = runtime.NodeStatus()

	exp, err := exporter.NewExporter(r.state.Node().Info.Hostname, models.EmptyString)
	if err != nil {
		return fmt.Errorf("Can not initialize collector: %v", err)
	}

	r.csi = _csi
	r.cri = _cri
	r.cii = _cii
	r.exporter = exp

	if err := r.Restore(); err != nil {
		return err
	}

	r.Subscribe()
	r.Loop()

	if r.config.ManifestDir != models.EmptyString {
		r.Provision(r.config.ManifestDir)
	}

	go r.exporter.Listen()

	return nil
}

func (r *Runtime) Restore() error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s:restore:> restore init", logNodeRuntimePrefix)

	if r.network != nil {
		if err := r.network.SubnetRestore(r.ctx); err != nil {
			log.Errorf("can not restore subnets: %s", err.Error())
			return err
		}

		if err := r.network.EndpointRestore(r.ctx); err != nil {
			log.Errorf("can not restore endpoint: %s", err.Error())
			return err
		}

		if err := r.network.ResolverManage(r.ctx); err != nil {
			log.Errorf("%s:> can not manage resolver:%s", logNodeRuntimePrefix, err.Error())
		}
	}

	if err := r.VolumeRestore(r.ctx); err != nil {
		log.Errorf("can not restore volumes: %s", err.Error())
		return err
	}

	if err := r.ImageRestore(r.ctx); err != nil {
		log.Errorf("Can not restore images: %s", err.Error())
		return err
	}

	if err := r.PodRestore(r.ctx); err != nil {
		log.Errorf("Can not restore pods: %s", err.Error())
		return err
	}

	return nil

}

// Provision node manifest
func (r *Runtime) Provision(dir string) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s:provision:> local init", logNodeRuntimePrefix)

	log.Debugf("%s:provision:> read manifests from dir: %s", logNodeRuntimePrefix, dir)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var (
		mf = new(models.NodeManifest)
	)
	mf.Configs = make(map[string]*models.ConfigManifest)
	mf.Pods = make(map[string]*models.PodManifest)
	mf.Volumes = make(map[string]*models.VolumeManifest)

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
			case models.KindConfig:
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
			case models.KindPod:

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
				r.state.Pods().SetLocal(*m.Meta.Name)
				break
			case models.KindVolume:

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

	r.Sync(mf)
	return nil
}

// Sync node runtime with new spec
func (r *Runtime) Sync(spec *models.NodeManifest) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s:sync:> sync runtime state", logNodeRuntimePrefix)
	r.spec <- spec
	return nil
}

// Loop runtime method defines single runtime loop
func (r *Runtime) Loop() {
	log := logger.WithContext(context.Background())
	log.Debugf("%s:loop:> start runtime loop", logNodeRuntimePrefix)

	go func(ctx context.Context) {

		for {
			select {
			case spec := <-r.spec:

				log.Debugf("%s:loop:> provision new spec", logNodeRuntimePrefix)

				if spec.Meta.Initial {

					if r.network != nil {

						log.Debugf("%s:> clean up endpoints", logNodeRuntimePrefix)
						endpoints := r.network.Endpoints().GetEndpoints()
						for e := range endpoints {

							// skip resolver endpoint
							if e == r.network.GetResolverEndpointKey() {
								continue
							}

							if _, ok := spec.Endpoints[e]; !ok {
								r.network.EndpointDestroy(context.Background(), e, endpoints[e])
							}
						}
					}

					log.Debugf("%s:> clean up pods", logNodeRuntimePrefix)
					pods := r.state.Pods().GetPods()

					for k := range pods {
						if _, ok := spec.Pods[k]; !ok {
							if !r.state.Pods().IsLocal(k) {
								r.PodDestroy(context.Background(), k, pods[k])
							}
						}
					}

					log.Debugf("%s:> clean up volumes", logNodeRuntimePrefix)
					volumes := r.state.Volumes().GetVolumes()

					for k := range volumes {
						if _, ok := spec.Volumes[k]; !ok {
							if !r.state.Volumes().IsLocal(k) {
								r.VolumeDestroy(context.Background(), k)
							}
						}
					}

					if r.network != nil {
						log.Debugf("%s:> clean up subnets", logNodeRuntimePrefix)
						nets := r.network.Subnets().GetSubnets()

						for cidr := range nets {
							if _, ok := spec.Network[cidr]; !ok {
								r.network.SubnetDestroy(ctx, cidr)
							}
						}
					}
				}

				if r.network != nil {
					if len(spec.Resolvers) != 0 {
						log.Debugf("%s:> set cluster dns ips: %#v", logNodeRuntimePrefix, spec.Resolvers)
						for key, res := range spec.Resolvers {
							r.network.Resolvers().SetResolver(key, res)
							r.network.ResolverManage(ctx)
						}
					}
				}

				if spec.Exporter != nil {
					log.Debugf("%s:> set cluster exporter endpoint: %s", logNodeRuntimePrefix, spec.Exporter.Endpoint)
					r.exporter.Reconnect(spec.Exporter.Endpoint)
				}

				log.Debugf("%s:> provision init", logNodeRuntimePrefix)

				if r.network != nil {
					log.Debugf("%s:> provision networks", logNodeRuntimePrefix)
					for cidr, n := range spec.Network {
						log.Debugf("network: %v", n)
						if err := r.network.SubnetManage(ctx, cidr, n); err != nil {
							log.Errorf("Subnet [%s] create err: %s", n.CIDR, err.Error())
						}
					}
				}

				log.Debugf("%s:> update secrets %d", logNodeRuntimePrefix, len(spec.Secrets))
				for s, spec := range spec.Secrets {
					log.Debugf("secret: %s > %s", s, spec.State)
				}

				log.Debugf("%s:> provision configs %d", logNodeRuntimePrefix, len(spec.Configs))
				for s, spec := range spec.Configs {
					log.Debugf("config: %s > %s", s, spec.State)
					if err := r.ConfigManage(ctx, s, spec); err != nil {
						log.Errorf("Config [%s] manage err: %s", s, err.Error())
					}
				}

				log.Debugf("%s:> provision pods", logNodeRuntimePrefix)
				for p, spec := range spec.Pods {
					log.Debugf("pod: %v", p)
					if err := r.PodManage(ctx, p, spec); err != nil {
						log.Errorf("Pod [%s] manage err: %s", p, err.Error())
					}
				}

				if r.network != nil {
					log.Debugf("%s:> provision endpoints", logNodeRuntimePrefix)
					for e, spec := range spec.Endpoints {
						log.Debugf("endpoint: %v", e)
						if err := r.network.EndpointManage(ctx, e, spec); err != nil {
							log.Errorf("Endpoint [%s] manage err: %s", e, err.Error())
						}
					}
				}

				log.Debugf("%s:> provision volumes", logNodeRuntimePrefix)
				for v, spec := range spec.Volumes {
					log.Debugf("volume: %v", v)
					if err := r.VolumeManage(ctx, v, spec); err != nil {
						log.Errorf("Volume [%s] manage err: %s", v, err.Error())
					}
				}
			}
		}
	}(r.ctx)
}

// Subscribe runtime for container events
func (r *Runtime) Subscribe() {
	log := logger.WithContext(context.Background())
	log.Debugf("%s:subscribe:> subscribe init", logNodeRuntimePrefix)
	go func() {
		if err := r.containerSubscribe(r.ctx); err != nil {
			log.Errorf("container subscribe err: %v", err)
		}
	}()
}

func (r *Runtime) Stop() {
	r.cancel()
}

func (r *Runtime) GetConfig() config.Config {
	return r.config
}

func (r *Runtime) GetState() *state.State {
	return r.state
}

func (r *Runtime) GetNetwork() *network.Network {
	return r.network
}
