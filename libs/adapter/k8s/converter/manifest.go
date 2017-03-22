//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package converter

import (
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/common"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func Convert_ReplicationController_to_Deployment(config *v1.ReplicationController) *v1beta1.Deployment {

	var deployment = new(v1beta1.Deployment)

	common.Set_defaults_v1beta1_deployment(deployment)

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

	common.Set_defaults_v1beta1_deployment(deployment)

	deployment.ObjectMeta = config.ObjectMeta
	deployment.Spec.Replicas = &replicas
	deployment.Spec.Template.Spec = config.Spec

	return deployment
}
