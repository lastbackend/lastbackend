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

package distribution

import (
	"context"

	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"regexp"
)

const (
	logVolumePrefix = "distribution:volume"
)

type Volume struct {
	context context.Context
	storage storage.Storage
}

func (v *Volume) Runtime() (*types.Runtime, error) {

	log.V(logLevel).Debugf("%s:get:> get services runtime info", logVolumePrefix)
	runtime, err := v.storage.Info(v.context, v.storage.Collection().Volume(), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> get runtime info error: %s", logVolumePrefix, err)
		return &runtime.Runtime, err
	}
	return &runtime.Runtime, nil
}

func (v *Volume) Get(namespace, name string) (*types.Volume, error) {
	log.V(logLevel).Debugf("%s:get:> get volume by id %s/%s", logVolumePrefix, namespace, name)

	item := new(types.Volume)

	err := v.storage.Get(v.context, v.storage.Collection().Volume(), v.storage.Key().Volume(namespace, name), &item, nil)
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> in namespace %s by name %s not found", logVolumePrefix, namespace, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> in namespace %s by name %s error: %v", logVolumePrefix, namespace, name, err)
		return nil, err
	}

	return item, nil
}

func (v *Volume) ListByNamespace(namespace string) (*types.VolumeList, error) {
	log.V(logLevel).Debugf("%s:list:> get volumes list", logVolumePrefix)

	list := types.NewVolumeList()
	filter := v.storage.Filter().Volume().ByNamespace(namespace)
	err := v.storage.List(v.context, v.storage.Collection().Volume(), filter, list, nil)
	if err != nil {
		log.V(logLevel).Error("%s:list:> get volumes list err: %v", logVolumePrefix, err)
		return list, err
	}

	log.V(logLevel).Debugf("%s:list:> get volumes list result: %d", logVolumePrefix, len(list.Items))

	return list, nil
}

func (v *Volume) Create(namespace *types.Namespace, vol *types.Volume) (*types.Volume, error) {
	log.V(logLevel).Debugf("%s:crete:> create volume %s", logVolumePrefix, vol.SelfLink())

	vol.Meta.SetDefault()
	vol.Meta.Namespace = namespace.Meta.Name
	vol.Status.State = types.StateCreated
	vol.SelfLink()

	if err := v.storage.Put(v.context, v.storage.Collection().Volume(),
		v.storage.Key().Volume(vol.Meta.Namespace, vol.Meta.Name), vol, nil); err != nil {
		log.V(logLevel).Errorf("%s:crete:> insert volume err: %v", logVolumePrefix, err)
		return nil, err
	}

	return vol, nil
}

func (v *Volume) Update(volume *types.Volume) error {
	log.V(logLevel).Debugf("%s:update:> update volume %s", logVolumePrefix, volume.Meta.Name)

	if err := v.storage.Set(v.context, v.storage.Collection().Volume(),
		v.storage.Key().Volume(volume.Meta.Namespace, volume.Meta.Name), volume, nil); err != nil {
		log.V(logLevel).Errorf("%s:update:> update volume err: %v", logVolumePrefix, err)
		return err
	}

	return nil
}

func (v *Volume) Destroy(volume *types.Volume) error {

	if volume == nil {
		log.V(logLevel).Warnf("%s:destroy:> invalid argument %v", logVolumePrefix, volume)
		return nil
	}

	log.V(logLevel).Debugf("%s:destroy:> volume %s", logVolumePrefix, volume.Meta.Name)

	volume.Status.State = types.StateDestroy
	volume.Spec.State.Destroy = true

	if err := v.storage.Set(v.context, v.storage.Collection().Volume(),
		v.storage.Key().Volume(volume.Meta.Namespace, volume.Meta.Name), volume, nil); err != nil {
		log.Errorf("%s:destroy:> volume err: %v", logVolumePrefix, err)
		return err
	}

	return nil
}

func (v *Volume) Remove(volume *types.Volume) error {
	log.V(logLevel).Debugf("%s:remove:> remove volume %#v", logVolumePrefix, volume)

	if err := v.storage.Del(v.context, v.storage.Collection().Volume(),
		v.storage.Key().Volume(volume.Meta.Namespace, volume.Meta.Name)); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove volume  err: %v", logVolumePrefix, err)
		return err
	}

	return nil
}

// Watch service changes
func (v *Volume) Watch(ch chan types.VolumeEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch volume by spec changes", logVolumePrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-v.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.VolumeEvent{}
				res.Action = e.Action
				res.Name = e.Name

				volume := new(types.Volume)

				if err := json.Unmarshal(e.Data.([]byte), volume); err != nil {
					log.Errorf("%s:> parse data err: %v", logServicePrefix, err)
					continue
				}

				res.Data = volume

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := v.storage.Watch(v.context, v.storage.Collection().Volume(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func (v *Volume) ManifestMap(node string) (*types.VolumeManifestMap, error) {
	log.V(logLevel).Debugf("%s:VolumeManifestMap:> ", logVolumePrefix)

	var (
		mf = types.NewVolumeManifestMap()
	)

	if err := v.storage.Map(v.context, v.storage.Collection().Manifest().Volume(node), types.EmptyString, mf, nil); err != nil {
		log.Errorf("%s:VolumeManifestMap:> err :%s", logVolumePrefix, err.Error())
		return nil, err
	}
	return mf, nil
}

func (v *Volume) ManifestGet(node, volume string) (*types.VolumeManifest, error) {
	log.V(logLevel).Debugf("%s:VolumeManifestGet:> ", logVolumePrefix)

	var (
		mf = new(types.VolumeManifest)
	)

	if err := v.storage.Get(v.context, v.storage.Collection().Manifest().Volume(node), volume, &mf, nil); err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		log.Errorf("%s:VolumeManifestGet:> err :%s", logVolumePrefix, err.Error())
		return nil, err
	}

	return mf, nil
}

func (v *Volume) ManifestAdd(node, volume string, manifest *types.VolumeManifest) error {
	log.V(logLevel).Debugf("%s:VolumeManifestAdd:> ", logVolumePrefix)

	if err := v.storage.Put(v.context, v.storage.Collection().Manifest().Volume(node), volume, manifest, nil); err != nil {
		log.Errorf("%s:VolumeManifestAdd:> err :%s", logVolumePrefix, err.Error())
		return err
	}

	return nil
}

func (v *Volume) ManifestSet(node, volume string, manifest *types.VolumeManifest) error {
	log.V(logLevel).Debugf("%s:VolumeManifestSet:> ", logVolumePrefix)

	if err := v.storage.Set(v.context, v.storage.Collection().Manifest().Volume(node), volume, manifest, nil); err != nil {
		log.Errorf("%s:VolumeManifestSet:> err :%s", logVolumePrefix, err.Error())
		return err
	}

	return nil
}

func (v *Volume) ManifestDel(node, volume string) error {
	log.V(logLevel).Debugf("%s:DelVolumeManifest:> ", logVolumePrefix)

	if err := v.storage.Del(v.context, v.storage.Collection().Manifest().Volume(node), volume); err != nil {
		log.Errorf("%s:VolumeManifestDel:> err :%s", logVolumePrefix, err.Error())
		return err
	}

	return nil
}

func (p *Volume) ManifestWatch(node string, ch chan types.VolumeManifestEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch volume manifest ", logVolumePrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	var f, c string

	if node != types.EmptyString {
		f = fmt.Sprintf(`\b.+\/%s\/%s\/(.+)\b`, node, storage.VolumeKind)
		c = p.storage.Collection().Manifest().Volume(node)
	} else {
		f = fmt.Sprintf(`\b.+\/(.+)\/%s\/(.+)\b`, storage.VolumeKind)
		c = p.storage.Collection().Manifest().Node()
	}

	r, err := regexp.Compile(f)
	if err != nil {
		log.Errorf("%s:> filter compile err: %v", logVolumePrefix, err.Error())
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

				res := types.VolumeManifestEvent{}
				res.Action = e.Action
				res.Name = e.Name
				res.SelfLink = e.SelfLink
				if node != types.EmptyString {
					res.Node = node
				} else {
					res.Node = keys[1]
				}

				manifest := new(types.VolumeManifest)

				if err := json.Unmarshal(e.Data.([]byte), manifest); err != nil {
					log.Errorf("%s:> parse data err: %v", logVolumePrefix, err)
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

func NewVolumeModel(ctx context.Context, stg storage.Storage) *Volume {
	return &Volume{ctx, stg}
}
