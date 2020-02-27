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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/minion/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
	"strings"
)

const (
	logVolumePrefix = "node:runtime:volume:>"
)

func VolumeManage(ctx context.Context, key string, manifest *types.VolumeManifest) error {

	log.V(logLevel).Debugf("%s provision volume: %s", logVolumePrefix, key)

	//==========================================================================
	// Destroy volume ==========================================================
	//==========================================================================

	// Call destroy volume
	if manifest.State.Destroy {

		v := envs.Get().GetState().Volumes().GetVolume(key)
		if v == nil {

			vs := types.NewVolumeStatus()
			vs.SetDestroyed()
			envs.Get().GetState().Volumes().AddVolume(key, vs)

			return nil
		}

		log.V(logLevel).Debugf("%s volume found > destroy it: %s", logVolumePrefix, key)

		if err := VolumeDestroy(ctx, key); err != nil {
			log.Errorf("%s can not destroy volume: %s", logVolumePrefix, err.Error())
			return err
		}

		v.SetDestroyed()
		envs.Get().GetState().Volumes().SetVolume(key, v)
		return nil
	}

	//==========================================================================
	// Check containers volume status =============================================
	//==========================================================================

	// Get volume list from current state
	v := envs.Get().GetState().Volumes().GetVolume(key)
	if v != nil {
		if v.State != types.StateDestroyed {
			return nil
		}
	}

	log.V(logLevel).Debugf("%s volume not found > create it: %s", logVolumePrefix, key)

	status, err := VolumeCreate(ctx, key, manifest)
	if err != nil {
		log.Errorf("%s can not create volume: %s err: %s", logVolumePrefix, key, err.Error())
		status.SetError(err)
	}

	envs.Get().GetState().Volumes().SetVolume(key, status)
	return nil
}

func VolumeCreate(ctx context.Context, name string, mf *types.VolumeManifest) (*types.VolumeStatus, error) {

	var status = new(types.VolumeStatus)

	log.V(logLevel).Debugf("%s create volume: %s", logVolumePrefix, mf)
	if mf.Type == types.EmptyString {
		mf.Type = types.KindVolumeHostDir
	}

	si, err := envs.Get().GetCSI(mf.Type)
	if err != nil {
		log.Errorf("%s can-not get storage interface: %s", logVolumePrefix, err)
		return nil, err
	}

	st, err := si.Create(ctx, name, mf)
	if err != nil {
		log.Errorf("%s can-not get secret from api: %s", logVolumePrefix, err)
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

	log.V(logLevel).Debugf("%s destroy volume: %s", logVolumePrefix, name)

	vol := envs.Get().GetState().Volumes().GetVolume(name)

	if vol == nil {
		return nil
	}

	if vol.Status.Type == types.EmptyString {
		vol.Status.Type = types.KindVolumeHostDir
	}

	si, err := envs.Get().GetCSI(vol.Status.Type)
	if err != nil {
		log.Errorf("%s remove volume failed: %s", logVolumePrefix, err.Error())
		return err
	}

	if err := si.Remove(ctx, &vol.Status); err != nil {
		log.Warnf("%s can not remove volume: %s: %s", logVolumePrefix, name, err.Error())
	}

	vol.SetDestroyed()
	envs.Get().GetState().Volumes().SetVolume(name, vol)

	return nil
}

func VolumeRestore(ctx context.Context) error {

	log.Debugf("%s start volumes restore", logVolumePrefix)

	tp := envs.Get().ListCSI()

	for _, t := range tp {

		log.Debugf("%s restore volumes type: %s", logVolumePrefix, t)
		sci, err := envs.Get().GetCSI(t)
		if err != nil {
			log.Errorf("%s storage interface init err: %s", logVolumePrefix, err.Error())
			return err
		}

		if sci == nil {
			return errors.New("container storage runtime interface not supported")
		}

		states, err := sci.List(ctx)
		if err != nil {
			log.Errorf("%s volumes restore err: %s", logVolumePrefix, err.Error())
			return err
		}

		for name, state := range states {
			status := new(types.VolumeStatus)
			if state.Ready {
				status.State = types.StateReady
			}
			status.Status = *state
			envs.Get().GetState().Volumes().SetVolume(strings.Replace(name, "_", ":", -1), status)
		}

	}

	return nil
}

func VolumeSetSecretData(ctx context.Context, name string, secret string) error {
	log.Debugf("%s volume set secret data: %s > %s", logVolumePrefix, secret, name)
	return nil
}

func VolumeCheckSecretData(ctx context.Context, name string, secret string) (bool, error) {
	log.Debugf("%s volume check secret data: %s > %s", logVolumePrefix, secret, name)
	return true, nil
}

func VolumeCheckConfigData(ctx context.Context, name string, config string) (bool, error) {
	log.Debugf("%s volume check config data: %s > %s", logVolumePrefix, config, name)

	vol := envs.Get().GetState().Volumes().GetVolume(name)
	cfg := envs.Get().GetState().Configs().GetConfig(config)

	if vol == nil {
		return false, errors.New("volume not exists")
	}

	if vol.Status.Type == types.EmptyString {
		vol.Status.Type = types.KindVolumeHostDir
	}

	si, err := envs.Get().GetCSI(vol.Status.Type)
	if err != nil {
		log.Errorf("%s remove volume failed: %s", logVolumePrefix, err.Error())
		return false, err
	}

	return si.FilesCheck(ctx, &vol.Status, cfg.Data)
}

func VolumeSetConfigData(ctx context.Context, name string, config string) error {

	log.Debugf("%s set volume config data: %s > %s", logVolumePrefix, config, name)

	vol := envs.Get().GetState().Volumes().GetVolume(name)
	cfg := envs.Get().GetState().Configs().GetConfig(config)

	if vol == nil {
		return errors.New("volume not exists")
	}

	if cfg == nil {
		return errors.New("config not exists")
	}

	if vol.Status.Type == types.EmptyString {
		vol.Status.Type = types.KindVolumeHostDir
	}

	si, err := envs.Get().GetCSI(vol.Status.Type)
	if err != nil {
		log.Errorf("%s remove volume failed: %s", logVolumePrefix, err.Error())
		return err
	}

	return si.FilesPut(ctx, &vol.Status, cfg.Data)
}
