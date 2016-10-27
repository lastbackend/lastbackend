package k8s

import (
	"github.com/lastbackend/lastbackend/libs/adapter/k8s"
	"k8s.io/client-go/1.5/kubernetes"
)

type IK8S interface {
	kubernetes.Interface
	k8s.LBClientsetInterface
}
