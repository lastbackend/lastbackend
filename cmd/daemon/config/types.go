package config

// The structure of the config to run the daemon
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

	RethinkDB struct {
		Address    string   `yaml:"address"`
		Addresses  []string `yaml:"addresses"`
		MaxOpen    int      `yaml:"max_open"`
		InitialCap int      `yaml:"initial_cap"`
		Database   string   `yaml:"database"`
		AuthKey    string   `yaml:"auth_key"`
		SSL        struct {
			CA string `yaml:"ca"`
		} `yaml:"ssl"`
	} `yaml:"rethinkdb"`
}
