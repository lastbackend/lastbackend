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
	"strings"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type Runtime struct {
}

func New() *Runtime {
	return new(Runtime)
}

func (r *Runtime) Run() {

	var ctx = context.Background()

	go r.podManifestWatch(ctx, nil)
	go r.volumeManifestWatch(ctx, nil)
	go r.endpointManifestWatch(ctx, nil)
	go r.subnetManifestWatch(ctx, nil)

	go r.secretWatch(ctx, nil)

	go r.nodeWatch(ctx, nil)
	go r.ingressWatch(ctx, nil)
	go r.routeManifestWatch(ctx, nil)

	c := envs.Get().GetCache()

	nm := distribution.NewNamespaceModel(ctx, envs.Get().GetStorage())
	for _, n := range []string{types.SYSTEM_NAMESPACE, types.DEFAULT_NAMESPACE} {
		ns, err := nm.Get(n)
		if err != nil {
			return
		}

		if ns == nil {
			ns = new(types.Namespace)
			ns.Meta.SetDefault()
			ns.Meta.Name = n
			internal, _ := envs.Get().GetDomain()
			ns.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s.%s", ns.Meta.Name, internal))
			ns.Meta.SelfLink = types.NamespaceSelfLink{}
			_ = ns.Meta.SelfLink.Parse(n)

			if _, err := nm.Create(ns); err != nil {
				return
			}
		}

	}

	cm := distribution.NewConfigModel(ctx, envs.Get().GetStorage())
	cl, err := cm.List(types.EmptyString)
	if err != nil {
		return
	}
	go r.configWatch(ctx, &cl.Storage.Revision)

	for _, i := range cl.Items {
		m := new(types.ConfigManifest)
		m.Set(i)
		m.State = types.StateReady
		c.Node().SetConfigManifest(i.SelfLink().String(), m)
	}

	dm := distribution.NewDiscoveryModel(ctx, envs.Get().GetStorage())
	dl, err := dm.List()
	if err != nil {
		return
	}
	go r.discoveryWatch(ctx, &dl.Storage.Revision)
	for _, i := range dl.Items {
		c.Node().SetDiscovery(i)
		c.Ingress().SetDiscovery(i)
	}

	em := distribution.NewExporterModel(ctx, envs.Get().GetStorage())
	el, err := em.List()
	if err != nil {
		return
	}
	go r.exporterWatch(ctx, &cl.Storage.Revision)
	for _, i := range el.Items {
		c.Node().SetExporter(i)
	}

}

func (r *Runtime) podManifestWatch(ctx context.Context, rev *int64) {

	// Watch pods change
	var (
		p = make(chan types.PodManifestEvent)
		c = envs.Get().GetCache()
	)

	mm := distribution.NewPodModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-p:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					c.Node().DelPodManifest(w.Node, w.SelfLink)
					continue
				}

				c.Node().SetPodManifest(w.Node, w.SelfLink, w.Data)
			}
		}
	}()

	mm.ManifestWatch(types.EmptyString, p, rev)
}

func (r *Runtime) volumeManifestWatch(ctx context.Context, rev *int64) {

	// Watch volumes change
	var (
		v = make(chan types.VolumeManifestEvent)
		c = envs.Get().GetCache()
	)

	mm := distribution.NewVolumeModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-v:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					c.Node().DelVolumeManifest(w.Node, w.SelfLink)
					continue
				}

				c.Node().SetVolumeManifest(w.Node, w.SelfLink, w.Data)
			}
		}
	}()

	mm.ManifestWatch(types.EmptyString, v, rev)
}

func (r *Runtime) endpointManifestWatch(ctx context.Context, rev *int64) {

	// Watch volumes change
	var (
		v = make(chan types.EndpointManifestEvent)
		c = envs.Get().GetCache()
	)

	mm := distribution.NewEndpointModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-v:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					w.Data.State = types.StateDestroy
				}

				c.Node().SetEndpointManifest(w.Name, w.Data)
				c.Ingress().SetEndpointManifest(w.Name, w.Data)
			}
		}
	}()

	mm.ManifestWatch(v, rev)
}

func (r *Runtime) subnetManifestWatch(ctx context.Context, rev *int64) {

	// Watch volumes change
	var (
		v = make(chan types.SubnetManifestEvent)
		c = envs.Get().GetCache()
	)

	mm := distribution.NewNetworkModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-v:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					w.Data.State = types.StateDestroy
				}

				c.Node().SetSubnetManifest(w.Name, w.Data)
				c.Ingress().SetSubnetManifest(w.Name, w.Data)
				c.Discovery().SetSubnetManifest(w.Name, w.Data)
			}
		}
	}()

	mm.SubnetManifestWatch(v, rev)
}

func (r *Runtime) secretWatch(ctx context.Context, rev *int64) {

	var (
		n = make(chan types.SecretEvent)
		c = envs.Get().GetCache()
	)

	mm := distribution.NewSecretModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-n:

				if w.Data == nil {
					continue
				}

				sm := new(types.SecretManifest)
				sm.Created = w.Data.Meta.Created
				sm.Updated = w.Data.Meta.Updated
				sm.State = types.StateUpdated

				if w.IsActionRemove() {
					sm := new(types.SecretManifest)
					sm.State = types.StateDestroyed
				}

				c.Node().SetSecretManifest(w.Data.Meta.Name, sm)
			}
		}
	}()

	mm.Watch(n, rev)
}

func (r *Runtime) configWatch(ctx context.Context, rev *int64) {

	var (
		n = make(chan types.ConfigEvent)
		c = envs.Get().GetCache()
	)

	mm := distribution.NewConfigModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-n:

				if w.Data == nil {
					continue
				}

				sm := new(types.ConfigManifest)
				sm.Created = w.Data.Meta.Created
				sm.Updated = w.Data.Meta.Updated
				sm.State = types.StateUpdated

				if w.IsActionRemove() {
					sm := new(types.ConfigManifest)
					sm.State = types.StateDestroyed
				}

				c.Node().SetConfigManifest(w.Data.Meta.Name, sm)
			}
		}
	}()

	mm.Watch(n, rev)
}

func (r *Runtime) nodeWatch(ctx context.Context, rev *int64) {

	// Watch node changes
	var (
		n = make(chan types.NodeEvent)
		c = envs.Get().GetCache()
	)

	mm := distribution.NewNodeModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-n:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					c.Node().Clear(w.Name)
					continue
				}

				if !w.Data.Status.Online {
					c.Node().Clear(w.Name)
				}

			}
		}
	}()

	mm.Watch(n, rev)
}

func (r *Runtime) discoveryWatch(ctx context.Context, rev *int64) {

	// Watch node changes
	var (
		n = make(chan types.DiscoveryEvent)
		c = envs.Get().GetCache()
	)

	im := distribution.NewDiscoveryModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-n:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					c.Node().DelDiscovery(w.Name)
					c.Ingress().DelDiscovery(w.Name)
					continue
				}
				c.Node().SetDiscovery(w.Data)
				c.Ingress().SetDiscovery(w.Data)
			}
		}
	}()

	im.Watch(n, rev)
}

func (r *Runtime) exporterWatch(ctx context.Context, rev *int64) {

	// Watch node changes
	var (
		n = make(chan types.ExporterEvent)
		c = envs.Get().GetCache()
	)

	im := distribution.NewExporterModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-n:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					c.Node().DelExporter(w.Name)
					continue
				}
				c.Node().SetExporter(w.Data)
			}
		}
	}()

	im.Watch(n, rev)
}

func (r *Runtime) ingressWatch(ctx context.Context, rev *int64) {

	// Watch node changes
	var (
		n = make(chan types.IngressEvent)
		c = envs.Get().GetCache()
	)

	im := distribution.NewIngressModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-n:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					c.Ingress().DelIngress(w.Name)
					continue
				}
				c.Ingress().SetIngress(w.Data)

			}
		}
	}()

	im.Watch(n, rev)
}

func (r *Runtime) routeManifestWatch(ctx context.Context, rev *int64) {

	// Watch node changes
	var (
		n = make(chan types.RouteManifestEvent)
		c = envs.Get().GetCache()
	)

	im := distribution.NewRouteModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-n:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					c.Ingress().DelRouteManifest(w.Ingress, w.SelfLink)
					continue
				}

				c.Ingress().SetRouteManifest(w.Ingress, w.SelfLink, w.Data)
			}
		}
	}()

	im.ManifestWatch(types.EmptyString, n, rev)
}
