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

package cluster

import (
	"context"

	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/master/ipam"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel  = 3
	logPrefix = "observer:cluster"
)

// ClusterState is cluster current state struct
type ClusterState struct {
	storage storage.IStorage
	ipam    ipam.IPAM

	cluster *models.Cluster

	ingress struct {
		observer chan *models.Ingress
		list     map[string]*models.Ingress
	}

	discovery struct {
		observer chan *models.Discovery
		list     map[string]*models.Discovery
	}

	route struct {
		ingress  map[string]int
		observer chan *models.Route
		list     map[string]*models.Route
	}

	volume struct {
		observer chan *models.Volume
		list     map[string]*models.Volume
	}

	node struct {
		observer chan *models.Node
		lease    chan *NodeLease
		release  chan *NodeLease
		list     map[string]*models.Node
	}
}

// System cluster describes main cluster state loop
func (cs *ClusterState) Observe() {
	// Watch node changes
	for {
		select {
		case l := <-cs.node.lease:
			_ = handleNodeLease(cs, l)
			break
		case l := <-cs.node.release:
			_ = handleNodeRelease(cs, l)
			break
		case n := <-cs.node.observer:
			log.V(7).Debugf("node: %s", n.Meta.Name)
			cs.node.list[n.SelfLink().String()] = n
			_ = clusterStatusState(cs)
			break
		case v := <-cs.volume.observer:
			log.V(7).Debugf("volume: %s", v.SelfLink().String())
			if err := volumeObserve(cs, v); err != nil {
				log.Errorf("%s", err.Error())
			}
			break
		case r := <-cs.route.observer:
			log.V(7).Debugf("route: %s", r.SelfLink().String())
			if err := routeObserve(cs, r); err != nil {
				log.Errorf("%s", err.Error())
			}
			break
		case i := <-cs.ingress.observer:
			log.V(7).Debugf("ingress: %s", i.SelfLink().String())
			cs.ingress.list[i.SelfLink().String()] = i
			break
		}
	}
}

// Loop cluster state from storage
func (cs *ClusterState) Loop() error {

	log.Debug("restore cluster state")
	var err error

	// Get cluster info
	cm := service.NewClusterModel(context.Background(), cs.storage)
	cs.cluster, err = cm.Get()
	if err != nil {
		return err
	}

	// Get all nodes in cluster
	nm := service.NewNodeModel(context.Background(), cs.storage)
	nl, err := nm.List()
	if err != nil {
		return err
	}
	for _, n := range nl.Items {
		// Add node to local cache
		cs.SetNode(n)
		// Run node observers
	}

	_ = clusterStatusState(cs)

	// Get all ingress servers in cluster
	im := service.NewIngressModel(context.Background(), cs.storage)
	il, err := im.List()
	if err != nil {
		return err
	}
	for _, i := range il.Items {
		// Add ingress to local cache
		cs.SetIngress(i)
		// Run ingress observers
	}

	// Get all routes in cluster
	rm := service.NewRouteModel(context.Background(), cs.storage)
	rl, err := rm.List()
	if err != nil {
		return err
	}
	for _, r := range rl.Items {
		// Add route to local cache
		cs.SetRoute(r)
		// Run route observers
	}

	go cs.watchNode(context.Background(), &nl.Storage.Revision)
	go cs.watchIngress(context.Background(), &il.Storage.Revision)
	go cs.watchRoute(context.Background(), &rl.Storage.Revision)

	return nil
}

func (cs *ClusterState) watchNode(ctx context.Context, rev *int64) {

	var (
		p = make(chan models.NodeEvent)
	)

	nm := service.NewNodeModel(ctx, cs.storage)

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
					cs.DelNode(w.Data)
					continue
				}

				cs.SetNode(w.Data)
			}
		}
	}()

	nm.Watch(p, rev)
}

func (cs *ClusterState) watchIngress(ctx context.Context, rev *int64) {

	var (
		p = make(chan models.IngressEvent)
	)

	nm := service.NewIngressModel(ctx, cs.storage)

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
					cs.DelIngress(w.Data)
					continue
				}

				cs.SetIngress(w.Data)
			}
		}
	}()

	nm.Watch(p, rev)
}

func (cs *ClusterState) watchRoute(ctx context.Context, rev *int64) {
	var (
		p = make(chan models.RouteEvent)
	)

	rm := service.NewRouteModel(ctx, cs.storage)

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
					cs.DelRoute(w.Data)
					continue
				}

				cs.SetRoute(w.Data)
			}
		}
	}()

	rm.Watch(p, rev)
}

// lease new node for requests by parameters
func (cs *ClusterState) lease(opts NodeLeaseOptions) (*models.Node, error) {

	// Work as node lease requests queue
	req := new(NodeLease)
	req.Request = opts
	req.done = make(chan bool)
	cs.node.lease <- req
	req.Wait()

	return req.Response.Node, req.Response.Err
}

// release node
func (cs *ClusterState) release(opts NodeLeaseOptions) (*models.Node, error) {
	// Work as node release
	req := new(NodeLease)
	req.Request = opts
	req.done = make(chan bool)
	cs.node.release <- req
	req.Wait()
	return req.Response.Node, req.Response.Err
}

// lease node in sync mode
func (cs *ClusterState) leaseSync(opts NodeLeaseOptions) (*models.Node, error) {

	// Work as node lease requests queue
	req := new(NodeLease)
	req.Request = opts
	req.sync = true
	if err := handleNodeLease(cs, req); err != nil {
		log.Errorf("sync lease error: %s", err.Error())
		return nil, err
	}

	return req.Response.Node, req.Response.Err
}

// release node in sync mode
func (cs *ClusterState) releaseSync(opts NodeLeaseOptions) (*models.Node, error) {
	// Work as node release
	req := new(NodeLease)
	req.Request = opts
	req.sync = true

	if err := handleNodeRelease(cs, req); err != nil {
		log.Errorf("sync release error: %s", err.Error())
		return nil, err
	}

	return req.Response.Node, req.Response.Err
}

// IPAM management
func (cs *ClusterState) IPAM() ipam.IPAM {
	return cs.ipam
}

func (cs *ClusterState) SetNode(n *models.Node) {
	cs.node.observer <- n
}

func (cs *ClusterState) DelNode(n *models.Node) {
	delete(cs.node.list, n.SelfLink().String())
}

func (cs *ClusterState) SetIngress(i *models.Ingress) {
	cs.ingress.observer <- i
}

func (cs *ClusterState) DelIngress(i *models.Ingress) {
	delete(cs.ingress.list, i.SelfLink().String())
}

func (cs *ClusterState) SetDiscovery(d *models.Discovery) {
	cs.discovery.observer <- d
}

func (cs *ClusterState) DelDiscovery(d *models.Discovery) {
	delete(cs.discovery.list, d.SelfLink().String())
}

func (cs *ClusterState) SetVolume(v *models.Volume) {
	cs.volume.observer <- v
}

func (cs *ClusterState) DelVolume(v *models.Volume) {
	delete(cs.volume.list, v.SelfLink().String())
}

func (cs *ClusterState) SetRoute(r *models.Route) {
	cs.route.observer <- r
}

func (cs *ClusterState) DelRoute(r *models.Route) {
	delete(cs.route.list, r.SelfLink().String())
}

func (cs *ClusterState) PodLease(p *models.Pod) (*models.Node, error) {

	var RAM int64

	for _, s := range p.Spec.Template.Containers {
		RAM += s.Resources.Request.RAM
	}

	opts := NodeLeaseOptions{
		Selector: p.Spec.Selector,
		RAM:      &RAM,
	}

	node, err := cs.lease(opts)
	if err != nil {
		log.Errorf("%s:> pod lease err: %s", logPrefix, err)
		return nil, err
	}

	return node, err
}

func (cs *ClusterState) PodRelease(p *models.Pod) (*models.Node, error) {
	var RAM int64

	for _, s := range p.Spec.Template.Containers {
		RAM += s.Resources.Request.RAM
	}

	opts := NodeLeaseOptions{
		Node: &p.Meta.Node,
		RAM:  &RAM,
	}

	node, err := cs.release(opts)
	if err != nil {
		log.Errorf("%s:> pod lease err: %s", logPrefix, err)
		return nil, err
	}

	return node, err
}

func (cs *ClusterState) VolumeLease(v *models.Volume) (*models.Node, error) {

	opts := NodeLeaseOptions{
		Node:     &v.Spec.Selector.Node,
		Selector: v.Spec.Selector,
		Storage:  &v.Spec.Capacity.Storage,
	}

	node, err := cs.leaseSync(opts)
	if err != nil {
		log.Errorf("%s:> volume lease err: %s", logPrefix, err)
		return nil, err
	}

	return node, err
}

func (cs *ClusterState) VolumeRelease(v *models.Volume) (*models.Node, error) {

	opts := NodeLeaseOptions{
		Node:    &v.Meta.Node,
		Storage: &v.Spec.Capacity.Storage,
	}

	node, err := cs.releaseSync(opts)
	if err != nil {
		log.Errorf("%s:> volume lease err: %s", logPrefix, err)
		return nil, err
	}

	return node, err
}

// NewClusterState returns new cluster state instance
func NewClusterState(stg storage.IStorage, ipam ipam.IPAM) *ClusterState {

	var cs = new(ClusterState)

	cs.storage = stg
	cs.ipam = ipam

	cs.cluster = new(models.Cluster)

	cs.ingress.observer = make(chan *models.Ingress)
	cs.ingress.list = make(map[string]*models.Ingress)

	cs.discovery.observer = make(chan *models.Discovery)
	cs.discovery.list = make(map[string]*models.Discovery)

	cs.volume.list = make(map[string]*models.Volume)
	cs.volume.observer = make(chan *models.Volume)

	cs.node.observer = make(chan *models.Node)
	cs.node.list = make(map[string]*models.Node)

	cs.node.lease = make(chan *NodeLease)
	cs.node.release = make(chan *NodeLease)

	cs.node.observer = make(chan *models.Node)

	cs.route.observer = make(chan *models.Route)
	cs.route.ingress = make(map[string]int, 0)
	cs.route.list = make(map[string]*models.Route, 0)

	go cs.Observe()

	return cs
}
