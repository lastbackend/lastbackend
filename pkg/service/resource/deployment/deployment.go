package deployment

import (
	"fmt"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/converter"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/pkg/service/resource/common"
	"github.com/lastbackend/lastbackend/pkg/service/resource/pod"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/apis/extensions"
)

// DeploymentDetail is a presentation layer view of Kubernetes Deployment resource.
type Deployment struct {
	ObjectMeta common.ObjectMeta `json:"objectMeta"`
	TypeMeta   common.TypeMeta   `json:"typeMeta"`
	// Detailed information about Pods belonging to this Deployment.
	PodList pod.PodList `json:"podList"`
	// Label selector of the service.
	Selector map[string]string `json:"selector"`
}

// GetDeploymentDetail returns model object of deployment and error, if any.
func GetDeployment(client k8s.IK8S, namespace string, deploymentName string) (*Deployment, error) {

	fmt.Printf("Getting details of %s deployment in %s namespace", deploymentName, namespace)

	deployment, err := client.Extensions().Deployments(namespace).Get(deploymentName)
	if err != nil {
		return nil, err
	}

	// Pods
	podList, err := GetDeploymentPods(client, namespace, deploymentName)
	if err != nil {
		return nil, err
	}

	var meta = new(api.ObjectMeta)

	err = converter.Convert_v1_ObjectMeta_to_api_ObjectMeta(&deployment.ObjectMeta, meta)
	if err != nil {
		return nil, err
	}

	return &Deployment{
		ObjectMeta: common.NewObjectMeta(*meta),
		TypeMeta:   common.NewTypeMeta(common.ResourceKindDeployment),
		PodList:    *podList,
		Selector:   deployment.Spec.Selector.MatchLabels,
	}, nil

}

// getJobPods returns list of pods targeting deployment.
func GetDeploymentPods(client k8s.IK8S, namespace string, deploymentName string) (*pod.PodList, error) {

	deployment, err := client.Extensions().Deployments(namespace).Get(deploymentName)
	if err != nil {
		return nil, err
	}

	var extensionDeployment = new(extensions.Deployment)

	err = converter.Convert_v1beta1_Deployment_To_extensions_Deployment(deployment, extensionDeployment)
	if err != nil {
		return nil, err
	}

	selector, err := unversioned.LabelSelectorAsSelector(extensionDeployment.Spec.Selector)
	if err != nil {
		return nil, err
	}

	options := api.ListOptions{LabelSelector: selector}

	podChannels := common.GetPodListChannelWithOptions(client, common.NewSameNamespaceQuery(namespace), options, 1)

	rawPods := <-podChannels.List
	if err := <-podChannels.Error; err != nil {
		return nil, err
	}

	pods := common.FilterNamespacedPodsBySelector(rawPods.Items, deployment.ObjectMeta.Namespace, deployment.Spec.Selector.MatchLabels)

	podList := pod.CreatePodList(pods)

	return &podList, nil
}
