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
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"github.com/lastbackend/lastbackend/pkg/util"
	"net"
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

	var ips = cache.EndpointCache.Get(name)

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

	cache.EndpointCache.Set(name, ips)

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

	cache.EndpointCache.Set(name, ips)

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

	err = cache.EndpointCache.Del(name)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf(`Endpoint: Remove name %s from cache successfully `, name)

	return nil
}
