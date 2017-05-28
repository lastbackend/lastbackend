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

const logLevel = 2

func Get(name string) ([]net.IP, error) {

	var (
		err     error
		log     = context.Get().GetLogger()
		storage = context.Get().GetStorage()
		cache   = context.Get().GetCache()
		data    = []string{}
	)

	log.V(logLevel).Debugf("Endpoint: get endpoint `%s` ip list from cache", name)

	var ips = cache.Endpoints().Get(name)

	if ips != nil && len(ips) != 0 {
		return ips, nil
	}

	if len(data) == 0 {
		log.V(logLevel).Debugf("Endpoint: try find endpoint `%s` in storage", name)

		result, err := storage.Endpoint().Get(context.Get().Background(), name)
		if err != nil {
			log.V(logLevel).Errorf("Endpoint: get endpoint `%s` from storage err: %s", name, err.Error())
			return nil, err
		}

		for _, ip := range result {
			data = append(data, ip)
		}
	}

	data = util.RemoveDuplicates(data)

	ips, err = util.ConvertStringIPToNetIP(data)
	if err != nil {
		log.Errorf("Endpoint: convert endpoint `%s` ips to net ips err: %s", name, err.Error())
		return ips, err
	}

	if err := cache.Endpoints().Set(name, ips); err != nil {
		log.V(logLevel).Errorf("Endpoint: set endpoint `%s` to cache err: %s", name, err.Error())
		return nil, err
	}

	log.V(logLevel).Debugf("Endpoint: get ip list from cache for `%s` successfully: %v", name, data)

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

	log.V(logLevel).Debugf("Endpoint: update name `%s` in cache", name)

	result, err := storage.Endpoint().Get(context.Get().Background(), name)
	if err != nil {
		log.V(logLevel).Errorf("Endpoint: get endpoint `%s` from storage err: %s", name, err.Error())
		return err
	}

	for _, ip := range result {
		data = append(data, ip)
	}

	data = util.RemoveDuplicates(data)

	ips, err := util.ConvertStringIPToNetIP(data)
	if err != nil {
		log.V(logLevel).Errorf("Endpoint: convert endpoint `%s` ip tp net.IP err: %s", name, err.Error())
		return err
	}

	if err := cache.Endpoints().Set(name, ips); err != nil {
		log.V(logLevel).Errorf("Endpoint: set endpoint `%s` to cache err: %s", name, err.Error())
		return err
	}

	log.V(logLevel).Debugf("Endpoint: update name `%s` in cache successfully ", name)

	return nil
}

func Remove(name string) error {

	var (
		err   error
		log   = context.Get().GetLogger()
		cache = context.Get().GetCache()
	)

	log.V(logLevel).Debugf("Endpoint: remove name `%s` from cache", name)

	err = cache.Endpoints().Del(name)
	if err != nil {
		log.V(logLevel).Errorf("Endpoint: delete endpoint `%s` from cache err: %s", name, err.Error())
		return err
	}

	log.V(logLevel).Debugf("Endpoint: remove name `%s` from cache successfully", name)

	return nil
}
