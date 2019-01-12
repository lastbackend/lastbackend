//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"strings"
)

func SecretGet(ctx context.Context, selflink string) (*types.Secret, error) {

	secret := envs.Get().GetState().Secrets().GetSecret(selflink)
	if secret != nil {
		return secret, nil
	}

	namespace, name := parseSecretSelflink(selflink)

	sr, err := envs.Get().GetRestClient().Namespace(namespace).Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("can not receive secret from api, err: %s", err.Error())
		return nil, err
	}

	return sr.Decode(), nil
}

func SecretCreate(ctx context.Context, selflink string) error {

	ok := envs.Get().GetState().Secrets().GetSecret(selflink)
	if ok != nil {
		return nil
	}

	namespace, name := parseSecretSelflink(selflink)
	secret, err := envs.Get().GetRestClient().Namespace(namespace).Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("get secret err: %s", err.Error())
		return err
	}

	envs.Get().GetState().Secrets().AddSecret(secret.Meta.SelfLink, secret.Decode())
	return nil
}

func SecretUpdate(ctx context.Context, selflink string) error {

	namespace, name := parseSecretSelflink(selflink)
	secret, err := envs.Get().GetRestClient().Namespace(namespace).Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("get secret err: %s", err.Error())
		return err
	}

	envs.Get().GetState().Secrets().AddSecret(secret.Meta.SelfLink, secret.Decode())
	return nil

}


func SecretRemove(ctx context.Context, selflink string) {
	envs.Get().GetState().Secrets().DelSecret(selflink)
}


func parseSecretSelflink(selflink string) (string, string) {
	var namespace, name string

	parts := strings.Split(selflink, ":")

	if len(parts) == 1 {
		namespace = types.DEFAULT_NAMESPACE
		name = parts[0]
	}

	if len(parts) >1 {
		namespace = parts[0]
		name = parts[1]
	}

	return namespace, name

}