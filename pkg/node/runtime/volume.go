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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)

func VolumeCreate(ctx context.Context, name string, mf *types.VolumeManifest) (*types.VolumeState, error) {

	log.V(logLevel).Debugf("Create volume: %s", mf)
	if mf.Type == types.EmptyString {
		mf.Type = types.VOLUMETYPELOCAL
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

	envs.Get().GetState().Volumes().AddVolume(name, st)

	return st, nil
}

func VolumeDestroy(ctx context.Context, name string) error {

	vol := envs.Get().GetState().Volumes().GetVolume(name)

	if vol == nil {
		return nil
	}

	if vol.Type == types.EmptyString {
		vol.Type = types.VOLUMETYPELOCAL
	}

	si, err := envs.Get().GetCSI(vol.Type)
	if err != nil {
		log.Errorf("Remove volume failed: %s", err.Error())
		return err
	}

	mf := types.VolumeManifest{
		Type: vol.Type,
		Path: vol.Path,
	}

	if err := si.Remove(ctx, name, &mf); err != nil {
		log.Warnf("can note remove volume: %s: %s", name, err.Error())
	}

	envs.Get().GetState().Volumes().DelVolume(name)

	return nil
}

func VolumeUpdate(ctx context.Context, name string, manifest *types.VolumeManifest) (*types.VolumeState, error) {
	return nil, nil
}

func VolumeRestore(ctx context.Context) error {
	return nil
}
