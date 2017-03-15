package v1

import (
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

// PodsGetter has a method to return a PodInterface.
// A group's client should implement this interface.
type PodsGetter interface {
	Pods(namespace string) PodInterface
}

// PodInterface has methods to work with Pod resources.
type PodInterface interface {
	Attach(name string, opts *v1.PodAttachOptions) *rest.Request
}

// pods implements PodInterface
type pods struct {
	client rest.Interface
	ns     string
}

// newPods returns a Pods
func newPods(c *LBClient, namespace string) *pods {
	return &pods{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

func (c *pods) Attach(name string, opts *v1.PodAttachOptions) *rest.Request {
	return c.client.Post().
		Namespace(c.ns).
		Name(name).
		Resource("pods").
		SubResource("attach").
		VersionedParams(opts, api.ParameterCodec)
}
