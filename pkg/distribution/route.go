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

package distribution

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
)

type IRoute interface {
	Get(namespace, name string) (*types.Route, error)
	ListByNamespace(namespace string) (map[string]*types.Route, error)
	Create(namespace *types.Namespace, services map[string]*types.Service, opts *request.RouteCreateOptions) (*types.Route, error)
	Update(route *types.Route, namespace *types.Namespace, services map[string]*types.Service, opts *request.RouteUpdateOptions) (*types.Route, error)
	SetState(route *types.Route, state *types.RouteState) error
	Remove(route *types.Route) error
}

type Route struct {
	context context.Context
	storage storage.Storage
}

func (n *Route) Get(namespace, name string) (*types.Route, error) {

	log.V(logLevel).Debug("Route: get route by id %s/%s", namespace, name)

	item, err := n.storage.Route().Get(n.context, namespace, name)
	if err != nil {
		log.V(logLevel).Error("Route: get route err: %s", err)
		return nil, err
	}

	return item, nil
}

func (n *Route) ListByNamespace(namespace string) (map[string]*types.Route, error) {

	log.V(logLevel).Debug("Route: list route")

	items, err := n.storage.Route().ListByNamespace(n.context, namespace)
	if err != nil {
		log.V(logLevel).Error("Route: list route err: %s", err)
		return items, err
	}

	log.V(logLevel).Debugf("Route: list route result: %d", len(items))

	return items, nil
}

func (n *Route) Create(namespace *types.Namespace, services map[string]*types.Service, opts *request.RouteCreateOptions) (*types.Route, error) {

	log.V(logLevel).Debugf("Route: create route %#v", opts)

	var route = types.Route{}
	//route.Meta.SetDefault()
	//route.Meta.Namespace = namespace.Meta.Name
	//route.Meta.Security = opts.Security
	//route.State.Provision = true
	//route.Meta.Hash = generator.GenerateRandomString(5)
	//
	//if len(opts.Domain) != 0 && opts.Custom {
	//	route.Meta.Domain = strings.ToLower(opts.Domain)
	//}
	//
	//if len(opts.Domain) == 0 && len(opts.Subdomain) != 0 && !opts.Custom {
	//	route.Meta.Domain = strings.ToLower(fmt.Sprintf("%s-%s.%s", opts.Subdomain, namespace.Meta.Endpoint, viper.GetString("domain.external")))
	//}
	//
	//if len(opts.Domain) == 0 && len(opts.Subdomain) == 0 && !opts.Custom {
	//	route.Meta.Domain = strings.ToLower(strings.Join([]string{generator.GenerateRandomString(5), namespace.Meta.Endpoint}, "-"))
	//}
	//
	//ss := make(map[string]*types.Service)
	//for _, service := range services {
	//	ss[service.Meta.Name] = service
	//}
	//
	////route.Rules = make(map[string]*types.RouteRule, 0)
	////for _, rule := range opts.Rules {
	////	route.Rules = append(route.Rules, &types.RouteRule{
	////		Service:  *rule.Service,
	////		Port:     *rule.Port,
	////		Path:     rule.Path,
	////		Endpoint: ss[*rule.Service].Meta.Endpoint,
	////	})
	////}

	if err := n.storage.Route().Insert(n.context, &route); err != nil {
		log.V(logLevel).Errorf("Route: insert Route err: %s", err)
		return nil, err
	}

	return &route, nil
}

func (n *Route) Update(route *types.Route, namespace *types.Namespace, services map[string]*types.Service, opts *request.RouteUpdateOptions) (*types.Route, error) {

	log.V(logLevel).Debugf("Route: update route %s", route.Meta.Name)

	route.Meta.SetDefault()
	//route.Meta.Domain = opts.Domain
	//route.Meta.Namespace = namespace.Meta.Name
	//route.Meta.Security = opts.Security
	//route.State.Provision = true
	//route.Meta.Hash = generator.GenerateRandomString(5)

	//route.Rules = make(map[string]*types.RouteRule, 0)
	//for _, rule := range opts.Rules {
	//
	//	var svc *types.Service
	//	for _, s := range services {
	//		if *rule.Service == s.Meta.Name {
	//			svc = s
	//			break
	//		}
	//	}
	//
	//	route.Rules = append(route.Rules, &types.RouteRule{
	//		Service:  *rule.Service,
	//		Port:     *rule.Port,
	//		Path:     rule.Path,
	//		Endpoint: svc.Meta.Endpoint,
	//	})
	//}

	//if len(opts.Domain) == 0 {
	//	route.Meta.Domain = strings.Join([]string{generator.GenerateRandomString(5), namespace.Meta.Endpoint}, "-")
	//}

	if err := n.storage.Route().Update(n.context, route); err != nil {
		log.V(logLevel).Errorf("Route: update route err: %s", err)
		return nil, err
	}

	return route, nil
}

func (n *Route) SetState(route *types.Route, state *types.RouteState) error {

	log.V(logLevel).Debugf("Route: set state route %s -> %#v", route.Meta.Name, state)

	route.State.Destroy = state.Destroy
	route.State.Provision = state.Provision

	if err := n.storage.Route().Update(n.context, route); err != nil {
		log.V(logLevel).Errorf("Route: mark route as destroy err: %s", err)
		return err
	}

	return nil
}

func (n *Route) Remove(route *types.Route) error {

	log.V(logLevel).Debugf("Route: remove route %#v", route)

	if err := n.storage.Route().Remove(n.context, route); err != nil {
		log.V(logLevel).Errorf("Route: remove route  err: %s", err)
		return err
	}

	return nil
}

func NewRouteModel(ctx context.Context, stg storage.Storage) IRoute {
	return &Route{ctx, stg}
}
