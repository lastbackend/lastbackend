package volume

import (
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
)

type Volume struct {
	config *v1.PersistentVolume
}

type VolumeList []Volume

func Get(name string) (*v1.PersistentVolume, *e.Err) {

	var (
		er  error
		ctx = context.Get()
	)

	volume, er := ctx.K8S.Core().PersistentVolumes().Get(name)
	if er != nil {
		return nil, e.New("volume").Unknown(er)
	}

	return volume, nil
}

func List() (*v1.PersistentVolumeList, *e.Err) {

	var (
		er  error
		ctx = context.Get()
	)

	pv, er := ctx.K8S.Core().PersistentVolumes().List(api.ListOptions{})
	if er != nil {
		return nil, e.New("volume").Unknown(er)
	}

	return pv, nil
}

func Create(user, project string, config *v1.PersistentVolume) (*Volume, *e.Err) {

	var (
		ctx    = context.Get()
		volume = new(model.Volume)
		pv     = new(Volume)
	)

	volume.User = user
	volume.Project = project
	volume.Name = fmt.Sprintf("%s-%s", config.Name, generator.GetUUIDV4()[0:12])

	volume, err := ctx.Storage.Volume().Insert(volume)
	if err != nil {
		return nil, err
	}

	pv.config = config
	pv.config.Name = volume.Name

	return pv, nil
}

func (v *Volume) Deploy() *e.Err {

	var (
		er  error
		ctx = context.Get()
	)

	_, er = ctx.K8S.Core().PersistentVolumes().Create(v.config)
	if er != nil {
		return e.New("volume").Unknown(er)
	}

	return nil
}
