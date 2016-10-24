package k8s

import (
	"k8s.io/client-go/1.5/kubernetes"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s"
)

type IK8S interface {
	kubernetes.Interface
	k8s.LBClientsetInterface
}
