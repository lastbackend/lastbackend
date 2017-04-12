package v1

import (
	"github.com/lastbackend/lastbackend/pkg/daemon/pod/views/v1"
)

type Spec struct {
	// Pods spec
	Pods []v1.Pod `json:"pods"`
}
