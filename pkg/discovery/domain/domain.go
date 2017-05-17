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
	"github.com/lastbackend/lastbackend/pkg/util"
	"net"
)

func Update(domain string) error {

	var (
		err     error
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage()
		cache   = context.Get().GetCache()
		data    = []string{}
	)

	log.Debug(`Update domain in cache `, domain)

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
	cache.Insert(domain, ips)
	cache.Unlock()

	log.Debug(`Update domain in cache successfully `, domain)

	return nil
}

func Remove(domain string) error {

	var (
		err   error
		log   = context.Get().GetLogger()
		cache = context.Get().GetCache()
	)

	log.Debug(`Remove domain from cache `, domain)

	err = cache.Remove(domain)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug(`Remove domain from cache successfully `, domain)

	return nil
}

func GetIPList(domain string) ([]net.IP, error) {

	var (
		err     error
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage()
		cache   = context.Get().GetCache()
		data    = []string{}
	)

	log.Debug(`Get ip list from cache `, domain)

	var ips = cache.IPList(domain)

	if ips != nil && len(ips) != 0 {
		return ips, nil
	}

	if len(data) == 0 {
		log.Debug(`Try find to db `, domain)

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
		log.Error(err)
		return nil, err
	}

	log.Debug(`ips`, ips)

	cache.Lock()
	cache.Insert(domain, ips)
	cache.Unlock()

	log.Debug(`Get ip list from cache successfully `, domain, data)

	return ips, nil
}
