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
	"github.com/lastbackend/lastbackend/pkg/api/build"
	c "github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
)

const logLevel = 3

var Util IUtil = util{}

func Create(ctx context.Context, registry string, source types.ServiceSource) (*types.Image, error) {

	var (
		log     = c.Get().GetLogger()
		storage = c.Get().GetStorage()
		image   = types.Image{}
	)

	log.V(logLevel).Debugf("Image: create new image with registry: %s and source: %#v", registry, source)

	image.Meta.SetDefault()
	image.Meta.Name = Util.Name(ctx, registry, source.Repo)

	image.Source = types.ImageSource{
		Hub:   source.Hub,
		Owner: source.Owner,
		Repo:  source.Repo,
		Tag:   source.Branch,
	}

	if err := storage.Image().Insert(ctx, &image); err != nil {
		log.V(logLevel).Debugf("Image: insert image err: %s", err.Error())
		return &image, err
	}

	if _, err := build.Create(ctx, image.Meta.Name, &image.Source); err != nil {
		log.V(logLevel).Debugf("Image: create build err: %s", err.Error())
		return &image, err
	}

	auth := Util.RegistryAuth(ctx, image.Meta.Name)
	if auth != nil {
		image.Registry.Auth = types.RegistryAuth{
			Server:   auth.Server,
			Username: auth.Username,
			Password: auth.Password,
		}
	}

	return &image, nil
}

func Get(ctx context.Context, registry string, source types.ServiceSource) (*types.Image, error) {

	var (
		log     = c.Get().GetLogger()
		storage = c.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Image: get image by registry: %s and source: %#v", registry, source)

	name := Util.Name(ctx, registry, source.Repo)

	image, err := storage.Image().Get(ctx, name)
	if err != nil {
		log.V(logLevel).Debugf("Image: create build err: %s", err.Error())
		return nil, err
	}

	auth := Util.RegistryAuth(ctx, image.Meta.Name)
	if auth != nil {
		image.Registry.Auth = types.RegistryAuth{
			Server:   auth.Server,
			Username: auth.Username,
			Password: auth.Password,
		}
	}

	return image, nil
}
