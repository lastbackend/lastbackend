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

package runtime

const HaproxyTemplate = `
#---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
  log         127.0.0.1 local2

  maxconn     2048
  daemon

  # turn on stats unix socket
  # stats socket /var/lib/haproxy/stats

  # ssl settings, as we want to get pretty result
  # @ https://www.ssllabs.com/ssltest
  tune.ssl.default-dh-param 2048
  ssl-default-bind-options no-sslv3 no-tls-tickets
  ssl-default-bind-ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-DSS-AES128-GCM-SHA256:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-DSS-AES128-SHA256:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:DHE-RSA-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:AES:CAMELLIA:DES-CBC3-SHA:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!PSK:!aECDH:!EDH-DSS-DES-CBC3-SHA:!EDH-RSA-DES-CBC3-SHA:!KRB5-DES-CBC3-SHA

#---------------------------------------------------------------------
# common defaults that all the 'listen' and 'backend' sections will
# use if not designated in their block
#---------------------------------------------------------------------
defaults
  mode                    tcp
  log                     global
  option                  dontlog-normal
  option                  tcpka
  retries                 3
  timeout http-request    10s
  timeout queue           1m
  timeout connect         10s
  timeout client          1m
  timeout server          1m
  timeout http-keep-alive 10s
  timeout check           5s
  default-server init-addr none

resolvers lstbknd
  {{range $endpoint, $port := .Resolvers}}nameserver dns {{$endpoint}}:{{$port}}
  {{end}}

#---------------------------------------------------------------------
# frontend which proxys stats
#---------------------------------------------------------------------
#listen stats # Define a listen section called "stats"
#  bind :92967 # Listen on localhost:9000
#  mode http
#  stats enable  # Enable stats page
#  stats hide-version  # Hide HAProxy version
#  stats realm Haproxy\ Statistics  # Title text for popup window
#  stats uri /stats  # Stats URI
  

#---------------------------------------------------------------------
# frontend which proxys raw/ssl request to the backends
#---------------------------------------------------------------------
{{range $port, $f := .Frontend}}{{if eq $f.Type "http" }}
frontend http
  mode http
  bind :80
  http-send-name-header Host
  http-request set-header Host %[req.hdr(Host)]
  http-request set-header X-Forwarded-Host %[req.hdr(Host)]

  {{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}acl r_{{$backend}}  hdr_dom(host) -i {{$domain}}  path_beg {{$path}}
  {{end}}{{end}}
  {{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}acl r_{{$backend}}_down  nbsrv({{$backend}}) lt 1
  {{end}}{{end}}
  {{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}use_backend local_http if r_{{$backend}}_down r_{{$backend}}
  {{end}}{{end}}
  {{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}use_backend {{$backend}} if r_{{$backend}} {{if eq $path "/"}}{{range $p, $b := $acl}}{{if ne $p "/"}}!r_{{$b}} {{end}}{{end}}{{end}}
  {{end}}{{end}}
  default_backend local_http
{{else if eq $f.Type "https" }}
frontend https
  bind :443
  option socket-stats
  tcp-request inspect-delay 5s
  tcp-request content accept if { req_ssl_hello_type 1 }
  {{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}acl r_{{$backend}}  hdr_dom(host) -i {{$domain}}  path_beg {{$path}}
  {{end}}{{end}}
	{{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}acl r_{{$backend}}_down  nbsrv({{$backend}}) lt 1
  {{end}}{{end}}
  {{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}use_backend local_http if r_{{$backend}}_down r_{{$backend}}
  {{end}}{{end}}
  {{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}use_backend {{$backend}} if r_{{$backend}}
  {{end}}{{end}}
  default_backend local_http
{{else if eq $f.Type "tcp" }}
frontend {{$port}}_tcp
  bind 0.0.0.0:{{$port}}
  {{range $domain, $acl := .Rules}}{{range $path, $backend := $acl}}use_backend {{$backend}} if r_{{$backend}}
  {{end}}{{end}}
  default_backend local_http
{{end}}{{end}}

#---------------------------------------------------------------------
# local proxy configuration
#---------------------------------------------------------------------
backend local_http
    mode http
    errorfile 503 /var/run/html/errors/503.html

#---------------------------------------------------------------------
# balancing between the various backends
#---------------------------------------------------------------------
{{range $name, $b := .Backend}}{{if eq $b.Type "http" }}
backend {{$name}}
  mode http
  balance roundrobin
  option forwardfor
  server {{$b.Endpoint}} {{$b.Endpoint}}:{{$b.Port}} check init-addr last,libc,none resolvers lstbknd
{{else if eq $b.Type "https" }}
backend {{$name}}
  mode tcp
  # maximum SSL session ID length is 32 bytes.
  stick-table type binary len 32 size 30k expire 30m
  acl clienthello req_ssl_hello_type 1
  acl serverhello rep_ssl_hello_type 2
  # use tcp content accepts to detects ssl client and server hello.
  tcp-request inspect-delay 5s
  tcp-request content accept if clienthello
  # no timeout on response inspect delay by default.
  tcp-response content accept if serverhello
  stick on payload_lv(43,1) if clienthello
  # Learn on response if server hello.
  stick store-response payload_lv(43,1) if serverhello
  option ssl-hello-chk
  http-request set-header Host {{$b.Domain}}
  server {{$b.Endpoint}} {{$b.Endpoint}}:{{$b.Port}} check init-addr last,libc,none resolvers lstbknd
{{else if eq $b.Type "tcp" }}
backend {{$name}}
  server {{$b.Endpoint}} {{$b.Endpoint}}:{{$b.Port}} check init-addr last,libc,none resolvers lstbknd
{{end}}{{end}}
`
