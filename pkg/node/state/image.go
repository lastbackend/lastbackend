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

package state

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
)

type ImageState struct {
	lock   sync.RWMutex
	images map[string]*types.Image
}

func (s *ImageState) GetImages() map[string]*types.Image {
	log.V(logLevel).Debug("Cache: ImageCache: get pods")
	return s.images
}

func (s *ImageState) SetImages(images map[string]*types.Image) {
	log.V(logLevel).Debugf("Cache: ImageCache: set images: %#v", images)
	for h, image := range images {
		s.images[h] = image
	}
}

func (s *ImageState) GetImage(link string) *types.Image {
	log.V(logLevel).Debugf("Cache: ImageCache: get image: %s", link)
	s.lock.Lock()
	defer s.lock.Unlock()
	image, ok := s.images[link]
	if !ok {
		return nil
	}
	return image
}

func (s *ImageState) AddImage(link string, image *types.Image) {
	log.V(logLevel).Debugf("Cache: ImageCache: add image: %#v", image)
	s.SetImage(link, image)
}

func (s *ImageState) SetImage(link string, image *types.Image) {
	log.V(logLevel).Debugf("Cache: ImageCache: set image: %#v", image)
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.images[image.SelfLink()]; ok {
		delete(s.images, image.SelfLink())
	}

	s.images[image.SelfLink()] = image
}

func (s *ImageState) DelImage(link string) {
	log.V(logLevel).Debugf("Cache: ImageCache: del image: %s", link)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.images[link]; ok {
		delete(s.images, link)
	}
}
