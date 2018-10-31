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

const logImagePrefix = "state:images:>"

type ImageState struct {
	lock   sync.RWMutex
	images map[string]*types.Image
}

func (s *ImageState) GetImages() map[string]*types.Image {
	log.V(logLevel).Debugf("%s get images", logImagePrefix)
	return s.images
}

func (s *ImageState) SetImages(images map[string]*types.Image) {
	log.V(logLevel).Debugf("%s set images: %d", logImagePrefix, len(images))
	for _, image := range images {
		for _, tag := range image.Meta.Tags {
			s.images[tag] = image
		}

	}
}

func (s *ImageState) GetImage(tag string) *types.Image {
	log.V(logLevel).Debugf("%s get image: %s", logImagePrefix, tag)
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(strings.Split(tag, ":")) == 1 {
		tag += ":latest"
	}

	image, ok := s.images[tag]
	if !ok {
		return nil
	}
	return image
}

func (s *ImageState) AddImage(tag string, image *types.Image) {
	log.V(logLevel).Debugf("%s add image: %s", logImagePrefix, image.Meta.Name)
	s.SetImage(tag, image)
}

func (s *ImageState) SetImage(tag string, image *types.Image) {
	log.V(logLevel).Debugf("%s set image: %s", logImagePrefix, image.Meta.Name)
	s.lock.Lock()
	defer s.lock.Unlock()

	s.images[image.SelfLink()] = image

	for _, i := range image.Meta.Tags {
		s.images[i] = image
	}
}

func (s *ImageState) DelImage(link string) {
	log.V(logLevel).Debugf("%s del image: %s", logImagePrefix, link)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.images[link]; ok {
		delete(s.images, link)
	}
}
