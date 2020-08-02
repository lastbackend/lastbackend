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
//
//import (
//	"context"
//	"encoding/json"
//
//	"github.com/lastbackend/lastbackend/internal/pkg/errors"
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/internal/pkg/storage"
//	log "github.com/lastbackend/lastbackend/tools/logger"
//)
//
//const (
//	logClusterPrefix = "distribution:cluster"
//)
//
//// Cluster - distribution model
//type Cluster struct {
//	context context.Context
//	storage storage.IStorage
//}
//
//// Info - get cluster info
//func (c *Cluster) Get() (*models.Cluster, error) {
//
//	log.Debugf("%s:get:> get info", logClusterPrefix)
//
//	cluster := new(models.Cluster)
//	err := c.storage.Get(c.context, c.storage.Collection().Cluster(), models.EmptyString, cluster, nil)
//	if err != nil {
//		if errors.Storage().IsErrEntityNotFound(err) {
//			log.Warnf("%s:get:> cluster not found", logClusterPrefix)
//			return cluster, nil
//		}
//
//		log.Errorf("%s:get:> get cluster err: %v", logClusterPrefix, err)
//		return nil, err
//	}
//
//	return cluster, nil
//}
//
//func (c *Cluster) Set(cluster *models.Cluster) error {
//
//	log.Debugf("%s:set:> update Cluster %#v", logClusterPrefix, cluster)
//	opts := storage.GetOpts()
//	opts.Force = true
//
//	if err := c.storage.Set(c.context, c.storage.Collection().Cluster(), models.EmptyString, cluster, opts); err != nil {
//		log.Errorf("%s:set:> update Cluster err: %v", logClusterPrefix, err)
//		return err
//	}
//
//	return nil
//}
//
//// Watch cluster changes
//func (c *Cluster) Watch(ch chan models.ClusterEvent) {
//
//	log.Debugf("%s:watch:> watch cluster", logClusterPrefix)
//
//	done := make(chan bool)
//	watcher := storage.NewWatcher()
//
//	go func() {
//		for {
//			select {
//			case <-c.context.Done():
//				done <- true
//				return
//			case e := <-watcher:
//				if e.Data == nil {
//					continue
//				}
//
//				res := models.ClusterEvent{}
//				res.Name = e.Name
//				res.Action = e.Action
//
//				cluster := new(models.Cluster)
//
//				if err := json.Unmarshal(e.Data.([]byte), *cluster); err != nil {
//					log.Errorf("%s:> parse data err: %v", logClusterPrefix, err)
//					continue
//				}
//
//				res.Data = cluster
//
//				ch <- res
//			}
//		}
//	}()
//
//	opts := storage.GetOpts()
//	go c.storage.Watch(c.context, c.storage.Collection().Cluster(), watcher, opts)
//
//	<-done
//}
//
//// NewClusterModel - return new cluster model
//func NewClusterModel(ctx context.Context, stg storage.IStorage) *Cluster {
//	return &Cluster{ctx, stg}
//}
