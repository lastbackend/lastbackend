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

package runtime

import (
	"context"
	"strings"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

func (r Runtime) SecretGet(ctx context.Context, namespace, name string) (*types.Secret, error) {

	cli := r.retClient
	if cli == nil {
		return nil, nil
	}

	secret := r.state.Secrets().GetSecret(types.NewSecretSelfLink(namespace, name).String())
	if secret != nil {
		return secret, nil
	}

	sr, err := cli.V1().Namespace(namespace).Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("can not receive secret from api, err: %s", err.Error())
		return nil, err
	}

	return sr.Decode(), nil
}

func (r Runtime) SecretCreate(ctx context.Context, namespace, name string) error {

	cli := r.retClient
	if cli == nil {
		return nil
	}

	ok := r.state.Secrets().GetSecret(types.NewSecretSelfLink(namespace, name).String())
	if ok != nil {
		return nil
	}

	secret, err := cli.V1().Namespace(namespace).Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("get secret err: %s", err.Error())
		return err
	}

	r.state.Secrets().AddSecret(secret.Meta.SelfLink, secret.Decode())

	return nil
}

func (r Runtime) SecretUpdate(ctx context.Context, selflink string) error {

	cli := r.retClient
	if cli == nil {
		return nil
	}

	namespace, name := r.parseSecretSelflink(selflink)

	secret, err := cli.V1().Namespace(namespace).Secret(name).Get(ctx)
	if err != nil {
		log.Errorf("get secret err: %s", err.Error())
		return err
	}

	r.state.Secrets().AddSecret(secret.Meta.SelfLink, secret.Decode())

	return nil

}

func (r Runtime) SecretRemove(ctx context.Context, selflink string) {
	r.state.Secrets().DelSecret(selflink)
}

func (r Runtime) parseSecretSelflink(selflink string) (string, string) {
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
