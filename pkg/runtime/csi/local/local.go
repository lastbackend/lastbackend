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

package local

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/runtime/csi"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Storage struct {
	csi.CSI
	root string
}

type StorageOpts struct {
	root string
}

func (s *Storage) List(ctx context.Context) (map[string]*types.VolumeState, error) {
	var vols = make(map[string]*types.VolumeState, 0)

	var dirs []string
	f, err := os.Open(s.root)
	if err != nil {
		return vols, err
	}

	items, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return vols, err
	}

	for _, item := range items {
		if item.IsDir() {
			dirs = append(dirs, item.Name())
		}
	}

	for _, dir := range dirs {
		vol := new(types.VolumeState)

		vol.Path = filepath.Join(s.root, dir)
		vol.Type = types.VOLUMETYPELOCAL
		vol.Ready = true
		vols[dir] = vol
	}

	return vols, nil
}

func (s *Storage) Create(ctx context.Context, name string, manifest *types.VolumeManifest) (*types.VolumeState, error) {

	var (
		status = new(types.VolumeState)
		path   = filepath.Join(s.root, manifest.HostPath)
	)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return status, err
		}
	}

	status.Path = path
	status.Ready = true

	return status, nil
}

func (s *Storage) Remove(ctx context.Context, name string, manifest *types.VolumeManifest) error {
	return os.Remove(filepath.Join(s.root, manifest.HostPath))
}

func Get() (*Storage, error) {

	log.Debug("Initialize local storage interface")
	var s = new(Storage)

	if viper.GetString("runtime.csi.local.root") != "" {
		s.root = viper.GetString("runtime.csi.local.root")
		log.Debugf("Initialize local storage interface root: %s", s.root)
	}

	if _, err := os.Stat(s.root); os.IsNotExist(err) {
		err = os.MkdirAll(s.root, 0755)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
