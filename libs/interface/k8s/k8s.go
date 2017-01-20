package k8s

import (
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/lb"
	"k8s.io/client-go/kubernetes"
)

type IK8S interface {
	kubernetes.Interface
	lb.Interface
}
