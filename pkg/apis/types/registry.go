package types

type Registry struct {
	// Registry ID
	ID string
	// Registry name
	Name string
	// Registry owner
	Owner string
	// Registry hub in http(s)://host:port format
	Hub string
	// Registry authentication information
	Auth RegistryAuth
}

type RegistryAuth struct {
	// Registry auth username
	Username string
	// Registry auth password
	Password string
	// Registry auth email
	Email string
	// Registry host
	Host string
}
