//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package image

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/build"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
)

var Util IUtil = util{}

func Create(ctx context.Context, registry string, source *types.ServiceSource) (*types.Image, error) {

	var (
		log     = c.Get().GetLogger()
		storage = c.Get().GetStorage()
	)

	log.Debug("Create image")

	name := Util.Name(ctx, registry, source.Repo)
	isource := &types.ImageSource{
		Hub:   source.Hub,
		Owner: source.Owner,
		Repo:  source.Repo,
		Tag:   source.Branch,
	}

	img, err := storage.Image().Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if img == nil {
		img, err = storage.Image().Insert(ctx, name, isource)
		if err != nil {
			return nil, err
		}
	}

	_, err = build.Create(ctx, img.Meta.ID, source)
	if err != nil {
		return nil, err
	}

	auth := Util.RegistryAuth(ctx, name)
	if auth != nil {
		img.Registry.Auth = new(types.RegistryAuth)
		img.Registry.Auth.Server = auth.Server
		img.Registry.Auth.Username = auth.Username
		img.Registry.Auth.Password = auth.Password
	}

	return img, nil
}
