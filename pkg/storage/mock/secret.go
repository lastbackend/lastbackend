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

package mock

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

type SecretStorage struct {
	storage.Secret
	data map[string]*types.Secret
}

// Get secret by name
func (s *SecretStorage) Get(ctx context.Context, namespace, name string) (*types.Secret, error) {
	if r, ok := s.data[s.keyCreate(namespace, name)]; ok {
		return r, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

// Get secrets by namespace name
func (s *SecretStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Secret, error) {
	list := make(map[string]*types.Secret, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Insert new secret
func (s *SecretStorage) Insert(ctx context.Context, secret *types.Secret) error {
	if err := s.checkSecretArgument(secret); err != nil {
		return err
	}

	s.data[s.keyGet(secret)] = secret

	return nil
}

// Update secret info
func (s *SecretStorage) Update(ctx context.Context, secret *types.Secret) error {

	if err := s.checkSecretExists(secret); err != nil {
		return err
	}

	s.data[s.keyGet(secret)] = secret

	return nil
}

// Remove secret from storage
func (s *SecretStorage) Remove(ctx context.Context, secret *types.Secret) error {

	if err := s.checkSecretExists(secret); err != nil {
		return err
	}

	delete(s.data, s.keyGet(secret))

	return nil
}

// Clear secret storage
func (s *SecretStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]*types.Secret)
	return nil
}

// keyCreate util function
func (s *SecretStorage) keyCreate(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

// keyCreate util function
func (s *SecretStorage) keyGet(r *types.Secret) string {
	return r.SelfLink()
}

// newSecretStorage returns new storage
func newSecretStorage() *SecretStorage {
	s := new(SecretStorage)
	s.data = make(map[string]*types.Secret)
	return s
}

// checkSecretArgument - check if argument is valid for manipulations
func (s *SecretStorage) checkSecretArgument(secret *types.Secret) error {

	if secret == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if secret.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

// checkSecretArgument - check if secret exists in store
func (s *SecretStorage) checkSecretExists(secret *types.Secret) error {

	if err := s.checkSecretArgument(secret); err != nil {
		return err
	}

	if _, ok := s.data[s.keyGet(secret)]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}
