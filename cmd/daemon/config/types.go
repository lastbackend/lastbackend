package config

type Config struct {
	Debug bool `yaml:"debug"`

	HttpServer struct {
		Port int `yaml:"port"`
	} `yaml:"http_server"`

	K8S struct {
		Host string `yaml:"host"`
		SSL  struct {
			CA   string `yaml:"ca"`
			Key  string `yaml:"key"`
			Cert string `yaml:"cert"`
		} `yaml:"ssl"`
	} `yaml:"k8s"`
}
