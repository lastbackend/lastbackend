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

package namespace

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logPrefix = "api:handler:namespace:fetchFromRequest"
	logLevel  = 3
)

func FetchFromRequest(ctx context.Context, selflink string) (*types.Namespace, *errors.Err) {

	nm := model.NewNamespaceModel(ctx, envs.Get().GetStorage())
	ns, err := nm.Get(selflink)

	if err != nil {
		log.V(logLevel).Errorf("%s:> get namespace err: %s", logPrefix, err.Error())
		return nil, errors.New("namespace").InternalServerError(err)
	}

	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:> get namespace err: %s", logPrefix, err.Error())
		return nil, errors.New("namespace").NotFound()
	}

	return ns, nil
}
