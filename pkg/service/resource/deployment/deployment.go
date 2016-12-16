package deployment

import (
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/converter"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/pkg/service/resource/common"
	"github.com/lastbackend/lastbackend/pkg/service/resource/pod"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/apis/extensions"
)

const kind = "deployment"

type Deployment struct {
	ObjectMeta common.ObjectMeta `json:"meta"`
	TypeMeta   common.TypeMeta   `json:"spec"`
	PodList    pod.PodList       `json:"pods"`
	Selector   map[string]string `json:"selector"`
}

func Get(client k8s.IK8S, namespace string, name string) (*Deployment, error) {

	deployment, err := client.Extensions().Deployments(namespace).Get(name)
	if err != nil {
		return nil, err
	}

	var deploymentNew = new(extensions.Deployment)

	err = converter.Convert_Deployment_v1beta1_to_extensions(deployment, deploymentNew)
	if err != nil {
		return nil, err
	}

	selector, err := unversioned.LabelSelectorAsSelector(deploymentNew.Spec.Selector)
	if err != nil {
		return nil, err
	}

	options := api.ListOptions{LabelSelector: selector}

	podChannel := common.GetPodListChannelWithOptions(client, common.NewSameNamespaceQuery(namespace), options, 1)

	podListRaw := <-podChannel.List
	if err := <-podChannel.Error; err != nil {
		return nil, err
	}

	pods := common.FilterNamespacedPodsBySelector(podListRaw.Items, deployment.ObjectMeta.Namespace, deployment.Spec.Selector.MatchLabels)

	podList := pod.CreatePodList(pods)

	return &Deployment{
		ObjectMeta: common.NewObjectMeta(deploymentNew.ObjectMeta),
		TypeMeta:   common.NewTypeMeta(kind),
		PodList:    *podList,
		Selector:   deploymentNew.Spec.Selector.MatchLabels,
	}, nil
}

func List(client k8s.IK8S, namespace string) ([]Deployment, error) {

	deploymentList, err := client.Extensions().Deployments(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentNewList = []Deployment{}

	for _, val := range deploymentList.Items {

		var deploymentNew = extensions.Deployment{}

		err = converter.Convert_Deployment_v1beta1_to_extensions(&val, &deploymentNew)
		if err != nil {
			return nil, err
		}

		selector, err := unversioned.LabelSelectorAsSelector(deploymentNew.Spec.Selector)
		if err != nil {
			return nil, err
		}

		options := api.ListOptions{LabelSelector: selector}

		podChannel := common.GetPodListChannelWithOptions(client, common.NewSameNamespaceQuery(namespace), options, 1)

		podListRaw := <-podChannel.List
		if err := <-podChannel.Error; err != nil {
			return nil, err
		}

		pods := common.FilterNamespacedPodsBySelector(podListRaw.Items, deploymentNew.ObjectMeta.Namespace, deploymentNew.Spec.Selector.MatchLabels)

		podList := pod.CreatePodList(pods)

		deploymentNewList = append(deploymentNewList, Deployment{
			ObjectMeta: common.NewObjectMeta(deploymentNew.ObjectMeta),
			TypeMeta:   common.NewTypeMeta(kind),
			PodList:    *podList,
			Selector:   deploymentNew.Spec.Selector.MatchLabels,
		})
	}

	return nil, nil
}
