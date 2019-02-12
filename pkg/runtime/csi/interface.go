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

package csi

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type CSI interface {
	List(ctx context.Context) (map[string]*types.VolumeState, error)
	Create(ctx context.Context, name string, manifest *types.VolumeManifest) (*types.VolumeState, error)
	FilesPut(ctx context.Context, state *types.VolumeState, files map[string]string) error
	FilesCheck(ctx context.Context, state *types.VolumeState, files map[string]string) (bool, error)
	FilesDel(ctx context.Context, state *types.VolumeState, files []string) error
	Remove(ctx context.Context, state *types.VolumeState) error
}
