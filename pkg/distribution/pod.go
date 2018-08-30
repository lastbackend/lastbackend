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

	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"github.com/spf13/viper"
	"regexp"
)

const (
	logPodPrefix = "distribution:pod"
)

type Pod struct {
	context context.Context
	storage storage.Storage
}

func (p *Pod) Runtime() (*types.Runtime, error) {

	log.V(logLevel).Debugf("%s:get:> get pod runtime info", logPodPrefix)
	runtime, err := p.storage.Info(p.context, p.storage.Collection().Pod(), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> get runtime info error: %s", logPodPrefix, err)
		return &runtime.Runtime, err
	}
	return &runtime.Runtime, nil
}

// Get pod info from storage
func (p *Pod) Get(namespace, service, deployment, name string) (*types.Pod, error) {
	log.V(logLevel).Debugf("%s:get:> get by name %s", logPodPrefix, name)

	pod := new(types.Pod)

	err := p.storage.Get(p.context, p.storage.Collection().Pod(),
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
		s.Labels["LBC"] = pod.SelfLink()
		s.DNS = types.SpecTemplateContainerDNS{
			Server: ips,
			Search: ips,
		}
		pod.Spec.Template.Containers = append(pod.Spec.Template.Containers, s)
	}

	for _, s := range deployment.Spec.Template.Volumes {
		pod.Spec.Template.Volumes = append(pod.Spec.Template.Volumes, s)
	}

	if err := p.storage.Put(p.context, p.storage.Collection().Pod(),
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

	err := p.storage.List(p.context, p.storage.Collection().Pod(), filter, list, nil)
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

	err := p.storage.List(p.context, p.storage.Collection().Pod(), filter, list, nil)
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

	err := p.storage.List(p.context, p.storage.Collection().Pod(), filter, list, nil)
	if err != nil {
		log.V(logLevel).Debugf("%s:listbydeployment:> get pod list by deployment id `%s/%s/%s` err: %v",
			logPodPrefix, namespace, service, deployment, err)
		return nil, err
	}

	return list, nil
}

// SetNode - set node info to pod
func (p *Pod) SetNode(pod *types.Pod, node *types.Node) error {
	log.V(logLevel).Debugf("%s:setnode:> set node for pod: %s", logPodPrefix, pod.Meta.Name)

	pod.Meta.Node = node.Meta.Name

	if err := p.storage.Set(p.context, p.storage.Collection().Pod(),
		p.storage.Key().Pod(pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name), pod, nil); err != nil {
		log.Errorf("%s:setnode:> pod set node err: %v", logPodPrefix, err)
		return err
	}

	return nil
}

// SetStatus - set state for pod
func (p *Pod) Update(pod *types.Pod) error {

	log.V(logLevel).Debugf("%s:update:> update pod: %s", logPodPrefix, pod.Meta.Name)

	if err := p.storage.Set(p.context, p.storage.Collection().Pod(),
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

	if err := p.storage.Set(p.context, p.storage.Collection().Pod(),
		p.storage.Key().Pod(pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name), pod, nil); err != nil {
		log.Errorf("%s:destroy:> mark pod for destroy error: %v", logPodPrefix, err)
		return err
	}
	return nil
}

// Remove pod from storage
func (p *Pod) Remove(pod *types.Pod) error {
	if err := p.storage.Del(p.context, p.storage.Collection().Pod(),
		p.storage.Key().Pod(pod.Meta.Namespace, pod.Meta.Service, pod.Meta.Deployment, pod.Meta.Name)); err != nil {
		log.Errorf("%s:remove:> mark pod for destroy error: %v", logPodPrefix, err)
		return err
	}
	return nil
}

func (p *Pod) Watch(ch chan types.PodEvent, rev *int64) error {
	log.V(logLevel).Debugf("%s:watch:> watch pod, from revision %d", logPodPrefix, *rev)

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

				if err := json.Unmarshal(e.Data.([]byte), obj); err != nil {
					log.Errorf("%s:watch:> parse json", logPodPrefix)
					continue
				}

				res.Data = obj

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := p.storage.Watch(p.context, p.storage.Collection().Pod(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func (p *Pod) ManifestMap(node string) (*types.PodManifestMap, error) {
	log.V(logLevel).Debugf("%s:PodManifestMap:> ", logPodPrefix)

	var (
		mf = types.NewPodManifestMap()
	)

	if err := p.storage.Map(p.context, p.storage.Collection().Manifest().Pod(node), types.EmptyString, mf, nil); err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s:PodManifestMap:> err: %s", logPodPrefix, err.Error())
			return nil, err
		}

		return nil, nil
	}

	return mf, nil
}

func (p *Pod) ManifestGet(node, pod string) (*types.PodManifest, error) {
	log.V(logLevel).Debugf("%s:PodManifestGet:> ", logPodPrefix)

	var (
		mf = new(types.PodManifest)
	)

	if err := p.storage.Get(p.context, p.storage.Collection().Manifest().Pod(node), pod, &mf, nil); err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return mf, nil
}

func (p *Pod) ManifestAdd(node, pod string, manifest *types.PodManifest) error {
	log.V(logLevel).Debugf("%s:PodManifestAdd:> ", logPodPrefix)

	if err := p.storage.Put(p.context, p.storage.Collection().Manifest().Pod(node), pod, manifest, nil); err != nil {
		log.Errorf("%s:PodManifestAdd:> err :%s", logPodPrefix, err.Error())
		return err
	}

	return nil
}

func (p *Pod) ManifestSet(node, pod string, manifest *types.PodManifest) error {
	log.V(logLevel).Debugf("%s:PodManifestSet:> ", logPodPrefix)

	if err := p.storage.Set(p.context, p.storage.Collection().Manifest().Pod(node), pod, manifest, nil); err != nil {
		log.Errorf("%s:PodManifestSet:> err :%s", logPodPrefix, err.Error())
		return err
	}

	return nil
}

func (p *Pod) ManifestDel(node, pod string) error {
	log.V(logLevel).Debugf("%s:PodManifestDel:> %s on node %s", logPodPrefix, pod, node)

	if err := p.storage.Del(p.context, p.storage.Collection().Manifest().Pod(node), pod); err != nil {
		log.Errorf("%s:PodManifestDel:> err :%s", logPodPrefix, err.Error())
		return err
	}

	return nil
}

func (p *Pod) ManifestWatch(node string, ch chan types.PodManifestEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch pod manifest ", logPodPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	var f, c string

	if node != types.EmptyString {
		f = fmt.Sprintf(`\b.+\/%s\/%s\/(.+)\b`, node, storage.PodKind)
		c = p.storage.Collection().Manifest().Pod(node)
	} else {
		f = fmt.Sprintf(`\b.+\/(.+)\/%s\/(.+)\b`, storage.PodKind)
		c = p.storage.Collection().Manifest().Node()
	}

	r, err := regexp.Compile(f)
	if err != nil {
		log.Errorf("%s:> filter compile err: %v", logPodPrefix, err.Error())
		return err
	}

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

				keys := r.FindStringSubmatch(e.System.Key)
				if len(keys) == 0 {
					continue
				}

				res := types.PodManifestEvent{}
				res.Action = e.Action
				res.Name = e.Name
				res.SelfLink = e.SelfLink
				if node != types.EmptyString {
					res.Node = node
				} else {
					res.Node = keys[1]
				}

				manifest := new(types.PodManifest)

				if err := json.Unmarshal(e.Data.([]byte), manifest); err != nil {
					log.Errorf("%s:> parse data err: %v", logPodPrefix, err)
					continue
				}

				res.Data = manifest

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := p.storage.Watch(p.context, c, watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewPodModel(ctx context.Context, stg storage.Storage) *Pod {
	return &Pod{ctx, stg}
}
