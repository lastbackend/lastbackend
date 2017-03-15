package k8s

import (
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/lb"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	*kubernetes.Clientset
	*lb.LBClientset
}

func Get(conf *rest.Config) (*Client, error) {

	kb, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err.Error())
	}

	lb, err := lb.NewForConfig(conf)
	if err != nil {
		panic(err.Error())
	}

	return &Client{kb, lb}, nil
}
