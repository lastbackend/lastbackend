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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logManifestPrefix = "distribution:manifest"
)

type IManifest interface {
	GetNodeManifest(node string) (*types.NodeManifest, error)
	DelNodeManifest(node string) error

	AddPodManifest (node, pod string, manifest *types.PodManifest) error
	GetPodManifest (node, pod string) (*types.PodManifest, error)
	SetPodManifest (node, pod string, manifest *types.PodManifest) error
	DelPodManifest (node, pod string) error

	AddVolumeManifest (node, volume string, manifest *types.VolumeManifest) error
	GetVolumeManifest (node, volume string) (*types.VolumeManifest, error)
	SetVolumeManifest (node, volume string, manifest *types.VolumeManifest) error
	DetVolumeManifest (node, volume string) error
}

func GetNodeManifest(node string) (*types.NodeManifest, error) {
	log.Debugf("%s:GetNodeManifest:> ", logManifestPrefix)
	return new(types.NodeManifest), nil
}

func DelNodeManifest(node string) error {
	log.Debugf("%s:DelNodeManifest:> ", logManifestPrefix)
	return nil
}


func AddPodManifest (node, pod string, manifest *types.PodManifest) error {
	log.Debugf("%s:AddPodManifest:> ", logManifestPrefix)
	return nil
}


func GetPodManifest (node, pod string) (*types.PodManifest, error) {
	log.Debugf("%s:GetPodManifest:> ", logManifestPrefix)
	return new(types.PodManifest), nil
}


func SetPodManifest (node, pod string, manifest *types.PodManifest) error {
	log.Debugf("%s:SetPodManifest:> ", logManifestPrefix)
	return nil
}


func DelPodManifest (node, pod string) error {
	log.Debugf("%s:DelPodManifest:> ", logManifestPrefix)
	return nil
}


func AddVolumeManifest (node, volume string, manifest *types.VolumeManifest) error {
	log.Debugf("%s:AddVolumeManifest:> ", logManifestPrefix)
	return nil
}


func GetVolumeManifest (node, volume string) (*types.VolumeManifest, error) {
	log.Debugf("%s:GetVolumeManifest:> ", logManifestPrefix)
	return new(types.VolumeManifest), nil
}


func SetVolumeManifest (node, volume string, manifest *types.VolumeManifest) error {
	log.Debugf("%s:SetVolumeManifest:> ", logManifestPrefix)
	return nil
}


func DelVolumeManifest (node, volume string) error {
	log.Debugf("%s:DelVolumeManifest:> ", logManifestPrefix)
	return nil
}
