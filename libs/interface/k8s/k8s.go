package k8s

import (
	"k8s.io/client-go/1.5/kubernetes"
)

type IK8S interface {
	kubernetes.Interface
}
