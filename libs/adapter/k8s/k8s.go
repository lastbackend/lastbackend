package k8s

import (
	"k8s.io/client-go/1.5/rest"
	"k8s.io/client-go/1.5/kubernetes"
)

// Interface exposes methods on k8s resources.
type Interface interface {

}


type Client struct {
	*kubernetes.Clientset
}

func Get (conf *rest.Config) (*Client, error) {
	c, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err.Error())
	}

	return &Client{c}, nil
}