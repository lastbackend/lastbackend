package k8s

import (
	"k8s.io/client-go/rest"
)

// LBExtend is used to interact with features provided by the Last.Backend group.

type LBClientInterface interface {
	GetRESTClient() rest.Interface
}

type LBClient struct {
	rest.Interface
}

func (c *LBClient) GetRESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c
}
