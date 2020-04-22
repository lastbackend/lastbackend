package templates

import (
	"bytes"
	"text/template"
)

type ContainerdConfig struct {
	IsRunningInUserNS bool
	SELinuxEnabled    bool
	Opt               string
	RegistryConfig    *Registry
}

const ContainerdConfigTemplate = `
[plugins.opt] 
	# The containerd managed opt directory provides a way for users to install containerd 
	# dependencies using the existing distribution infrastructure.
	# default: /opt/containerd
  path = "{{ .Opt }}"

[plugins.cri]
	# stream_server_address is the ip address streaming server is listening on.
  stream_server_address = "127.0.0.1"

	# stream_server_port is the port streaming server is listening on.
  stream_server_port = "10001"

  # enable_selinux indicates to enable the selinux support.
  enable_selinux = {{ .SELinuxEnabled }}

{{- if .IsRunningInUserNS }}
  # disable_cgroup indicates to disable the cgroup support.
  # This is useful when the daemon does not have permission to access cgroup.
  disable_cgroup = true

  # disable_apparmor indicates to disable the apparmor support.
  # This is useful when the daemon does not have permission to access apparmor.
  disable_apparmor = true

  # restrict_oom_score_adj indicates to limit the lower bound of OOMScoreAdj to
  # the containerd's current OOMScoreAdj.
  # This is useful when the containerd does not have permission to decrease OOMScoreAdj.
  restrict_oom_score_adj = true
{{end}}

# is a map from CRI RuntimeHandler strings, which specify types
# of runtime configurations, to the matching configurations.
# In this example, 'runc' is the RuntimeHandler string to match.
[plugins.cri.containerd.runtimes.runc]
	# runtime_type is the runtime type to use in containerd.
	# The default value is "io.containerd.runc.v2" since containerd 1.4.
	# The default value was "io.containerd.runc.v1" in containerd 1.3, "io.containerd.runtime.v1.linux" in prior releases.
	runtime_type = "io.containerd.runc.v2"

# contains config related to the registry
{{ if .RegistryConfig }}
{{ if .RegistryConfig.Mirrors }}
[plugins.cri.registry.mirrors]{{end}}
{{range $k, $v := .RegistryConfig.Mirrors }}
[plugins.cri.registry.mirrors."{{$k}}"]
  endpoint = [{{range $i, $j := $v.Endpoints}}{{if $i}}, {{end}}{{printf "%q" .}}{{end}}]
{{end}}

{{range $k, $v := .RegistryConfig.Configs }}
{{ if $v.Auth }}
[plugins.cri.registry.configs."{{$k}}".auth]
  {{ if $v.Auth.Username }}username = "{{ $v.Auth.Username }}"{{end}}
  {{ if $v.Auth.Password }}password = "{{ $v.Auth.Password }}"{{end}}
  {{ if $v.Auth.Auth }}auth = "{{ $v.Auth.Auth }}"{{end}}
  {{ if $v.Auth.IdentityToken }}identitytoken = "{{ $v.Auth.IdentityToken }}"{{end}}
{{end}}
{{ if $v.TLS }}
[plugins.cri.registry.configs."{{$k}}".tls]
  {{ if $v.TLS.CAFile }}ca_file = "{{ $v.TLS.CAFile }}"{{end}}
  {{ if $v.TLS.CertFile }}cert_file = "{{ $v.TLS.CertFile }}"{{end}}
  {{ if $v.TLS.KeyFile }}key_file = "{{ $v.TLS.KeyFile }}"{{end}}
{{end}}
{{end}}
{{end}}
`

func ParseTemplateFromConfig(templateBuffer string, config interface{}) (string, error) {
	out := new(bytes.Buffer)
	t := template.Must(template.New("compiled_template").Parse(templateBuffer))
	if err := t.Execute(out, config); err != nil {
		return "", err
	}
	return out.String(), nil
}
