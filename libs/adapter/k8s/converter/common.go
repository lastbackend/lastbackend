package converter

import (
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/1.5/pkg/util/intstr"
)

func Set_defaults_extensions_deployment(obj *extensions.Deployment) {

	// Default labels and selector to labels from pod template spec.
	var (
		replicas       int32 = 1
		labels               = obj.Spec.Template.Labels
		maxUnavailable       = intstr.FromInt(1)
		maxSurge             = intstr.FromInt(1)
	)

	obj.APIVersion = "extensions/v1beta1"
	obj.Kind = "Deployment"

	if labels != nil {
		if obj.Spec.Selector == nil {
			obj.Spec.Selector = &unversioned.LabelSelector{MatchLabels: labels}
		}

		if len(obj.Labels) == 0 {
			obj.Labels = labels
		}
	}

	// Set DeploymentSpec.Replicas to 1 if it is not set.
	obj.Spec.Replicas = replicas

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
		strategy.RollingUpdate.MaxUnavailable = maxUnavailable

		// Set default MaxSurge as 1 by default.
		strategy.RollingUpdate.MaxSurge = maxSurge
	}
}

func Set_defaults_v1beta1_deployment(obj *v1beta1.Deployment) {

	// Default labels and selector to labels from pod template spec.
	var (
		replicas       int32 = 1
		labels               = obj.Spec.Template.Labels
		maxUnavailable       = intstr.FromInt(1)
		maxSurge             = intstr.FromInt(1)
	)

	obj.APIVersion = "extensions/v1beta1"
	obj.Kind = "Deployment"

	if labels != nil {
		if obj.Spec.Selector == nil {
			obj.Spec.Selector = &v1beta1.LabelSelector{MatchLabels: labels}
		}
		if len(obj.Labels) == 0 {
			obj.Labels = labels
		}
	}

	// Set DeploymentSpec.Replicas to 1 if it is not set.
	obj.Spec.Replicas = &replicas

	strategy := &obj.Spec.Strategy

	// Set default DeploymentStrategyType as RollingUpdate.
	if strategy.Type == "" {
		strategy.Type = v1beta1.RollingUpdateDeploymentStrategyType
	}

	if strategy.Type == v1beta1.RollingUpdateDeploymentStrategyType {
		if strategy.RollingUpdate == nil {
			rollingUpdate := v1beta1.RollingUpdateDeployment{}
			strategy.RollingUpdate = &rollingUpdate
		}

		// Set default MaxUnavailable as 1 by default.
		strategy.RollingUpdate.MaxUnavailable = &maxUnavailable

		// Set default MaxSurge as 1 by default.
		strategy.RollingUpdate.MaxSurge = &maxSurge
	}
}

func Set_defaults_v1_Pod(obj *v1.Pod) {

	for i := range obj.Spec.Containers {
		// set requests to limits if requests are not specified, but limits are
		if obj.Spec.Containers[i].Resources.Limits != nil {
			if obj.Spec.Containers[i].Resources.Requests == nil {
				obj.Spec.Containers[i].Resources.Requests = make(v1.ResourceList)
			}

			for key, value := range obj.Spec.Containers[i].Resources.Limits {
				if _, exists := obj.Spec.Containers[i].Resources.Requests[key]; !exists {
					obj.Spec.Containers[i].Resources.Requests[key] = *(value.Copy())
				}
			}
		}
	}
}
