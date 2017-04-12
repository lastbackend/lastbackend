package v1

type Meta struct {
	Hostname     string `json:"hostname"`
	OSName       string `json:"os_name"`
	OSType       string `json:"os_type"`
	Architecture string `json:"architecture"`

	CRI     CRIMeta     `json:"cri"`
	CPU     HostCPU     `json:"cpu"`
	Memory  HostMemory  `json:"memory"`
	Network HostNetwork `json:"network"`
	Storage HostStorage `json:"storage"`
}

type CRIMeta struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type HostCPU struct {
	Name  string `json:"name"`
	Cores int64  `json:"cores"`
}

type HostMemory struct {
	Total     int64 `json:"total"`
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
}

type HostNetwork struct {
	Interface string   `json:"interface,omitempty"`
	IP        []string `json:"ip,omitempty"`
}

type HostStorage struct {
	Available string `json:"available"`
	Used      string `json:"used"`
	Total     string `json:"total"`
}
