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

package config

// The structure of the config to run the daemon
type Config struct {
	Debug bool `yaml:"debug"`

	TokenSecret string `yaml:"secret"`

	TemplateRegistry struct {
		Host string `yaml:"host"`
	} `yaml:"template_registry"`

	ProxyServer struct {
		Port int `yaml:"port"`
	} `yaml:"proxy_server"`

	HttpServer struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"http_server"`

	Etcd struct {
		Endpoints []string `yaml:"endpoints"`
		TLS       struct {
			Key  string `yaml:"key"`
			Cert string `yaml:"cert"`
			CA   string `yaml:"ca"`
		} `yaml:"tls"`
		Quorum bool `yaml:"quorum"`
	} `yaml:"etcd"`

	Registry struct {
		Server   string `yaml:"server"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"registry"`

	VCS struct {
		Github struct {
			Client struct {
				ID       string `yaml:"id"`
				SecretID string `yaml:"secret_id"`
			} `yaml:"client"`
		} `yaml:"github"`

		Bitbucket struct {
			Client struct {
				ID       string `yaml:"id"`
				SecretID string `yaml:"secret_id"`
			} `yaml:"client"`
			RedirectUri string `yaml:"redirect_uri"`
		} `yaml:"bitbucket"`

		Gitlab struct {
			Client struct {
				ID       string `yaml:"id"`
				SecretID string `yaml:"secret_id"`
			} `yaml:"client"`
			RedirectUri string `yaml:"redirect_uri"`
		} `yaml:"gitlab"`
	} `yaml:"vcs"`
}
