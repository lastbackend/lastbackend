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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"fmt"
)

type ISecret interface {
	Get(namespace, name string) (*types.Secret, error)
	ListByNamespace(namespace string) (map[string]*types.Secret, error)
	Create(namespace *types.Namespace, opts *types.SecretCreateOptions) (*types.Secret, error)
	Update(secret *types.Secret, namespace *types.Namespace, opts *types.SecretUpdateOptions) (*types.Secret, error)
	Remove(secret *types.Secret) error
}

type Secret struct {
	context context.Context
	storage storage.Storage
}

func (n *Secret) Get(namespace, name string) (*types.Secret, error) {

	log.V(logLevel).Debug("api:distribution:secret: get secret by id %s/%s", namespace, name)

	item, err := n.storage.Secret().Get(n.context, namespace, name)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("api:distribution:secret:get: in namespace %s by name %s not found", namespace, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("api:distribution:secret:get: in namespace %s by name %s error: %s", namespace, name, err)
		return nil, err
	}

	return item, nil
}

func (n *Secret) ListByNamespace(namespace string) (map[string]*types.Secret, error) {

	log.V(logLevel).Debug("api:distribution:secret: list secret")

	items, err := n.storage.Secret().ListByNamespace(n.context, namespace)
	if err != nil {
		log.V(logLevel).Error("api:distribution:secret: list secret err: %s", err)
		return items, err
	}

	log.V(logLevel).Debugf("api:distribution:secret: list secret result: %d", len(items))

	return items, nil
}

func (n *Secret) Create(namespace *types.Namespace, opts *types.SecretCreateOptions) (*types.Secret, error) {

	log.V(logLevel).Debugf("api:distribution:secret:crete create secret %#v", opts)

	secret := new(types.Secret)
	secret.Meta.SetDefault()
	secret.Meta.Name = generator.GenerateRandomString(10)
	secret.Meta.Namespace = namespace.Meta.Name
	if opts.Data != nil {
		secret.Data = *opts.Data
	}

	if err := n.storage.Secret().Insert(n.context, secret); err != nil {
		log.V(logLevel).Errorf("api:distribution:secret:crete insert secret err: %s", err)
		return nil, err
	}

	return secret, nil
}

func (n *Secret) Update(secret *types.Secret, namespace *types.Namespace, opts *types.SecretUpdateOptions) (*types.Secret, error) {

	log.V(logLevel).Debugf("api:distribution:secret:update update secret %s", secret.Meta.Name)

	if opts.Data != nil {
		secret.Data = *opts.Data
	}
fmt.Println(">>>>>>>>", secret.Data)
	if err := n.storage.Secret().Update(n.context, secret); err != nil {
		log.V(logLevel).Errorf("api:distribution:secret:update update secret err: %s", err)
		return nil, err
	}

	return secret, nil
}

func (n *Secret) Remove(secret *types.Secret) error {

	log.V(logLevel).Debugf("api:distribution:secret:remove remove secret %#v", secret)

	if err := n.storage.Secret().Remove(n.context, secret); err != nil {
		log.V(logLevel).Errorf("api:distribution:secret:remove remove secret  err: %s", err)
		return err
	}

	return nil
}

func NewSecretModel(ctx context.Context, stg storage.Storage) ISecret {
	return &Secret{ctx, stg}
}
