package config

type Config struct {
	Debug bool `yaml:"debug"`

	HttpServer struct {
		Port int `yaml:"port"`
	} `yaml:"http_server"`
}
