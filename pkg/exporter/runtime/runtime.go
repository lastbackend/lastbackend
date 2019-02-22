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
	"github.com/lastbackend/lastbackend/pkg/exporter/envs"
	"github.com/lastbackend/lastbackend/pkg/exporter/logger"
)

type Runtime struct {
	logger *logger.Logger
}

func (r *Runtime) Logger(ctx context.Context) error {
	return r.logger.Listen()
}

func NewRuntime() (*Runtime, error) {

	var (
		err error
	)

	r := new(Runtime)
	r.logger, err = logger.NewLogger()
	if err != nil {
		return nil, err
	}

	envs.Get().SetLogger(r.logger)

	return r, nil
}
