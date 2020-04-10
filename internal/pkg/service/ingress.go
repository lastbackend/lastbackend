//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package service

import (
	"context"
	"encoding/json"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logIngressPrefix = "distribution:ingress"
)

type Ingress struct {
	context context.Context
	storage storage.IStorage
}

func (n *Ingress) List() (*models.IngressList, error) {
	list := models.NewIngressList()

	if err := n.storage.List(n.context, n.storage.Collection().Ingress().Info(), "", list, nil); err != nil {
		log.V(logLevel).Errorf("%s:list:> get ingress list err: %v", logIngressPrefix, err)
		return nil, err
	}

	return list, nil
}

func (n *Ingress) Put(ingress *models.Ingress) error {

	log.V(logLevel).Debugf("%s:create:> create ingress in cluster", logIngressPrefix)

	if err := n.storage.Put(n.context, n.storage.Collection().Ingress().Info(),
		ingress.SelfLink().String(), ingress, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert ingress err: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func (n *Ingress) Set(ingress *models.Ingress) error {

	log.V(logLevel).Debugf("%s:create:> create ingress in cluster", logIngressPrefix)

	if err := n.storage.Set(n.context, n.storage.Collection().Ingress().Info(),
		ingress.SelfLink().String(), ingress, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert ingress err: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func (n *Ingress) Get(selflink string) (*models.Ingress, error) {

	log.V(logLevel).Debugf("%s:get:> get by selflink %s", logIngressPrefix, selflink)

	ingress := new(models.Ingress)
	err := n.storage.Get(n.context, n.storage.Collection().Ingress().Info(), selflink, ingress, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get: ingress %s not found", logIngressPrefix, selflink)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get ingress `%s` err: %v", logIngressPrefix, selflink, err)
		return nil, err
	}

	return ingress, nil
}

func (n *Ingress) Remove(ingress *models.Ingress) error {

	log.V(logLevel).Debugf("%s:remove:> remove ingress %s", logIngressPrefix, ingress.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Ingress().Info(), ingress.SelfLink().String()); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove ingress err: %v", logIngressPrefix, err)
		return err
	}

	return nil
}

func (n *Ingress) Watch(ch chan models.IngressEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch routes", logIngressPrefix)

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

				res := models.IngressEvent{}
				res.Action = e.Action
				res.Name = e.Name

				ingress := new(models.Ingress)

				if err := json.Unmarshal(e.Data.([]byte), ingress); err != nil {
					log.Errorf("%s:> parse data err: %v", logIngressPrefix, err)
					continue
				}

				res.Data = ingress

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Ingress().Info(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewIngressModel(ctx context.Context, stg storage.IStorage) *Ingress {
	return &Ingress{ctx, stg}
}
