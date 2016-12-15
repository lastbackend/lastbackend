package converter

import (
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
)

func Convert_ReplicationController_to_Deployment(config *v1.ReplicationController) *v1beta1.Deployment {

	var deployment = new(v1beta1.Deployment)

	Set_defaults_v1beta1_deployment(deployment)

	deployment.ObjectMeta = config.ObjectMeta
	deployment.Spec.Replicas = config.Spec.Replicas
	deployment.Spec.Template.Spec = config.Spec.Template.Spec
	deployment.Spec.Template.ObjectMeta = config.Spec.Template.ObjectMeta

	for key, val := range config.Spec.Selector {
		deployment.Spec.Selector.MatchLabels[key] = val
	}

	return deployment
}

func Convert_Pod_to_Deployment(config *v1.Pod) *v1beta1.Deployment {

	var (
		replicas   int32 = 1
		deployment       = new(v1beta1.Deployment)
	)

	Set_defaults_v1beta1_deployment(deployment)

	deployment.ObjectMeta = config.ObjectMeta
	deployment.Spec.Replicas = &replicas
	deployment.Spec.Template.Spec = config.Spec

	return deployment
}
