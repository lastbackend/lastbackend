package v1

import (
	"fmt"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/apimachinery/registered"
	"k8s.io/client-go/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

type LBInterface interface {
	RESTClient() rest.Interface
	PodsGetter
}

type LBClient struct {
	restClient rest.Interface
}

func (c *LBClient) Pods(namespace string) PodInterface {
	return newPods(c, namespace)
}

// NewForConfig creates a new LBClient for the given config.
func NewForConfig(c *rest.Config) (*LBClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &LBClient{client}, nil
}

// NewForConfigOrDie creates a new LBClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *LBClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new LBClient for the given RESTClient.
func New(c rest.Interface) *LBClient {
	return &LBClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv, err := unversioned.ParseGroupVersion("/v1")
	if err != nil {
		return err
	}
	// if /v1 is not enabled, return an error
	if !registered.IsEnabledVersion(gv) {
		return fmt.Errorf("/v1 is not enabled")
	}
	config.APIPath = "/api"
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	copyGroupVersion := gv
	config.GroupVersion = &copyGroupVersion

	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *LBClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
