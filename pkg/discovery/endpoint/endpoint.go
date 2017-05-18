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

package endpoint

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util"
	"net"
	"strings"
)

func Get(name string) ([]net.IP, error) {

	var (
		err     error
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage()
		cache   = context.Get().GetCache()
		data    = []string{}
	)

	log.Debugf(`Endpoint: Get ip list from cache %s`, name)

	var ips = cache.Get(name)

	if ips != nil && len(ips) != 0 {
		return ips, nil
	}

	if len(data) == 0 {
		log.Debugf(`Endpoint: Try find to db %s`, name)

		result, err := storage.Endpoint().Get(context.Get().Background(), name)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		for _, ip := range result {
			data = append(data, ip)
		}
	}

	data = util.RemoveDuplicates(data)

	ips, err = util.ConvertStringIPToNetIP(data)
	if err != nil {
		log.Errorf("Endpoint: convert ips to net ips error %s", err.Error())
		return ips, err
	}

	log.Debug(`ips`, ips)

	cache.Lock()
	cache.Set(name, ips)
	cache.Unlock()

	log.Debugf(`Endpoint: Get ip list from cache for %s successfully: %v`, name, data)

	return ips, nil
}

func Update(name string) error {

	var (
		err     error
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage()
		cache   = context.Get().GetCache()
		data    = []string{}
	)

	log.Debugf(`Endpoint: Update name %s in cache `, name)

	result, err := storage.Endpoint().Get(context.Get().Background(), name)
	if err != nil {
		log.Error(err)
		return err
	}

	for _, ip := range result {
		data = append(data, ip)
	}

	data = util.RemoveDuplicates(data)

	ips, err := util.ConvertStringIPToNetIP(data)
	if err != nil {
		log.Error(err)
		return err
	}

	cache.Lock()
	cache.Set(name, ips)
	cache.Unlock()

	log.Debugf(`Endpoint: Update name %s in cache successfully `, name)

	return nil
}

func Remove(name string) error {

	var (
		err   error
		log   = context.Get().GetLogger()
		cache = context.Get().GetCache()
	)

	log.Debugf(`Endpoint: Remove name %s from cache `, name)

	err = cache.Del(name)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf(`Endpoint: Remove name %s from cache successfully `, name)

	return nil
}

func Loop() {
	var (
		log       = context.Get().GetLogger()
		storage   = context.Get().GetStorage()
		cache     = context.Get().GetCache()
		endpoints = make(chan string, 100)
		services  = make(chan *types.Service, 100)
		pods      = make(chan *types.Pod, 100)
	)

	log.Debug("Endpoint: Endpoint: Watch")

	go func() {
		for {
			select {
			case e := <-endpoints:
				{
					fmt.Println("===================================")
					fmt.Println("endpoints", e)
					fmt.Println("===================================")
					i, err := storage.Endpoint().Get(context.Get().Background(), e)
					if err != nil {
						if err.Error() != store.ErrKeyNotFound {
							if err = cache.Del(e); err != nil {
								log.Errorf("Endpoint: remove ips from cache error %s", err.Error())
							}
						} else {
							log.Errorf("Endpoint: get ips for domain error %s", err.Error())
						}
						continue
					}

					ips, err := util.ConvertStringIPToNetIP(i)
					if err != nil {
						log.Errorf("Endpoint: convert ips to net ips error %s", err.Error())
						continue
					}

					if err = cache.Set(e, ips); err != nil {
						log.Errorf("Endpoint: save ips to cache error %s", err.Error())
						continue
					}
				}
			case s := <-services:
				{
					fmt.Println("===================================")
					fmt.Println("service", s)
					fmt.Println("===================================")

					if s == nil {
						continue
					}

					serviceEndpoint := fmt.Sprintf("%s-%s.%s", s.Meta.Name, s.Meta.Namespace, *context.Get().GetConfig().SystemDomain)
					serviceEndpoint = strings.Replace(serviceEndpoint, ":", "-", -1)

					hosts := make(map[string]string)
					ips := []string{}
					for _, pod := range s.Pods {
						if _, ok := hosts[pod.Meta.Hostname]; ok || pod.State.State == types.StateDestroy {
							continue
						}

						node, err := storage.Node().Get(context.Get().Background(), pod.Meta.Hostname)
						if err != nil {
							log.Errorf("Endpoint: get node error %s", err.Error())
							break
						}

						hosts[pod.Meta.Hostname] = node.Meta.IP
						ips = append(ips, node.Meta.IP)
					}

					if s.State.State == types.StateDestroy {
						if err := storage.Endpoint().Remove(context.Get().Background(), serviceEndpoint); err != nil {
							log.Errorf("Endpoint: remove service endpoint error %s", err.Error())
						}
						continue
					}

					if err := storage.Endpoint().Upsert(context.Get().Background(), serviceEndpoint, ips); err != nil {
						log.Errorf("Endpoint: upsert service endpoint error %s", err.Error())
						continue
					}
				}
			case p := <-pods:
				{
					fmt.Println("===================================")
					fmt.Println("pod", p)
					fmt.Println("===================================")

					if p == nil || /*!p.State.Provision ||*/ p.Meta.Hostname == "" {
						continue
					}

					podEndpoint := fmt.Sprintf("%s.%s", p.Meta.Name, *context.Get().GetConfig().SystemDomain)
					podEndpoint = strings.Replace(podEndpoint, ":", "-", -1)

					srv, err := storage.Service().GetByPodName(context.Get().Background(), p.Meta.Name)
					if err != nil {
						if err.Error() == store.ErrKeyNotFound {
							if err := storage.Endpoint().Remove(context.Get().Background(), podEndpoint); err != nil {
								log.Errorf("Endpoint: remove endpoint error %s", err.Error())
							}
						} else {
							log.Errorf("Endpoint: get service error %s", err.Error())
						}
						continue
					}

					node, err := storage.Node().Get(context.Get().Background(), p.Meta.Hostname)
					if err != nil {
						log.Errorf("Endpoint: get node error %s", err.Error())
						break
					}

					serviceEndpoint := fmt.Sprintf("%s-%s.%s", srv.Meta.Name, srv.Meta.Namespace, *context.Get().GetConfig().SystemDomain)
					serviceEndpoint = strings.Replace(serviceEndpoint, ":", "-", -1)

					fmt.Println("pod state", p.State.State, p.State.State == types.StateDestroy)

					if p.State.State == types.StateDestroy {
						if err := storage.Endpoint().Remove(context.Get().Background(), podEndpoint); err != nil {
							log.Errorf("Endpoint: remove endpoint error %s", err.Error())
						}
						continue
					}

					if err := storage.Endpoint().Upsert(context.Get().Background(), podEndpoint, []string{node.Meta.IP}); err != nil {
						log.Errorf("Endpoint: upsert endpoint error %s", err.Error())
						continue
					}

					services <- srv
				}
			}
		}
	}()

	go storage.Service().Watch(context.Get().Background(), services)
	go storage.Pod().Watch(context.Get().Background(), pods)
	go storage.Endpoint().Watch(context.Get().Background(), endpoints)
}
