package k8s

import (
	"k8s.io/client-go/1.5/rest"
	"k8s.io/client-go/1.5/kubernetes"

)

// Interface exposes methods on k8s resources.
type LBClientsetInterface interface {
	LB() LBClientInterface
}

type LBClientInterface interface {
	GetRESTClient() *rest.RESTClient
	ComponentAccountsGetter
}

type LBClientset struct {
	*LBClient
}

func (c *LBClientset) LB() LBClientInterface {
	if c == nil {
		return nil
	}
	return c.LBClient
}

// LBExtend is used to interact with features provided by the Last.Backend group.
type LBClient struct {
	*rest.RESTClient
}

func (c *LBClient) Accounts() ComponentAccountInterface {
	return newComponentAccounts(c)
}

func (c *LBClient) GetRESTClient() *rest.RESTClient {
	if c == nil {
		return nil
	}
	return c.RESTClient
}

type Client struct {
	*kubernetes.Clientset
	*LBClientset
}


func Get (conf *rest.Config) (*Client, error) {
	c, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err.Error())
	}

	return &Client{c}, nil
}

