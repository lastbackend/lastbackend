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

package etcd

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"time"
)

const secretStorage = "secrets"

type SecretStorage struct {
	storage.Secret
}

// Get secret by name
func (s *SecretStorage) Get(ctx context.Context, namespace, name string) (*types.Secret, error) {

	log.V(logLevel).Debugf("storage:etcd:secret:> get by name: %s", name)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:secret:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("storage:etcd:secret:> get by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + secretStorage + `\/.+\/(meta|status|spec)\b`
	var (
		secret = new(types.Secret)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> get by name err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyMeta := keyCreate(secretStorage, s.keyCreate(namespace, name))
	if err := client.Map(ctx, keyMeta, filter, secret); err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> get by name err: %s", name, err.Error())
		return nil, err
	}

	if secret.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	return secret, nil
}

// Get secrets by namespace name
func (s *SecretStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Secret, error) {

	log.V(logLevel).Debugf("storage:etcd:secret:> get list by namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:secret:> get list by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + secretStorage + `\/(.+)\/(meta|status|spec)\b`

	var (
		secrets = make(map[string]*types.Secret)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> get list by namespace err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyDirCreate(secretStorage, fmt.Sprintf("%s:", namespace))
	if err := client.MapList(ctx, key, filter, secrets); err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> err: %s", namespace, err.Error())
		return nil, err
	}

	return secrets, nil

}

// Insert new secret
func (s *SecretStorage) Insert(ctx context.Context, secret *types.Secret) error {

	log.V(logLevel).Debugf("storage:etcd:secret:> insert secret: %#v", secret)

	if err := s.checkSecretArgument(secret); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> insert secret err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(secretStorage, s.keyGet(secret), "meta")
	if err := tx.Create(keyMeta, secret.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> insert secret err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(secretStorage, s.keyGet(secret), "data")
	if err := tx.Create(keySpec, secret.Data, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> insert secret err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> insert secret err: %s", err.Error())
		return err
	}

	return nil
}

// Update secret info
func (s *SecretStorage) Update(ctx context.Context, secret *types.Secret) error {

	log.V(logLevel).Debugf("storage:etcd:secret:> update secret: %#v", secret)

	if err := s.checkSecretExists(ctx, secret); err != nil {
		return err
	}

	secret.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> update secret err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(secretStorage, s.keyGet(secret), "meta")
	if err := client.Upsert(ctx, keyMeta, secret.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> update secret err: %s", err.Error())
		return err
	}

	return nil
}

// Remove secret from storage
func (s *SecretStorage) Remove(ctx context.Context, secret *types.Secret) error {

	log.V(logLevel).Debugf("storage:etcd:secret:> remove secret: %#v", secret)

	if err := s.checkSecretExists(ctx, secret); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> remove err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(secretStorage, s.keyGet(secret))
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> remove secret err: %s", err.Error())
		return err
	}

	return nil
}

// Clear secret storage
func (s *SecretStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:secret:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, secretStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:secret:> clear err: %s", err.Error())
		return err
	}

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

func newSecretStorage() *SecretStorage {
	s := new(SecretStorage)
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
func (s *SecretStorage) checkSecretExists(ctx context.Context, secret *types.Secret) error {

	if err := s.checkSecretArgument(secret); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:secret:> check secret exists")

	if _, err := s.Get(ctx, secret.Meta.Namespace, secret.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:secret:> check secret exists err: %s", err.Error())
		return err
	}

	return nil
}
