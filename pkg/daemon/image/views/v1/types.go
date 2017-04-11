package v1

type ImageSpec struct {
	// Image full name
	Name string `json:"name"`
	// Image pull provision flag
	Pull bool `json:"pull"`
	// Image Auth base64 encoded string
	Auth string `json:"auth"`
}
