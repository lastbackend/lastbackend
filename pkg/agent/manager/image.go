package manager

import (
	"github.com/golang/glog"
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
