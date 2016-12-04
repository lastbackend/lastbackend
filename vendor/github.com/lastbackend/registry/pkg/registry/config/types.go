package config

// The structure of the config to run the daemon
type Config struct {
	Debug bool `yaml:"debug" json:"debug"`

	HttpServer struct {
		Port int `yaml:"port" json:"port"`
	} `yaml:"http_server" json:"http_server"`
}
