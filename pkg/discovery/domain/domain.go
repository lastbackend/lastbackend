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

package domain

import (
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util"
	"net"
)

func Get(domain string) ([]net.IP, error) {

	var (
		err     error
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage()
		cache   = context.Get().GetCache()
		data    = []string{}
	)

	log.Debugf(`Domain: Get ip list from cache %s`, domain)

	var ips = cache.Get(domain)

	if ips != nil && len(ips) != 0 {
		return ips, nil
	}

	if len(data) == 0 {
		log.Debugf(`Domain: Try find to db %s`, domain)

		result, err := storage.Endpoint().Get(context.Get().Background(), domain)
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
		log.Errorf("Domain: convert ips to net ips error %s", err.Error())
		return ips, err
	}

	log.Debug(`ips`, ips)

	cache.Lock()
	cache.Set(domain, ips)
	cache.Unlock()

	log.Debugf(`Domain: Get ip list from cache for %s successfully: %v`, domain, data)

	return ips, nil
}

func Update(domain string) error {

	var (
		err     error
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage()
		cache   = context.Get().GetCache()
		data    = []string{}
	)

	log.Debugf(`Domain: Update domain %s in cache `, domain)

	result, err := storage.Endpoint().Get(context.Get().Background(), domain)
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
	cache.Set(domain, ips)
	cache.Unlock()

	log.Debugf(`Domain: Update domain %s in cache successfully `, domain)

	return nil
}

func Remove(domain string) error {

	var (
		err   error
		log   = context.Get().GetLogger()
		cache = context.Get().GetCache()
	)

	log.Debugf(`Domain: Remove domain %s from cache `, domain)

	err = cache.Del(domain)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf(`Domain: Remove domain %s from cache successfully `, domain)

	return nil
}

func Watch() error {
	var (
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage().Endpoint()
		cache   = context.Get().GetCache()
		event   = make(chan string, 1)
	)

	log.Debug("Domain: Endpoint: Watch")

	go func() {
		for {
			select {
			case endpoint := <-event:
				{
					i, err := storage.Get(context.Get().Background(), endpoint)
					if err != nil && err.Error() != store.ErrKeyNotFound {
						if err = cache.Del(endpoint); err != nil {
							log.Errorf("Domain: remove ips from cache error %s", err.Error())
						}
						continue
					} else {
						log.Errorf("Domain: get ips for domain error %s", err.Error())
						continue
					}

					ips, err := util.ConvertStringIPToNetIP(i)
					if err != nil {
						log.Errorf("Domain: convert ips to net ips error %s", err.Error())
						continue
					}

					if err = cache.Set(endpoint, ips); err != nil {
						log.Errorf("Domain: save ips to cache error %s", err.Error())
						continue
					}
				}
			}
		}
	}()

	return storage.Watch(context.Get().Background(), event)
}
