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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"github.com/spf13/viper"
	"strings"
)

type IRoute interface {
	Get(namespace, name string) (*types.Route, error)
	ListByNamespace(namespace string) (map[string]*types.Route, error)
	Create(namespace *types.Namespace, opts *types.RouteCreateOptions) (*types.Route, error)
	Update(route *types.Route, namespace *types.Namespace, opts *types.RouteUpdateOptions) (*types.Route, error)
	SetStatus(route *types.Route, status *types.RouteStatus) error
	Remove(route *types.Route) error
}

type Route struct {
	context context.Context
	storage storage.Storage
}

func (n *Route) Get(namespace, name string) (*types.Route, error) {

	log.V(logLevel).Debug("api:distribution:route: get route by id %s/%s", namespace, name)

	item, err := n.storage.Route().Get(n.context, namespace, name)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("api:distribution:route:get: in namespace %s by name %s not found", namespace, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("api:distribution:route:get: in namespace %s by name %s error: %s", namespace, name, err)
		return nil, err
	}

	return item, nil
}

func (n *Route) ListByNamespace(namespace string) (map[string]*types.Route, error) {

	log.V(logLevel).Debug("api:distribution:route: list route")

	items, err := n.storage.Route().ListByNamespace(n.context, namespace)
	if err != nil {
		log.V(logLevel).Error("api:distribution:route: list route err: %s", err)
		return items, err
	}

	log.V(logLevel).Debugf("api:distribution:route: list route result: %d", len(items))

	return items, nil
}

func (n *Route) Create(namespace *types.Namespace, opts *types.RouteCreateOptions) (*types.Route, error) {

	log.V(logLevel).Debugf("api:distribution:route:crete create route %#v", opts)

	route := new(types.Route)
	route.Meta.SetDefault()
	route.Meta.Name = generator.GenerateRandomString(10)
	route.Meta.Namespace = namespace.Meta.Name
	route.Meta.Security = opts.Security
	route.Status.Stage = types.StageInitialized

	if len(opts.Domain) != 0 && opts.Custom {
		route.Spec.Domain = strings.ToLower(opts.Domain)
	}

	if len(opts.Domain) == 0 && len(opts.Subdomain) != 0 && !opts.Custom {
		route.Spec.Domain = strings.ToLower(fmt.Sprintf("%s-%s.%s", opts.Subdomain, namespace.Meta.Endpoint, viper.GetString("domain.external")))
	}

	if len(opts.Domain) == 0 && len(opts.Subdomain) == 0 && !opts.Custom {
		route.Spec.Domain = strings.ToLower(strings.Join([]string{generator.GenerateRandomString(5), namespace.Meta.Endpoint}, "-"))
	}

	route.Spec.Rules = make([]*types.RouteRule, 0)
	for _, rule := range opts.Rules {
		route.Spec.Rules = append(route.Spec.Rules, &types.RouteRule{
			Endpoint: *rule.Endpoint,
			Port:     *rule.Port,
			Path:     rule.Path,
		})
	}

	if err := n.storage.Route().Insert(n.context, route); err != nil {
		log.V(logLevel).Errorf("api:distribution:route:crete insert route err: %s", err)
		return nil, err
	}

	return route, nil
}

func (n *Route) Update(route *types.Route, namespace *types.Namespace, opts *types.RouteUpdateOptions) (*types.Route, error) {

	log.V(logLevel).Debugf("api:distribution:route:update update route %s", route.Meta.Name)

	route.Meta.SetDefault()
	route.Meta.Security = opts.Security
	route.Status.Stage = types.StageProvision

	route.Spec.Domain = opts.Domain
	route.Spec.Rules = make([]*types.RouteRule, 0)
	for _, rule := range opts.Rules {
		route.Spec.Rules = append(route.Spec.Rules, &types.RouteRule{
			Endpoint: *rule.Endpoint,
			Port:     *rule.Port,
			Path:     rule.Path,
		})
	}

	if len(opts.Domain) == 0 {
		route.Spec.Domain = strings.Join([]string{generator.GenerateRandomString(5), namespace.Meta.Endpoint}, "-")
	}

	if err := n.storage.Route().Update(n.context, route); err != nil {
		log.V(logLevel).Errorf("api:distribution:route:update update route err: %s", err)
		return nil, err
	}

	return route, nil
}

func (n *Route) SetStatus(route *types.Route, status *types.RouteStatus) error {

	if route == nil {
		log.V(logLevel).Warnf("api:distribution:route:setstatus: invalid argument %v", route)
		return nil
	}

	log.V(logLevel).Debugf("api:distribution:route:setstate set state route %s -> %#v", route.Meta.Name, status)

	route.Status = *status
	if err := n.storage.Route().SetStatus(n.context, route); err != nil {
		log.Errorf("Pod set status err: %s", err.Error())
		return err
	}


	return nil
}

func (n *Route) Remove(route *types.Route) error {

	log.V(logLevel).Debugf("api:distribution:route:remove remove route %#v", route)

	if err := n.storage.Route().Remove(n.context, route); err != nil {
		log.V(logLevel).Errorf("api:distribution:route:remove remove route  err: %s", err)
		return err
	}

	return nil
}

func NewRouteModel(ctx context.Context, stg storage.Storage) IRoute {
	return &Route{ctx, stg}
}
