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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)


func VolumeManage(ctx context.Context, key string, manifest *types.VolumeManifest) error {

	log.V(logLevel).Debugf("Provision volume: %s", key)

	//==========================================================================
	// Destroy pod =============================================================
	//==========================================================================

	// Call destroy pod
	if manifest.State.Destroy {

		v := envs.Get().GetState().Volumes().GetVolume(key)
		if v == nil {

			vs := types.NewVolumeStatus()
			vs.SetDestroyed()
			envs.Get().GetState().Volumes().AddVolume(key, vs)

			return nil
		}

		log.V(logLevel).Debugf("Volume found > destroy it: %s", key)

		VolumeDestroy(ctx, key)
		v.SetDestroyed()
		envs.Get().GetState().Volumes().DelVolume(key)
		return nil
	}

	//==========================================================================
	// Check containers pod status =============================================
	//==========================================================================

	// Get pod list from current state
	v := envs.Get().GetState().Volumes().GetVolume(key)
	if v != nil {
		return nil
	}

	log.V(logLevel).Debugf("Volume not found > create it: %s", key)

	status, err := VolumeCreate(ctx, key, manifest)
	if err != nil {
		log.Errorf("Can not create pod: %s err: %s", key, err.Error())
		status.SetError(err)
	}

	envs.Get().GetState().Volumes().SetVolume(key, status)
	return nil
}

func VolumeCreate(ctx context.Context, name string, mf *types.VolumeManifest) (*types.VolumeStatus, error) {

	var status   = new(types.VolumeStatus)

	log.V(logLevel).Debugf("Create volume: %s", mf)
	if mf.Type == types.EmptyString {
		mf.Type = types.KindVolumeHostDir
	}

	si, err := envs.Get().GetCSI(mf.Type)
	if err != nil {
		log.Errorf("Can-not get storage interface: %s", err)
		return nil, err
	}

	st, err := si.Create(ctx, name, mf)
	if err != nil {
		log.Errorf("Can-not get secret from api: %s", err)
		return nil, err
	}

	if st.Ready {
		status.State = types.StateReady
	}

	status.Status = *st
	envs.Get().GetState().Volumes().AddVolume(name, status)

	return status, nil
}


func VolumeDestroy(ctx context.Context, name string) error {

	vol := envs.Get().GetState().Volumes().GetVolume(name)

	if vol == nil {
		return nil
	}

	if vol.Status.Type == types.EmptyString {
		vol.Status.Type = types.KindVolumeHostDir
	}

	si, err := envs.Get().GetCSI(vol.Status.Type)
	if err != nil {
		log.Errorf("Remove volume failed: %s", err.Error())
		return err
	}


	if err := si.Remove(ctx, &vol.Status); err != nil {
		log.Warnf("can note remove volume: %s: %s", name, err.Error())
	}

	envs.Get().GetState().Volumes().DelVolume(name)

	return nil
}

func VolumeRestore(ctx context.Context) error {

	log.Debug("Start volumes restore")

	tp := envs.Get().ListCSI()

	for _, t := range tp {

		log.Debugf("restore volumes type: %s", t)
		sci, err := envs.Get().GetCSI(t)
		if err != nil {
			log.Errorf("storage interface init err: %s", err.Error())
			return err
		}

		if sci == nil {
			return errors.New("container storage runtime interface not supported")
		}

		states, err := sci.List(ctx)
		if err != nil {
			log.Errorf("volumes restore err: %s", err.Error())
			return err
		}

		for name, state := range states {
			status := new(types.VolumeStatus)
			if state.Ready {
				status.State = types.StateReady
			}
			status.Status = *state

			envs.Get().GetState().Volumes().SetVolume(name, status)
		}

	}

	return nil
}

func VolumeSetSecretData (ctx context.Context, name string, secret string) error {
	return nil
}

func VolumeCheckSecretData(ctx context.Context, name string, secret string) (bool, error) {
	log.Debugf("volume check secret data: %s > %s", secret, name)
	return true, nil
}

func VolumeCheckConfigData(ctx context.Context, name string, config string) (bool, error) {
	log.Debugf("volume check config data: %s > %s", config, name)

	vol := envs.Get().GetState().Volumes().GetVolume(name)
	cfg := envs.Get().GetState().Configs().GetConfig(config)

	if vol == nil {
		return  false, errors.New("volume not exists")
	}

	if vol.Status.Type == types.EmptyString {
		vol.Status.Type = types.KindVolumeHostDir
	}

	si, err := envs.Get().GetCSI(vol.Status.Type)
	if err != nil {
		log.Errorf("Remove volume failed: %s", err.Error())
		return false, err
	}

	return si.FilesCheck(ctx, &vol.Status, cfg.Data)
}

func VolumeSetConfigData (ctx context.Context, name string, config string) error {

	log.Debugf("set volume config data: %s > %s", config, name)

	vol := envs.Get().GetState().Volumes().GetVolume(name)
	cfg := envs.Get().GetState().Configs().GetConfig(config)

	if vol == nil {
		return errors.New("volume not exists")
	}

	if vol.Status.Type == types.EmptyString {
		vol.Status.Type = types.KindVolumeHostDir
	}

	si, err := envs.Get().GetCSI(vol.Status.Type)
	if err != nil {
		log.Errorf("Remove volume failed: %s", err.Error())
		return err
	}

	return si.FilesPut(ctx, &vol.Status, cfg.Data)
}
