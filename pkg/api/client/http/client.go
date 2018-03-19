//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package http

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
)

type Client struct {
	*ClusterClient
	*DeploymentClient
	*EventsClient
	*NamespaceClient
	*NodeClient
	*RouteClient
	*ServiceClient
	*TriggerClient
	*VolumeClient
}

func (s *Client) Cluster() interfaces.Cluster {
	if s == nil {
		return nil
	}
	return s.ClusterClient
}

func (s *Client) Deployment() interfaces.Deployment {
	if s == nil {
		return nil
	}
	return s.DeploymentClient
}

func (s *Client) Events() interfaces.Events {
	if s == nil {
		return nil
	}
	return s.EventsClient
}

func (s *Client) Namespace() interfaces.Namespace {
	if s == nil {
		return nil
	}
	return s.NamespaceClient
}

func (s *Client) Node() interfaces.Node {
	if s == nil {
		return nil
	}
	return s.NodeClient
}

func (s *Client) Route() interfaces.Route {
	if s == nil {
		return nil
	}
	return s.RouteClient
}

func (s *Client) Service() interfaces.Service {
	if s == nil {
		return nil
	}
	return s.ServiceClient
}

func (s *Client) Trigger() interfaces.Trigger {
	if s == nil {
		return nil
	}
	return s.TriggerClient
}

func (s *Client) Volume() interfaces.Volume {
	if s == nil {
		return nil
	}
	return s.VolumeClient
}

func New(ctx context.Context) (*Client, error) {

	s := new(Client)

	s.ClusterClient = newClusterClient()
	s.DeploymentClient = newDeploymentClient()
	s.EventsClient = newEventsClient()
	s.NamespaceClient = newNamespaceClient()
	s.NodeClient = newNodeClient()
	s.RouteClient = newRouteClient()
	s.ServiceClient = newServiceClient()
	s.TriggerClient = newTriggerClient()
	s.VolumeClient = newVolumeClient()

	return s, nil
}
