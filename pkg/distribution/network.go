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

package distribution

import (
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logNetworkPrefix = "distribution:network"
)

type Network struct {
	context context.Context
	storage storage.Storage
}

func (s *Network) Runtime() (*types.System, error) {

	log.V(logLevel).Debugf("%s:get:> get network runtime info", logNetworkPrefix)
	runtime, err := s.storage.Info(s.context, s.storage.Collection().Network(), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> get runtime info error: %s", logNetworkPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil

}

// Get network info
func (s *Network) Get() (*types.Network, error) {

	log.V(logLevel).Debugf("%s:get:> get network info", logNetworkPrefix)

	net := new(types.Network)

	err := s.storage.Get(s.context, s.storage.Collection().Network(), types.EmptyString, net, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get network info not found", logNetworkPrefix)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get network info  error: %v", logNetworkPrefix, err)
		return nil, err
	}

	return net, nil
}

// Create new network info
func (s *Network) Put(net *types.Network) (*types.Network, error) {

	log.V(logLevel).Debugf("%s:create:> put new network %#v", logNetworkPrefix)

	if err := s.storage.Put(s.context, s.storage.Collection().Network(), types.EmptyString, net, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert network err: %v", logNetworkPrefix, err)
		return nil, err
	}

	return net, nil
}

// Update network in namespace
func (s *Network) Set(net *types.Network) (*types.Network, error) {

	log.V(logLevel).Debugf("%s:create:> put new network %#v", logNetworkPrefix)

	if err := s.storage.Set(s.context, s.storage.Collection().Network(), types.EmptyString, net, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert network err: %v", logNetworkPrefix, err)
		return nil, err
	}

	return net, nil
}

// Remove network from storage
func (s *Network) Del(net *types.Network) error {

	log.V(logLevel).Debugf("%s:remove:> remove network %#v", logNetworkPrefix, net)

	err := s.storage.Del(s.context, s.storage.Collection().Network(), types.EmptyString)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove network err: %v", logNetworkPrefix, err)
		return err
	}

	return nil
}

// Watch network changes
func (s *Network) Watch(ch chan types.NetworkEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch network by spec changes", logNetworkPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-s.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.NetworkEvent{}
				res.Action = e.Action
				res.Name = e.Name

				network := new(types.Network)

				if err := json.Unmarshal(e.Data.([]byte), network); err != nil {
					log.Errorf("%s:> parse data err: %v", logNetworkPrefix, err)
					continue
				}

				res.Data = network

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := s.storage.Watch(s.context, s.storage.Collection().Network(), watcher, opts); err != nil {
		return err
	}

	return nil
}

// Get subnet list
func (s *Network) SubnetList() ([]*types.Subnet, error) {

	log.V(logLevel).Debugf("%s:SubnetList:> get snets %s", logNetworkPrefix)

	snets := make([]*types.Subnet, 0)

	err := s.storage.List(s.context, s.storage.Collection().Subnet(), "", snets, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:SubnetList:> getsnets not found", logNetworkPrefix)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:SubnetList:> get by name %s error: %v", logNetworkPrefix, err)
		return nil, err
	}

	return snets, nil
}

// Get subnet by name
func (s *Network) SubnetGet(cidr string) (*types.Subnet, error) {

	log.V(logLevel).Debugf("%s:SubnetGet:> get by name %s", logNetworkPrefix, cidr)

	name := types.SubnetGetNameFromCIDR(cidr)
	snet := new(types.Subnet)
	key := types.NewSubnetSelfLink(cidr).String()

	err := s.storage.Get(s.context, s.storage.Collection().Subnet(), key, snet, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:SubnetGet:> get by name %s not found", logNetworkPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:SubnetGet:> get by name %s error: %v", logNetworkPrefix, name, err)
		return nil, err
	}

	return snet, nil
}

// Create new subnet
func (s *Network) SubnetPut(hostname string, spec types.SubnetSpec) (*types.Subnet, error) {

	log.V(logLevel).Debugf("%s:SubnetPut:> put new subnet", logNetworkPrefix)

	snet := new(types.Subnet)
	snet.Meta.SetDefault()
	snet.Meta.Name = types.SubnetGetNameFromCIDR(spec.CIDR)
	snet.Meta.Node = hostname
	snet.Meta.SelfLink = *types.NewSubnetSelfLink(spec.CIDR)
	snet.Spec = spec

	if err := s.storage.Put(s.context, s.storage.Collection().Subnet(),
		snet.SelfLink().String(), snet, nil); err != nil {
		log.V(logLevel).Errorf("%s:SubnetPut:> insert subnet err: %v", logNetworkPrefix, err)
		return nil, err
	}

	if err := s.SubnetManifestAdd(snet); err != nil {
		log.V(logLevel).Errorf("%s:SubnetPut:> insert subnet manifest err: %v", logNetworkPrefix, err)
		return nil, err
	}

	return snet, nil
}

// Update subnet
func (s *Network) SubnetSet(snet *types.Subnet) error {

	log.V(logLevel).Debugf("%s:SubnetSet:> set subnet", logNetworkPrefix)

	if err := s.storage.Set(s.context, s.storage.Collection().Subnet(),
		snet.SelfLink().String(), snet, nil); err != nil {
		log.V(logLevel).Errorf("%s:SubnetSet:> err: %v", logNetworkPrefix, err)
		return err
	}

	m, err := s.SubnetManifestGet(snet.SelfLink().String())
	if err != nil {
		log.V(logLevel).Errorf("%s:SubnetSet:> get manifest err: %v", logNetworkPrefix, err)
		return err
	}

	if m == nil {
		return s.SubnetManifestAdd(snet)
	}

	if !types.SubnetSpecEqual(&m.SubnetSpec, &snet.Spec) {
		if err := s.SubnetManifestSet(m, snet); err != nil {
			log.V(logLevel).Errorf("%s:SubnetPut:> insert subnet manifest err: %v", logNetworkPrefix, err)
			return err
		}
	}

	return nil
}

// Remove subnet
func (s *Network) SubnetDel(name string) error {

	log.V(logLevel).Debugf("%s:SubnetDel:> remove subnet", logNetworkPrefix)
	key := types.NewSubnetSelfLink(name).String()

	err := s.storage.Del(s.context, s.storage.Collection().Network(), key)
	if err != nil {
		log.V(logLevel).Errorf("%s:SubnetDel:> remove subnet err: %v", logNetworkPrefix, err)
		return err
	}

	if err := s.SubnetManifestDel(name); err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Errorf("%s:SubnetDel:> get manifest err: %v", logNetworkPrefix, err)
			return err
		}
	}

	return nil
}

// Check subnet
func (s *Network) SubnetEqual(snet *types.Subnet, spec types.SubnetSpec) bool {

	if snet.Spec.CIDR != spec.CIDR {
		return false
	}

	if snet.Spec.Type != spec.Type {
		return false
	}

	if snet.Spec.Addr != spec.Addr {
		return false
	}

	if snet.Spec.IFace.Addr != spec.IFace.Addr {
		return false
	}

	if snet.Spec.IFace.Name != spec.IFace.Name {
		return false
	}

	if snet.Spec.IFace.HAddr != spec.IFace.HAddr {
		return false
	}

	if snet.Spec.IFace.Index != spec.IFace.Index {
		return false
	}

	return true
}

// Watch network changes
func (s *Network) SubnetWatch(ch chan types.SubnetEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch subnet spec changes ", logNetworkPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-s.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.SubnetEvent{}
				res.Action = e.Action
				res.Name = e.Name

				network := new(types.Subnet)

				if err := json.Unmarshal(e.Data.([]byte), network); err != nil {
					log.Errorf("%s:> parse data err: %v", logNetworkPrefix, err)
					continue
				}

				res.Data = network

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := s.storage.Watch(s.context, s.storage.Collection().Subnet(), watcher, opts); err != nil {
		return err
	}

	return nil
}

// Get network subnet manifests map
func (s *Network) SubnetManifestMap() (*types.SubnetManifestMap, error) {
	log.V(logLevel).Debugf("%s:SubnetManifestMap:> ", logNetworkPrefix)

	var (
		mf = types.NewSubnetManifestMap()
	)

	if err := s.storage.Map(s.context, s.storage.Collection().Manifest().Subnet(), types.EmptyString, mf, nil); err != nil {
		log.Errorf("%s:SubnetManifestMap:> err :%s", logNetworkPrefix, err.Error())
		return nil, err
	}

	return mf, nil
}

// Get particular network manifest
func (s *Network) SubnetManifestGet(name string) (*types.SubnetManifest, error) {
	log.V(logLevel).Debugf("%s:SubnetManifestGet:> ", logNetworkPrefix)

	var (
		mf = new(types.SubnetManifest)
	)

	if err := s.storage.Get(s.context, s.storage.Collection().Manifest().Subnet(), name, &mf, nil); err != nil {
		log.Errorf("%s:SubnetManifestGet:> err :%s", logNetworkPrefix, err.Error())

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return mf, nil
}

// Add particular network manifest
func (s *Network) SubnetManifestAdd(snet *types.Subnet) error {
	log.V(logLevel).Debugf("%s:SubnetManifestAdd:> ", logNetworkPrefix)

	m := new(types.SubnetManifest)
	m.SubnetSpec = snet.Spec

	if err := s.storage.Put(s.context, s.storage.Collection().Manifest().Subnet(), snet.SelfLink().String(),
		m, nil); err != nil {
		log.Errorf("%s:SubnetManifestAdd:> err :%s", logNetworkPrefix, err.Error())
		return err
	}

	return nil
}

// Set particular network manifest
func (s *Network) SubnetManifestSet(m *types.SubnetManifest, snet *types.Subnet) error {
	log.V(logLevel).Debugf("%s:SubnetManifestAdd:> ", logNetworkPrefix)
	m.SubnetSpec = snet.Spec
	if err := s.storage.Set(s.context, s.storage.Collection().Manifest().Subnet(), snet.SelfLink().String(), m, nil); err != nil {
		log.Errorf("%s:SubnetManifestAdd:> err :%s", logNetworkPrefix, err.Error())
		return err
	}
	return nil
}

// Del particular network manifest
func (s *Network) SubnetManifestDel(name string) error {
	log.V(logLevel).Debugf("%s:SubnetManifestDel:> ", logNetworkPrefix)

	if err := s.storage.Del(s.context, s.storage.Collection().Manifest().Subnet(), name); err != nil {
		log.Errorf("%s:SubnetManifestDel:> err :%s", logNetworkPrefix, err.Error())
		return err
	}

	return nil
}

// watch subnet manifests
func (s *Network) SubnetManifestWatch(ch chan types.SubnetManifestEvent, rev *int64) error {
	log.V(logLevel).Debugf("%s:SubnetManifestWatch:> watch manifest ", logNetworkPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-s.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.SubnetManifestEvent{}
				res.Action = e.Action
				res.Name = e.Name
				res.SelfLink = e.SelfLink

				manifest := new(types.SubnetManifest)

				if err := json.Unmarshal(e.Data.([]byte), manifest); err != nil {
					log.Errorf("%s:> parse data err: %v", logNetworkPrefix, err)
					continue
				}

				res.Data = manifest

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := s.storage.Watch(s.context, s.storage.Collection().Manifest().Subnet(), watcher, opts); err != nil {
		return err
	}

	return nil
}

// NewNetworkModel returns new network management model
func NewNetworkModel(ctx context.Context, stg storage.Storage) *Network {
	return &Network{ctx, stg}
}
