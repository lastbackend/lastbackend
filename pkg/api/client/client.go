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

package client

import (
	"github.com/lastbackend/lastbackend/pkg/api/client/config"
	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	ClientHTTP = "http"
	ClientGRPC = "grpc"
)

type IClient interface {
	V1() types.ClientV1
}

func New(driver string, endpoint string, config *config.Config) (IClient, error) {
	switch driver {
	case ClientHTTP:
		return http.New(endpoint, config)
	default:
		log.Panicf("driver %s not defined", driver)
	}
	return nil, nil
}

func NewConfig() *config.Config {
	return new(config.Config)
}

func NewTLSConfig() *config.TLSConfig {
	return new(config.TLSConfig)
}
