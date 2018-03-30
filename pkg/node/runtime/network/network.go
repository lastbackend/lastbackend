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

package network

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)

func Restore(ctx context.Context) error {

	sn, err := envs.Get().GetCNI().Subnets(ctx)
	if err != nil {
		log.Errorf("Can-not get subnets from CNI", err.Error())
	}

	for _, s := range sn {
		envs.Get().GetState().Networks().SetSubnet(s)
	}

	return nil
}

func Create(ctx context.Context, sn *types.NetworkSpec) error {
	if err := envs.Get().GetCNI().Create(ctx, sn); err != nil {
		log.Errorf("Can not create network subnet: %s", err.Error())
		return err
	}

	envs.Get().GetState().Networks().AddSubnet(sn)
	return nil
}

func Destroy(ctx context.Context, sn *types.NetworkSpec) error {
	if err := envs.Get().GetCNI().Destroy(ctx, sn); err != nil {
		log.Errorf("Can not destroy network subnet: %s", err.Error())
		return err
	}

	envs.Get().GetState().Networks().DelSubnet(sn)
	return nil
}
