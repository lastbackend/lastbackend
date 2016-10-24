package model

type ServiceSource struct {
	UUID      string `json:"uuid,omitempty"`
	ServiceID string `json:"service,omitempty"`
	Type      string `json:"type"`
	Hub       string `json:"hub"`
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
	Branch    string `json:"branch"`
	Readme    string `json:"readme"`
}

type ServiceSources []ServiceSource
