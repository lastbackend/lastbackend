//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

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
