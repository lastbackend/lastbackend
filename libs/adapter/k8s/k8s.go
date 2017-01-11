package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	rc := kb.CoreV1().RESTClient()
	lb := &LBClientset{&LBClient{rc}}

	return &Client{kb, lb}, nil
}
