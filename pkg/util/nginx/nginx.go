//
// Last.Backend LLC CONFidENTIAL
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

package nginx

import (
	"os"
	"path/filepath"
	"text/template"
)

var TplNginxConf = template.Must(template.New("").Parse(`
{{range $upstream := .Upstreams}}
upstream {{$upstream.Name}} {
	server {{$upstream.Address}};
}
{{end}}
server {
	{{if eq .Server.Protocol "http"}}
	listen {{.Server.Port}};
	{{else if eq .Server.Protocol "https"}}
	listen {{.Server.Port}} ssl;{{end}}
	server_name	{{.Server.Hostname}};
	{{if eq .Server.Protocol "https"}}
	ssl	on;
	ssl_certificate	{{.RootPath}}/ssl/server.crt;
	ssl_certificate_key	{{.RootPath}}/ssl/server.key;
	{{end}}
{{range $location := .Server.Locations}}
	location {{$location.Path}} {
		proxy_set_header	Host $host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
		proxy_pass		{{$location.ProxyPass}};
		proxy_http_version	1.1;
		proxy_set_header	Upgrade $http_upgrade;
		proxy_set_header	Connection 'upgrade';
	}
{{end}}
}
`))

type Nginx struct{}

func (n Nginx) GenerateConfig(path string, template interface{}) error {
	f, err := os.Create(filepath.Join(path))
	if err != nil {
		return err
	}
	defer f.Close()

	return TplNginxConf.Execute(f, template)
}

func (n Nginx) RemoveConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil
		}
	}
	return os.Remove(path)
}
