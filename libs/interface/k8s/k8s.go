package k8s

import (
	"k8s.io/client-go/kubernetes"
)

type IK8S interface {
	kubernetes.Interface
}
