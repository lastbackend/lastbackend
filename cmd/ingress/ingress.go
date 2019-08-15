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
	"github.com/lastbackend/lastbackend/pkg/ingress"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
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
		{Name: "access-token", Short: "", Value: "", Desc: "Access token to API server", Bind: "token"},
		{Name: "haproxy-config-path", Short: "", Value: "/var/run/lastbackend/ingress/haproxy", Desc: "HAProxy configuration path setup", Bind: "haproxy.config"},
		{Name: "haproxy-pid", Short: "", Value: "/var/run/lastbackend/ingress/haproxy/haproxy.pid", Desc: "HAProxy pid file path", Bind: "haproxy.pid"},
		{Name: "haproxy-exec", Short: "", Value: "/usr/sbin/haproxy", Desc: "HAProxy entrypoint path", Bind: "haproxy.exec"},
		{Name: "haproxy-stat-port", Short: "", Value: 1936, Desc: "HAProxy statistic port definition. If not provided - statistic will be disabled", Bind: "haproxy.stat.port"},
		{Name: "haproxy-stat-username", Short: "", Value: "", Desc: "HAProxy statistic access username", Bind: "haproxy.stat.username"},
		{Name: "haproxy-stat-password", Short: "", Value: "", Desc: "HAProxy statistic access password", Bind: "haproxy.stat.password"},
		{Name: "bind-interface", Short: "", Value: "eth0", Desc: "Exporter bind network interface", Bind: "network.interface"},
		{Name: "network-proxy", Short: "", Value: "ipvs", Desc: "Container proxy interface driver", Bind: "network.cpi.type"},
		{Name: "network-proxy-iface-internal", Short: "", Value: "docker0", Desc: "Network external interface binding", Bind: "network.cpi.interface.internal"},
		{Name: "network-proxy-iface-external", Short: "", Value: "eth0", Desc: "Network container bridge binding", Bind: "network.cpi.interface.external"},
		{Name: "network-driver", Short: "", Value: "vxlan", Desc: "Container overlay network driver", Bind: "network.cni.type"},
		{Name: "network-driver-iface-external", Short: "", Value: "eth0", Desc: "", Bind: "network.cni.interface.external"},
		{Name: "network-driver-iface-bridge", Short: "", Value: "docker0", Desc: "", Bind: "network.cni.interface.internal"},
		{Name: "network-resolvers", Short: "", Value: []string{"8.8.8.8", "8.8.4.4"}, Desc: "Additional resolvers IPS for Ingress", Bind: "resolver.servers"},
		{Name: "api-uri", Short: "", Value: "", Desc: "REST API endpoint", Bind: "api.uri"},
		{Name: "api-cert-file", Short: "", Value: "", Desc: "REST API TLS certificate file path", Bind: "api.tls.cert"},
		{Name: "api-private-key-file", Short: "", Value: "", Desc: "REST API TLS private key file path", Bind: "api.tls.key"},
		{Name: "api-ca-file", Short: "", Value: "", Desc: "REST API TSL certificate authority file path", Bind: "api.tls.ca"},
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

	ingress.Daemon(v)
}
