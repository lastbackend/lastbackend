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
	"strings"
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
	for _, image := range images {
		for _, tag := range image.Meta.Tags {
			s.images[tag] = image
		}

	}
}

func (s *ImageState) GetImage(tag string) *types.Image {
	log.V(logLevel).Debugf("Cache: ImageCache: get image: %s", tag)
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(strings.Split(tag, ":"))==1 {
		tag+=":latest"
	}

	image, ok := s.images[tag]
	if !ok {
		return nil
	}
	return image
}

func (s *ImageState) AddImage(tag string, image *types.Image) {
	log.V(logLevel).Debugf("Cache: ImageCache: add image: %#v", image)
	s.SetImage(tag, image)
}

func (s *ImageState) SetImage(tag string, image *types.Image) {
	log.V(logLevel).Debugf("Cache: ImageCache: set image: %#v", image)
	s.lock.Lock()
	defer s.lock.Unlock()

	s.images[image.SelfLink()] = image

	for _, i := range image.Meta.Tags {
		s.images[i] = image
	}
}

func (s *ImageState) DelImage(link string) {
	log.V(logLevel).Debugf("Cache: ImageCache: del image: %s", link)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.images[link]; ok {
		delete(s.images, link)
	}
}
