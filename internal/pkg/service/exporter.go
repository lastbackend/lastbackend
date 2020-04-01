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
	logExporterPrefix = "distribution:exporter"
)

type Exporter struct {
	context context.Context
	storage storage.IStorage
}

func (n *Exporter) List() (*models.ExporterList, error) {
	list := models.NewExporterList()

	if err := n.storage.List(n.context, n.storage.Collection().Exporter().Info(), "", list, nil); err != nil {
		log.V(logLevel).Errorf("%s:list:> get exporter list err: %v", logExporterPrefix, err)
		return nil, err
	}

	return list, nil
}

func (n *Exporter) Put(exporter *models.Exporter) error {

	log.V(logLevel).Debugf("%s:create:> create exporter in cluster", logExporterPrefix)

	if err := n.storage.Put(n.context, n.storage.Collection().Exporter().Info(),
		exporter.SelfLink().String(), exporter, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert exporter err: %v", logExporterPrefix, err)
		return err
	}

	return nil
}

func (n *Exporter) Set(exporter *models.Exporter) error {

	log.V(logLevel).Debugf("%s:create:> create exporter in cluster", logExporterPrefix)

	if err := n.storage.Set(n.context, n.storage.Collection().Exporter().Info(),
		exporter.SelfLink().String(), exporter, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert exporter err: %v", logExporterPrefix, err)
		return err
	}

	return nil
}

func (n *Exporter) Get(selflink string) (*models.Exporter, error) {

	log.V(logLevel).Debugf("%s:get:> get by selflink %s", logExporterPrefix, selflink)

	exporter := new(models.Exporter)
	err := n.storage.Get(n.context, n.storage.Collection().Exporter().Info(), selflink, exporter, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get: exporter %s not found", logExporterPrefix, selflink)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get exporter `%s` err: %v", logExporterPrefix, selflink, err)
		return nil, err
	}

	return exporter, nil
}

func (n *Exporter) Remove(exporter *models.Exporter) error {

	log.V(logLevel).Debugf("%s:remove:> remove exporter %s", logExporterPrefix, exporter.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Exporter().Info(), exporter.SelfLink().String()); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove exporter err: %v", logExporterPrefix, err)
		return err
	}

	return nil
}

func (n *Exporter) Watch(ch chan models.ExporterEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch routes", logExporterPrefix)

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

				res := models.ExporterEvent{}
				res.Action = e.Action
				res.Name = e.Name

				exporter := new(models.Exporter)

				if err := json.Unmarshal(e.Data.([]byte), exporter); err != nil {
					log.Errorf("%s:> parse data err: %v", logExporterPrefix, err)
					continue
				}

				res.Data = exporter

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Exporter().Info(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewExporterModel(ctx context.Context, stg storage.IStorage) *Exporter {
	return &Exporter{ctx, stg}
}
