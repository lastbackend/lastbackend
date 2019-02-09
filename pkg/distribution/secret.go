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
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logSecretPrefix = "distribution:secret"
)

type Secret struct {
	context context.Context
	storage storage.Storage
}

func (n *Secret) Runtime() (*types.System, error) {

	log.V(logLevel).Debugf("%s:get:> get secret runtime info", logSecretPrefix)
	runtime, err := n.storage.Info(n.context, n.storage.Collection().Secret(), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> get runtime info error: %s", logSecretPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}

func (n *Secret) Get(namespace, name string) (*types.Secret, error) {

	log.V(logLevel).Debugf("%s:get:> get secret by id %s/%s", logSecretPrefix, name)

	item := new(types.Secret)
	sl := types.NewSecretSelfLink(namespace, name).String()

	err := n.storage.Get(n.context, n.storage.Collection().Secret(), sl, &item, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> in namespace %s by name %s not found", logSecretPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> in namespace %s by name %s error: %s", logSecretPrefix, name, err)
		return nil, err
	}

	return item, nil
}

func (n *Secret) List(filter string) (*types.SecretList, error) {

	var f string

	log.V(logLevel).Debugf("%s:list:> get secrets list by namespace", logSecretPrefix)

	list := types.NewSecretList()
	if filter != types.EmptyString {
		f = n.storage.Filter().Secret().ByNamespace(filter)
	}

	err := n.storage.List(n.context, n.storage.Collection().Secret(), f, list, nil)
	if err != nil {
		log.V(logLevel).Error("%s:list:> get secrets list by namespace err: %s", logSecretPrefix, err)
		return list, err
	}

	log.V(logLevel).Debugf("%s:list:> get secrets list by namespace result: %d", logSecretPrefix, len(list.Items))

	return list, nil
}

func (n *Secret) Create(namespace *types.Namespace, secret *types.Secret) (*types.Secret, error) {

	log.V(logLevel).Debugf("%s:create:> create secret %#v", logSecretPrefix, secret.Meta.Name)

	secret.Meta.SetDefault()
	secret.Meta.Namespace = namespace.Meta.Name
	secret.SelfLink()

	if err := n.storage.Put(n.context, n.storage.Collection().Secret(),
		secret.SelfLink().String(), secret, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert secret err: %v", logSecretPrefix, err)
		return nil, err
	}

	return secret, nil
}

func (n *Secret) Update(secret *types.Secret) (*types.Secret, error) {

	log.V(logLevel).Debugf("%s:update:> update secret %s", logSecretPrefix, secret.Meta.Name)

	if err := n.storage.Set(n.context, n.storage.Collection().Secret(),
		secret.SelfLink().String(), secret, nil); err != nil {
		log.V(logLevel).Errorf("%s:update:> update secret err: %s", logSecretPrefix, err)
		return nil, err
	}

	return secret, nil
}

func (n *Secret) Remove(secret *types.Secret) error {

	log.V(logLevel).Debugf("%s:remove:> remove secret %#v", logSecretPrefix, secret)

	if err := n.storage.Del(n.context, n.storage.Collection().Secret(),
		secret.SelfLink().String()); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove secret  err: %s", logSecretPrefix, err)
		return err
	}

	return nil
}

func (n *Secret) Watch(ch chan types.SecretEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch secret", logSecretPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-n.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.SecretEvent{}
				res.Action = e.Action
				res.Name = e.Name

				secret := new(types.Secret)

				if err := json.Unmarshal(e.Data.([]byte), secret); err != nil {
					log.Errorf("%s:> parse data err: %v", logSecretPrefix, err)
					continue
				}

				res.Data = secret

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Secret(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewSecretModel(ctx context.Context, stg storage.Storage) *Secret {
	return &Secret{ctx, stg}
}
