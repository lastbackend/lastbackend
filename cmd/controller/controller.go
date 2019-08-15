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

// Last.Backend Open-source API
//
// Open-source system for automating deployment, scaling, and management of containerized applications.
//
// Terms Of Applications:
//
// https://lastbackend.com/legal/terms/
//
//     Schemes: https
//     Host: api.lastbackend.com
//     BasePath: /
//     Version: 0.9.4
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Last.Backend Teams <team@lastbackend.com> https://lastbackend.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearerToken:
//
//     SecurityDefinitions:
//       bearerToken:
//         description: Bearer Token authentication
//         type: apiKey
//         name: authorization
//         in: header
//
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
//
// swagger:meta
package main

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/controller"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const default_env_prefix = "LB"
const default_config_type = "yaml"
const default_config_name = "config"

var (
	flags = []struct {
		// flag name
		Name string
		// flag short name
		Short string
		// flag value
		Value interface{}
		// flag description
		Desc string
		// viper name for binding from flag
		Bind string
	}{
		{Name: "services-cidr", Short: "", Value: "172.0.0.0/24", Desc: "Services IP CIDR for internal IPAM service", Bind: "service.cidr"},
		{Name: "storage", Short: "", Value: "etcd", Desc: "Set storage driver (Allow: etcd, mock)", Bind: "storage.driver"},
		{Name: "etcd-cert-file", Short: "", Value: "", Desc: "ETCD database cert file path", Bind: "storage.etcd.tls.cert"},
		{Name: "etcd-private-key-file", Short: "", Value: "", Desc: "ETCD database private key file path", Bind: "storage.etcd.tls.key"},
		{Name: "etcd-ca-file", Short: "", Value: "", Desc: "ETCD database certificate authority file", Bind: "storage.etcd.tls.ca"},
		{Name: "etcd-endpoints", Short: "", Value: []string{"127.0.0.1:2379"}, Desc: "ETCD database endpoints list", Bind: "storage.etcd.endpoints"},
		{Name: "etcd-prefix", Short: "", Value: "lastbackend", Desc: "ETCD database storage prefix", Bind: "storage.etcd.prefix"},
		{Name: "verbose", Short: "v", Value: 0, Desc: "Set log level from 0 to 7", Bind: "verbose"},
		{Name: "config", Short: "c", Value: "", Desc: "Path for the configuration file", Bind: "config"},
	}
)

func main() {

	for _, item := range flags {
		switch item.Value.(type) {
		case string:
			flag.StringP(item.Name, item.Short, item.Value.(string), item.Desc)
		case int:
			flag.IntP(item.Name, item.Short, item.Value.(int), item.Desc)
		case []string:
			flag.StringSliceP(item.Name, item.Short, item.Value.([]string), item.Desc)
		default:
			panic(fmt.Sprintf("bad %s argument value", item.Name))
		}
	}

	flag.Parse()

	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetEnvPrefix(default_env_prefix)

	for _, item := range flags {

		if len(flag.Lookup(item.Name).Value.String()) != 0 {
			if err := v.BindPFlag(item.Bind, flag.Lookup(item.Name)); err != nil {
				panic(err)
			}
		} else {
			v.SetDefault(item.Bind, nil)
		}

		name := strings.Replace(strings.ToUpper(item.Name), "-", "_", -1)
		name = strings.Join([]string{default_env_prefix, name}, "_")

		if err := v.BindEnv(item.Bind, name); err != nil {
			panic(err)
		}

		if !validator.IsZeroOfUnderlyingType(item.Value) {
			v.SetDefault(item.Bind, item.Value)
		}

	}

	v.SetConfigType(default_config_type)
	v.SetConfigFile(v.GetString(default_config_name))

	if len(v.GetString("config")) != 0 {
		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Sprintf("Read config err: %v", err))
		}
	}

	controller.Daemon(v)
}
