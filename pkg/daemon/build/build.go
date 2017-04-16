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

package build

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
)

func Create(ctx context.Context, imageName string, source *types.ImageSource) (*types.Build, error) {
	var (
		log     = c.Get().GetLogger()
		storage = c.Get().GetStorage()
	)

	log.Debug("Create build")

	build := types.Build{
		Source: types.BuildSource{
			Hub:   source.Hub,
			Owner: source.Owner,
			Repo:  source.Repo,
			Tag:   source.Tag,
		},
	}

	if err := storage.Build().Insert(ctx, imageName, &build); err != nil {
		return &build, err
	}

	return &build, nil
}
