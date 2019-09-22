//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logLevel         = 3
	logMonitorPrefix = "monitor:"
)

type Monitor struct {
	sync     sync.RWMutex
	watchers map[chan *types.Event]bool
}

func (m *Monitor) Subscribe(subscriber chan *types.Event, done chan bool) {

	m.sync.Lock()
	m.watchers[subscriber] = true
	m.sync.Unlock()

	log.V(logLevel).Debugf("%s:watch:> subscribe ", logMonitorPrefix)
	<-done
	log.V(logLevel).Debugf("%s:watch:> unsubscribe ", logMonitorPrefix)

	m.sync.Lock()
	delete(m.watchers, subscriber)
	m.sync.Unlock()
}

func (m *Monitor) Watch(ctx context.Context, stg storage.Storage, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> start ", logMonitorPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	var (
		f = `\b(\w+)\/?(\w+)?(\/\w+)?\/(.+)\b`
		c = stg.Collection().Root()
	)

	r, err := regexp.Compile(f)
	if err != nil {
		log.Errorf("%s:> filter compile err: %v", logMonitorPrefix, err.Error())
		return err
	}

	go func() {

		for {
			select {
			case <-ctx.Done():
				done <- true
				return
			case e := <-watcher:

				if e.Data == nil {
					continue
				}

				keys := r.FindStringSubmatch(e.Storage.Key)

				if len(keys) == 0 {
					continue
				}

				res := types.Event{}
				res.Action = e.Action
				res.Name = e.Name
				res.SelfLink = e.SelfLink
				res.Timestamp = time.Now()

				switch keys[1] {
				case types.KindNamespace:
					res.Kind = types.KindNamespace
					entity := new(types.Namespace)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindService:
					res.Kind = types.KindService
					entity := new(types.Service)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}

					res.Data = entity

				case types.KindDeployment:
					res.Kind = types.KindDeployment
					entity := new(types.Deployment)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}

					res.Data = entity

				case types.KindJob:
					res.Kind = types.KindJob
					entity := new(types.Job)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}

					res.Data = entity

				case types.KindPod:
					res.Kind = types.KindPod
					entity := new(types.Pod)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindSecret:
					res.Kind = types.KindSecret
					entity := new(types.Secret)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindConfig:
					res.Kind = types.KindConfig
					entity := new(types.Config)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindVolume:
					res.Kind = types.KindVolume
					entity := new(types.Volume)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindRoute:
					res.Kind = types.KindRoute
					entity := new(types.Route)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindNode:

					if keys[2] != "info" {
						continue
					}

					res.Kind = types.KindNode
					entity := new(types.Node)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindDiscovery:
					res.Kind = types.KindDiscovery

					if keys[2] != "info" {
						continue
					}

					entity := new(types.Discovery)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindIngress:
					res.Kind = types.KindIngress

					if keys[2] != "info" {
						continue
					}

					entity := new(types.Ingress)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity

				case types.KindCluster:
					res.Kind = types.KindCluster

					if keys[1] != "info" {
						continue
					}

					entity := new(types.Cluster)
					if err := json.Unmarshal(e.Data.([]byte), entity); err != nil {
						log.Errorf("%s:> parse data err: %v", logMonitorPrefix, err)
						continue
					}
					res.Data = entity
				default:
					continue
				}

				if err := m.dispatch(ctx, &res); err != nil {
					log.Errorf("%s> dispatch err: %s", logMonitorPrefix, err.Error())
				}

			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := stg.Watch(ctx, c, watcher, opts); err != nil {
		return err
	}

	return nil
}

func (m *Monitor) dispatch(ctx context.Context, event *types.Event) error {

	m.sync.Lock()
	for c := range m.watchers {
		go func() {
			fmt.Println("sent event:> dispatcher", event.Action, event.SelfLink, event.Timestamp)
			c <- event
		}()
	}
	m.sync.Unlock()

	return nil
}

func New() *Monitor {

	var m = Monitor{}
	m.watchers = make(map[chan *types.Event]bool)

	return &m
}
