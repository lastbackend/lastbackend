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

package runtime

import (
	"encoding/base64"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"golang.org/x/net/context"
)

func imageCreateAuthString(username, password string) string {

	config := types.AuthConfig{
		Username: username,
		Password: password,
	}

	js, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(js)
}

func ImagePull(ctx context.Context, image *types.SpecTemplateContainerImage) error {

	var (
		mf = new(types.ImageManifest)
	)

	mf.Name = image.Name
	if image.Secret != types.EmptyString {
		secret, err := SecretGet(ctx, image.Secret)
		if err != nil {
			log.Errorf("can not get secret for image. err: %s", err.Error())
			return err
		}
		auth, err := secret.DecodeSecretAuthData()
		if err != nil {
			log.Errorf("can not get parse secret auth data. err: %s", err.Error())
			return err
		}
		mf.Auth = imageCreateAuthString(auth.Username, auth.Password)
	}

	img, err := envs.Get().GetIRI().Pull(ctx, mf)
	if err != nil {
		log.Errorf("can not pull image: %s", err.Error())
		return err
	}

	if img != nil {
		envs.Get().GetState().Images().AddImage(img.SelfLink(), img)
	}

	return nil
}

func ImageRemove(ctx context.Context, link string) error {
	if err := envs.Get().GetIRI().Remove(ctx, link); err != nil {
		log.Warnf("Can-not remove unnecessary image %s: %s", link, err)
	}

	envs.Get().GetState().Images().DelImage(link)

	return nil
}
