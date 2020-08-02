//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package volume
//
//import (
//	"context"
//	"net/http"
//
//	"github.com/lastbackend/lastbackend/internal/api/envs"
//	"github.com/lastbackend/lastbackend/internal/pkg/errors"
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
//	"github.com/lastbackend/lastbackend/tools/log"
//)
//
//const (
//	logPrefix = "api:handler:volume"
//	logLevel  = 3
//)
//
//func Fetch(ctx context.Context, namespace, name string) (*models.Volume, *errors.Err) {
//
//	vm := service.NewVolumeModel(ctx, envs.Get().GetStorage())
//	vol, err := vm.Get(namespace, name)
//
//	if err != nil {
//		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
//		return nil, errors.New("volume").InternalServerError(err)
//	}
//
//	if vol == nil {
//		err := errors.New("namespace not found")
//		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
//		return nil, errors.New("volume").NotFound()
//	}
//
//	return vol, nil
//}
//
//func Apply(ctx context.Context, ns *models.Namespace, mf *request.VolumeManifest) (*models.Volume, *errors.Err) {
//
//	if mf.Meta.Name == nil {
//		return nil, errors.New("volume").BadParameter("meta.name")
//	}
//
//	vol, err := Fetch(ctx, ns.Meta.Name, *mf.Meta.Name)
//	if err != nil {
//		if err.Code != http.StatusText(http.StatusNotFound) {
//			return nil, errors.New("volume").InternalServerError()
//		}
//	}
//
//	if vol == nil {
//		return Create(ctx, ns, mf)
//	}
//
//	return Update(ctx, ns, vol, mf)
//}
//
//func Create(ctx context.Context, ns *models.Namespace, mf *request.VolumeManifest) (*models.Volume, *errors.Err) {
//
//	vm := service.NewVolumeModel(ctx, envs.Get().GetStorage())
//	if mf.Meta.Name != nil {
//
//		srv, err := vm.Get(ns.Meta.Name, *mf.Meta.Name)
//		if err != nil {
//			log.Errorf("%s:create:> get volume by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
//			return nil, errors.New("volume").InternalServerError()
//		}
//
//		if srv != nil {
//			log.Warnf("%s:create:> volume name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
//			return nil, errors.New("volume").NotUnique("name")
//		}
//	}
//
//	vol := new(models.Volume)
//	vol.Meta.SetDefault()
//	vol.Meta.SelfLink = *models.NewVolumeSelfLink(ns.Meta.Name, *mf.Meta.Name)
//	vol.Meta.Namespace = ns.Meta.Name
//
//	mf.SetVolumeMeta(vol)
//	mf.SetVolumeSpec(vol)
//
//	if _, err := vm.Create(ns, vol); err != nil {
//		log.Errorf("%s:create:> create volume err: %s", logPrefix, ns.Meta.Name, err.Error())
//		return nil, errors.New("volume").InternalServerError()
//	}
//
//	return vol, nil
//}
//
////
//func Update(ctx context.Context, ns *models.Namespace, vol *models.Volume, mf *request.VolumeManifest) (*models.Volume, *errors.Err) {
//
//	vm := service.NewVolumeModel(ctx, envs.Get().GetStorage())
//
//	mf.SetVolumeMeta(vol)
//	mf.SetVolumeSpec(vol)
//
//	if err := vm.Update(vol); err != nil {
//		log.Errorf("%s:update:> update volume err: %s", logPrefix, err.Error())
//		return nil, errors.New("volume").InternalServerError()
//	}
//
//	return vol, nil
//}
