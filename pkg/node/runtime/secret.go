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
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)

func SecretGet(ctx context.Context, name string) (*types.Secret, error) {
	secret := envs.Get().GetState().Secrets().GetSecret(name)
	if secret != nil {
		return secret, nil
	}

	sr, err := envs.Get().GetRestClient().Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("can not receive secret from api, err: %s", err.Error())
		return nil, err
	}

	return sr.Decode(), nil
}

func SecretCreate(ctx context.Context, name string) error {

	ok := envs.Get().GetState().Secrets().GetSecret(name)
	if ok != nil {
		return nil
	}

	secret, err := envs.Get().GetRestClient().Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("get secret err: %s", err.Error())
		return err
	}

	envs.Get().GetState().Secrets().AddSecret(secret.Meta.Name, secret.Decode())
	return nil
}

func SecretUpdate(ctx context.Context, name string) error {

	secret, err := envs.Get().GetRestClient().Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("get secret err: %s", err.Error())
		return err
	}

	envs.Get().GetState().Secrets().AddSecret(secret.Meta.Name, secret.Decode())
	return nil

}

func SecretRemove(ctx context.Context, name string) {
	envs.Get().GetState().Secrets().DelSecret(name)
}
