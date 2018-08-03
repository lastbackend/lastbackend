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

package cluster

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/ipam/ipam"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

// ClusterState is cluster current state struct
type ClusterState struct {
	ec      *EndpointController
	cluster *types.Cluster
	node    struct {
		observer chan *types.Node
		lease    chan *types.NodeLease
		release  chan *types.NodeLease
		list     map[string]*types.Node
	}
	pod struct {
		observer chan *types.Pod
	}
}

// Runtime cluster describes main cluster state loop
func (cs *ClusterState) Observe() {
	// Watch node changes
	for {
		select {
		case n := <-cs.node.observer:
			log.V(7).Debugf("node: %s", n.Meta.Name)
			break
		case p := <-cs.pod.observer:

			switch p.Status.State {
			case types.StateCreated:
				log.Debugf("cluster: Pod provision: %s start", p.SelfLink())
				if err := PodProvision(p, cs); err != nil {
					log.Errorf("%s", err.Error())
				}
				log.Debugf("cluster: Pod provision: %s done", p.SelfLink())

				break
			case types.StateProvision:
				log.Debugf("Pod provision: %s", p.SelfLink())
				if err := PodProvision(p, cs); err != nil {
					log.Errorf("%s", err.Error())
				}
				break
			case types.StateDestroy:
				if err := PodDestroy(p); err != nil {
					log.Errorf("%s", err.Error())
				}
				break
			case types.StateDestroyed:
				if err := PodRemove(p, cs); err != nil {
					log.Errorf("%s", err.Error())
				}
				break
			}
		}
	}
}

// Restore cluster state from storage
func (cs *ClusterState) Restore() error {

	log.Debug("restore cluster state")
	var err error

	// Get cluster info
	cm := distribution.NewClusterModel(context.Background(), envs.Get().GetStorage())
	cs.cluster, err = cm.Get()
	if err != nil {
		return err
	}

	// Get all nodes in cluster
	nm := distribution.NewNodeModel(context.Background(), envs.Get().GetStorage())
	nl, err := nm.List()
	if err != nil {
		return err
	}

	for _, n := range nl {
		// Add node to local cache
		cs.node.list[n.Meta.SelfLink] = n
		// Run node observers
	}

	// Sync cluster state

	// Sync cluster network

	// Sync cluster endpoints

	// Sync cluster manifests

	nn := distribution.NewNamespaceModel(context.Background(), envs.Get().GetStorage())
	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	ns, err := nn.List()
	if err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	for _, n := range ns {
		pl, err := pm.ListByNamespace(n.Meta.Name)
		if err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		for _, p := range pl {
			cs.pod.observer <- p
		}
	}

	return nil
}

// Lease new node for requests by parameters
func (cs *ClusterState) Lease(opts types.NodeLeaseOptions) (*types.Node, error) {
	// Work as node lease requests queue
	req := new(types.NodeLease)
	req.Request = opts

	return req.Get()
}

// Release node
func (cs *ClusterState) Release(opts types.NodeLeaseOptions) (*types.Node, error) {
	// Work as node release
	req := new(types.NodeLease)
	req.Request = opts
	return req.Get()
}

// IPAM management
func (cs *ClusterState) IPAM() ipam.IPAM {
	return envs.Get().GetIPAM()
}

// Endpoint management caller
func (cs *ClusterState) Endpoint() *EndpointController {
	return cs.ec
}

func (cs *ClusterState) SetNode(n *types.Node) {
	cs.node.observer <- n
}

func (cs *ClusterState) DelNode(n *types.Node) {
	delete(cs.node.list, n.Meta.SelfLink)
}

func (cs *ClusterState) SetPod(p *types.Pod) {
	cs.pod.observer <- p
}

// NewClusterState returns new cluster state instance
func NewClusterState() *ClusterState {
	var cs = new(ClusterState)
	cs.ec = new(EndpointController)
	cs.node.observer = make(chan *types.Node)
	cs.node.list = make(map[string]*types.Node)
	cs.node.lease = make(chan *types.NodeLease)
	cs.node.release = make(chan *types.NodeLease)
	cs.node.observer = make(chan *types.Node)
	cs.pod.observer = make(chan *types.Pod)

	go cs.Observe()

	return cs
}
