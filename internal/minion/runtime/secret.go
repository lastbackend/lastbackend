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
	"github.com/lastbackend/lastbackend/internal/minion/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
	"strings"
)

func SecretGet(ctx context.Context, namespace, name string) (*types.Secret, error) {

	cli := envs.Get().GetRestClient()
	if cli == nil {
		return nil, nil
	}

	secret := envs.Get().GetState().Secrets().GetSecret(types.NewSecretSelfLink(namespace, name).String())
	if secret != nil {
		return secret, nil
	}

	sr, err := cli.Namespace(namespace).Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("can not receive secret from api, err: %s", err.Error())
		return nil, err
	}

	return sr.Decode(), nil
}

func SecretCreate(ctx context.Context, namespace, name string) error {

	cli := envs.Get().GetRestClient()
	if cli == nil {
		return nil
	}

	ok := envs.Get().GetState().Secrets().GetSecret(types.NewSecretSelfLink(namespace, name).String())
	if ok != nil {
		return nil
	}

	secret, err := cli.Namespace(namespace).Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("get secret err: %s", err.Error())
		return err
	}

	envs.Get().GetState().Secrets().AddSecret(secret.Meta.SelfLink, secret.Decode())

	return nil
}

func SecretUpdate(ctx context.Context, selflink string) error {

	cli := envs.Get().GetRestClient()
	if cli == nil {
		return nil
	}

	namespace, name := parseSecretSelflink(selflink)

	secret, err := cli.Namespace(namespace).Secret(name).Get(ctx)
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

	parts := strings.SplitN(selflink, ":", 1)

	if len(parts) == 1 {
		namespace = types.DefaultNamespace
		name = parts[0]
	}

	if len(parts) > 1 {
		namespace = parts[0]
		name = parts[1]
	}

	return namespace, name

}
