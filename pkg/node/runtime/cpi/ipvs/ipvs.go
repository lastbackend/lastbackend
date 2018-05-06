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

// +build linux

package ipvs

/*
#include <linux/types.h>
#include <linux/ip_vs.h>
*/
import "C"
import (
	"context"
)

// Proxy balancer
type IPVS struct {
	cms string
}


// GetServices method returns all services on current node
func (i *IPVS) GetServices(ctx context.Context) (map[string]*Service, error) {
	svcs := make(map[string]*Service)
	return svcs, nil
}

// GetService method returns particular service on current node
func (i *IPVS) GetService(ctx context.Context, name string) (*Service, error) {
	svc := new(Service)
	return svc, nil
}

// AddService method add particular service to current node
func (i *IPVS) AddService(ctx context.Context, svc *Service) error {
	return nil
}

// SetService method updates particular service on current node
func (i *IPVS) SetService(ctx context.Context, svc *Service) error {
	return nil
}

// DelService method delete service from current node by service name
func (i *IPVS) DelService(ctx context.Context, svc *Service) error {
	return nil
}

// GetBackends method returns all backend on current node depends particular service
func (i *IPVS) GetBackends(ctx context.Context, svc *Service) (map[string]*Backend, error) {
	bknds := make(map[string]*Backend)
	return bknds, nil
}

// GetBackend method returns particular backend on provided service
func (i *IPVS) GetBackend(ctx context.Context, svc *Service, name string) (*Backend, error) {
	bknd := new(Backend)
	return bknd, nil
}

// AddBackend method adds new backend to provided service
func (i *IPVS) AddBackend(ctx context.Context, svc *Service, bknd *Backend) error {
	return nil
}

// SetBackend method updates backend on provided service
func (i *IPVS) SetBackend(ctx context.Context, svc *Service, bknd *Backend) error {
	return nil
}

// DelBackend method removes backend from service and node
func (i *IPVS) DelBackend(ctx context.Context, svc *Service, bknd *Backend) error {
	return nil
}

// check is IPVS available for host
func (i *IPVS) check () error {
	return nil
}

