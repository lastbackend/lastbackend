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
	"strings"
	"time"

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"github.com/spf13/viper"
)

const (
	logPodPrefix = "distribution:pod"
)

type Pod struct {
	context context.Context
	storage storage.Storage
}

// Get pod info from storage
func (p *Pod) Get(namespace, service, deployment, name string) (*types.Pod, error) {
	log.V(logLevel).Debugf("%s:get:> get by name %s", logPodPrefix, name)

	pod := new(types.Pod)

	err := p.storage.Get(p.context, storage.PodKind,
		p.storage.Key().Pod(namespace, service, deployment, name), pod, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> `%s` not found", logPodPrefix, name)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get Pod `%s` err: %v", logPodPrefix, name, err)
		return nil, err
	}

	return pod, nil
}

// Create new pod
func (p *Pod) Create(deployment *types.Deployment) (*types.Pod, error) {

	pod := types.NewPod()
	pod.Meta.SetDefault()
	pod.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	pod.Meta.Deployment = deployment.Meta.Name
	pod.Meta.Service = deployment.Meta.Service
	pod.Meta.Namespace = deployment.Meta.Namespace

	pod.Status.SetCreated()
	pod.Status.Steps = make(map[string]types.PodStep)
	pod.Status.Steps[types.StepInitialized] = types.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}

	var ips = make([]string, 0)
	viper.UnmarshalKey("dns.ips", &ips)
	ips = append(ips, "8.8.8.8")

	for _, s := range deployment.Spec.Template.Containers {
		s.Labels = make(map[string]string)
		s.Labels["LB"] = pod.SelfLink()
		s.DNS = types.SpecTemplateContainerDNS{
			Server: ips,
			Search: ips,
		}
		pod.Spec.Template.Containers = append(pod.Spec.Template.Containers, s)
	}

	for _, s := range deployment.Spec.Template.Volumes {
		pod.Spec.Template.Volumes = append(pod.Spec.Template.Volumes, s)
	}

	if err := p.storage.Put(p.context, storage.PodKind,
		p.storage.Key().Pod(pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name), pod, nil); err != nil {
		log.Errorf("%s:create:> insert pod err %v", logPodPrefix, err)
		return nil, err
	}

	return pod, nil
}

// ListByNamespace returns pod list in selected namespace
func (p *Pod) ListByNamespace(namespace string) (*types.PodList, error) {
	log.V(logLevel).Debugf("%s:listbynamespace:> get pod list by namespace %s", logPodPrefix, namespace)

	list := types.NewPodList()
	filter := p.storage.Filter().Pod().ByNamespace(namespace)

	err := p.storage.List(p.context, storage.PodKind, filter, list, nil)
	if err != nil {
		log.V(logLevel).Debugf("%s:listbynamespace:> get pod list by deployment id `%s` err: %v", logPodPrefix, namespace, err)
		return nil, err
	}

	return list, nil
}

// ListByService returns pod list in selected service
func (p *Pod) ListByService(namespace, service string) (*types.PodList, error) {
	log.V(logLevel).Debugf("%s:listbyservice:> get pod list by service id %s/%s", logPodPrefix, namespace, service)

	list := types.NewPodList()
	filter := p.storage.Filter().Pod().ByService(namespace, service)

	err := p.storage.List(p.context, storage.PodKind, filter, list, nil)
	if err != nil {
		log.V(logLevel).Debugf("%s:listbyservice:> get pod list by service id `%s` err: %v", logPodPrefix, namespace, service, err)
		return nil, err
	}

	return list, nil
}

// ListByDeployment returns pod list in selected deployment
func (p *Pod) ListByDeployment(namespace, service, deployment string) (*types.PodList, error) {
	log.V(logLevel).Debugf("%s:listbydeployment:> get pod list by id %s/%s/%s", logPodPrefix, namespace, service, deployment)

	list := types.NewPodList()
	filter := p.storage.Filter().Pod().ByDeployment(namespace, service, deployment)

	err := p.storage.List(p.context, storage.PodKind, filter, list, nil)
	if err != nil {
		log.V(logLevel).Debugf("%s:listbydeployment:> get pod list by deployment id `%s/%s/%s` err: %v",
			logPodPrefix, namespace, service, deployment, err)
		return nil, err
	}

	return list, nil
}

// SetNode - set node info to pod
func (p *Pod) SetNode(pod *types.Pod, node *types.Node) error {
	log.Debugf("%s:setnode:> set node for pod: %s", logPodPrefix, pod.Meta.Name)

	pod.Meta.Node = node.Meta.Name

	if err := p.storage.Set(p.context, storage.PodKind,
		p.storage.Key().Pod(pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name), pod, nil); err != nil {
		log.Errorf("%s:setnode:> pod set node err: %v", logPodPrefix, err)
		return err
	}

	return nil
}

// SetStatus - set state for pod
func (p *Pod) Update(pod *types.Pod) error {

	log.Debugf("%s:update:> update pod: %s", logPodPrefix, pod.Meta.Name)

	if err := p.storage.Set(p.context, storage.PodKind,
		p.storage.Key().Pod(pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name),
		pod, nil); err != nil {
		log.Errorf("%s:update:> pod update err: %v", logPodPrefix, err)
		return err
	}

	return nil
}

// Destroy pod
func (p *Pod) Destroy(pod *types.Pod) error {

	pod.Spec.State.Destroy = true

	if err := p.storage.Set(p.context, storage.PodKind,
		p.storage.Key().Pod(pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name), pod, nil); err != nil {
		log.Errorf("%s:destroy:> mark pod for destroy error: %v", logPodPrefix, err)
		return err
	}
	return nil
}

// Remove pod from storage
func (p *Pod) Remove(pod *types.Pod) error {
	if err := p.storage.Del(p.context, storage.PodKind,
		p.storage.Key().Pod(pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name)); err != nil {
		log.Errorf("%s:remove:> mark pod for destroy error: %v", logPodPrefix, err)
		return err
	}
	return nil
}

func (p *Pod) Watch(ch chan types.PodEvent) error {
	log.Debugf("%s:watch:> watch service", logPodPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-p.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.PodEvent{}
				res.Action = e.Action
				res.Name = e.Name

				obj := new(types.Pod)

				if err := json.Unmarshal(e.Data.([]byte), &obj); err != nil {
					log.Errorf("%s:watch:> parse json", logPodPrefix)
					continue
				}

				res.Data = obj

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := p.storage.Watch(p.context, storage.PodKind, watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewPodModel(ctx context.Context, stg storage.Storage) *Pod {
	return &Pod{ctx, stg}
}
