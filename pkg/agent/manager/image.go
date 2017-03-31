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

package manager

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

type ImageManager struct {
	update chan types.ImageList
	close  chan bool
}

func NewImageManager() *ImageManager {
	ctx := context.Get()
	ctx.Log.Info("start image manager")
	var im = new(ImageManager)

	im.update = make(chan types.ImageList)
	im.close = make(chan bool)

	return im
}

func ReleaseImageManager(im *ImageManager) error {
	ctx := context.Get()
	ctx.Log.Info("release image manager")
	close(im.update)
	close(im.close)
	return nil
}

func (im *ImageManager) watch() error {
	ctx := context.Get()
	ctx.Log.Info("start image watcher")

	for {
		select {
		case _ = <-im.close:
			return ReleaseImageManager(im)
		case images := <-im.update:
			{
				for _, image := range images {
					im.sync(&image)
				}
			}
		}
	}

	return nil
}

func (im *ImageManager) sync(i *types.Image) {
	ctx := context.Get()
	ctx.Log.Info("image manager sync")

}
