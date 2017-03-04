package volume

import (
	"fmt"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"k8s.io/client-go/pkg/api/v1"
)

type Volume struct {
	config *v1.PersistentVolume
}

type VolumeList []Volume

func Get(name string) (*v1.PersistentVolume, error) {

	var (
		err error
		ctx = context.Get()
	)

	volume, err := ctx.K8S.CoreV1().PersistentVolumes().Get(name)
	if err != nil {
		return nil, err
	}

	return volume, nil
}

func List() (*v1.PersistentVolumeList, error) {

	var (
		err error
		ctx = context.Get()
	)

	pv, err := ctx.K8S.CoreV1().PersistentVolumes().List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pv, nil
}

func Create(user, project string, config *v1.PersistentVolume) (*Volume, error) {

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

func (v *Volume) Deploy() error {

	var (
		err error
		ctx = context.Get()
	)

	_, err = ctx.K8S.CoreV1().PersistentVolumes().Create(v.config)
	if err != nil {
		return err
	}

	return nil
}
