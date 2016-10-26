package k8s

import (
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/rest"
)

type Client struct {
	*kubernetes.Clientset
	*LBClientset
}

func Get(conf *rest.Config) (*Client, error) {

	kb, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err.Error())
	}

	rc := kb.Core().GetRESTClient()
	lb := &LBClientset{&LBClient{rc}}

	return &Client{kb, lb}, nil
}
