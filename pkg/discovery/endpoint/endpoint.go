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

package endpoint

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/discovery/envs"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util"
	"net"
)

const logLevel = 2

func Get(name string) ([]net.IP, error) {

	var (
		err   error
		stg   = envs.Get().GetStorage()
		cache = envs.Get().GetCache()
		data  = []string{}
	)

	log.Debugf(`Endpoint: Get ip list from cache %s`, name)

	var ips = cache.Endpoints().Get(name)

	if ips != nil && len(ips) != 0 {
		return ips, nil
	}

	if len(data) == 0 {
		log.Debugf(`Endpoint: Try find to db %s`, name)

		result, err := stg.Endpoint().Get(context.Background(), name)
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

	cache.Endpoints().Set(name, ips)

	log.Debugf(`Endpoint: Get ip list from cache for %s successfully: %v`, name, data)

	return ips, nil
}

func GetForService(name string) ([]net.IP, error) {

	var (
		err error
		//storage = envs.Get().GetStorage()
		cache = envs.Get().GetCache()
		data  = make([]string, 0)
	)

	log.V(logLevel).Debugf("Endpoint: get endpoint `%s` ip list from cache", name)

	var ips = cache.Endpoints().Get(name)

	if ips != nil && len(ips) != 0 {
		return ips, nil
	}

	if len(data) == 0 {
		log.V(logLevel).Debugf("Endpoint: try find endpoint `%s` in storage", name)

		//namespace := ""
		//service   := ""
		//
		//result, err := storage.Pod().ListByService(context.Background(), namespace, service)
		//if err != nil {
		//	log.V(logLevel).Errorf("Endpoint: get endpoint `%s` from storage err: %s", name, err)
		//	return nil, err
		//}

		for _, ip := range ips {
			data = append(data, ip.String())
		}
	}

	data = util.RemoveDuplicates(data)

	ips, err = util.ConvertStringIPToNetIP(data)
	if err != nil {
		log.Errorf("Endpoint: convert endpoint `%s` ips to net ips err: %s", name, err)
		return ips, err
	}

	// TODO: Update cache and add expire
	//if err := cache.Endpoint().Set(name, ips); err != nil {
	//	log.V(logLevel).Errorf("Endpoint: set endpoint `%s` to cache err: %s", name, err)
	//	return nil, err
	//}

	log.V(logLevel).Debugf("Endpoint: get ip list from cache for `%s` successfully: %v", name, data)

	return ips, nil
}

func GetForRoute(name string) ([]net.IP, error) {

	var (
		err error
		//storage = envs.Get().GetStorage()
		cache = envs.Get().GetCache()
		data  = make([]string, 0)
	)

	log.V(logLevel).Debugf("Endpoint: get endpoint `%s` ip list from cache", name)

	var ips = cache.Endpoints().Get(name)

	if ips != nil && len(ips) != 0 {
		return make([]net.IP, 0), nil
	}

	if len(data) == 0 {
		log.V(logLevel).Debugf("Endpoint: try find endpoint `%s` in storage", name)

		//namespace := ""
		//
		//result, err := storage.Route().Get(context.Background(), namespace, name)
		//if err != nil {
		//	log.V(logLevel).Errorf("Endpoint: get endpoint `%s` from storage err: %s", name, err)
		//	return nil, err
		//}

		for _, ip := range ips {
			data = append(data, ip.String())
		}
	}

	data = util.RemoveDuplicates(data)

	ips, err = util.ConvertStringIPToNetIP(data)
	if err != nil {
		log.Errorf("Endpoint: convert endpoint `%s` ips to net ips err: %s", name, err)
		return ips, err
	}

	// TODO: Update cache and add expire
	//if err := cache.Endpoint().Set(name, ips); err != nil {
	//	log.V(logLevel).Errorf("Endpoint: set endpoint `%s` to cache err: %s", name, err)
	//	return nil, err
	//}

	log.V(logLevel).Debugf("Endpoint: get ip list from cache for `%s` successfully: %v", name, data)

	return ips, nil
}

func Update(name string) error {

	var (
		err   error
		stg   = envs.Get().GetStorage()
		cache = envs.Get().GetCache()
		data  = []string{}
	)

	log.Debugf(`Endpoint: Update name %s in cache `, name)

	result, err := stg.Endpoint().Get(context.Background(), name)
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

	cache.Endpoints().Set(name, ips)

	log.Debugf(`Endpoint: Update name %s in cache successfully `, name)

	return nil
}

func Remove(name string) error {

	var (
		err   error
		cache = envs.Get().GetCache()
	)

	log.Debugf(`Endpoint: Remove name %s from cache `, name)

	err = cache.Endpoints().Del(name)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf(`Endpoint: Remove name %s from cache successfully `, name)

	return nil
}
