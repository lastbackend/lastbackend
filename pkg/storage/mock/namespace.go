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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

var data map[string]*types.Namespace

// Namespace Service type for interface in interfaces folder
type NamespaceStorage struct {
	storage.Namespace
}

// Get namespace by name
func (s *NamespaceStorage) GetByName(ctx context.Context, name string) (*types.Namespace, error) {

	if ns, ok := data[name]; !ok {
		return nil, errors.New("n")
	}

	return getByName(name), nil
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) ([]*types.Namespace, error) {
	return []*types.Namespace{getByName("demo")}, nil
}

// Insert new namespace into storage
func (s *NamespaceStorage) Insert(ctx context.Context, namespace *types.Namespace) error {
	data[namespace.Meta.Name] = namespace
	return nil
}

// Update namespace model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) error {
	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, name string) error {
	return nil
}

func newNamespaceStorage() *NamespaceStorage {
	s := new(NamespaceStorage)
	return s
}

/* ============================================================================================================== */
/* =============================================== HELPER METHODS =============================================== */
/* ============================================================================================================== */

func createNamespace(name, description string) *types.Namespace {
	ns := new(types.Namespace)
	ns.Meta.SetDefault()
	ns.Meta.Name = name
	ns.Meta.Description = description
	ns.Meta.Endpoint = fmt.Sprintf("%s.demo.io", name)
	ns.Meta.SelfLink = fmt.Sprintf("/%s", name)
	ns.Env = make(types.NamespaceEnvs, 0)
	return ns
}

func getByName(name string) *types.Namespace {
	switch name {
	case "demo":
		ns := createNamespace(name, "demo description")
		ns.Env = types.NamespaceEnvs{
			types.NamespaceEnv{Name: "DEBUG", Value: "true"},
			types.NamespaceEnv{Name: "NODE_ENV", Value: "development"},
		}
		ns.Resources.RAM = 128
		ns.Resources.Routes = 1
		ns.Quotas.Routes = 1
		ns.Quotas.RAM = 256
		ns.Labels = map[string]string{"ns": "demo"}
		return ns
	default:
		return nil
	}
}
