package converter

import (
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/apis/extensions"
	"k8s.io/client-go/1.5/pkg/util/intstr"
)

func Set_defaults_extensions_deployment(obj *extensions.Deployment) {
	// Default labels and selector to labels from pod template spec.
	labels := obj.Spec.Template.Labels

	if labels != nil {
		if obj.Spec.Selector == nil {
			obj.Spec.Selector = &unversioned.LabelSelector{MatchLabels: labels}
		}
		if len(obj.Labels) == 0 {
			obj.Labels = labels
		}
	}

	// Set DeploymentSpec.Replicas to 1 if it is not set.
	obj.Spec.Replicas = 1

	strategy := &obj.Spec.Strategy

	// Set default DeploymentStrategyType as RollingUpdate.
	if strategy.Type == "" {
		strategy.Type = extensions.RollingUpdateDeploymentStrategyType
	}

	if strategy.Type == extensions.RollingUpdateDeploymentStrategyType {
		if strategy.RollingUpdate == nil {
			rollingUpdate := extensions.RollingUpdateDeployment{}
			strategy.RollingUpdate = &rollingUpdate
		}

		// Set default MaxUnavailable as 1 by default.
		strategy.RollingUpdate.MaxUnavailable = intstr.FromInt(1)

		// Set default MaxSurge as 1 by default.
		strategy.RollingUpdate.MaxSurge = intstr.FromInt(1)
	}
}
