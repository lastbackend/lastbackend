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

func NetworkRestore(ctx context.Context) error {

	sn, err := envs.Get().GetCNI().Subnets(ctx)
	if err != nil {
		log.Errorf("Can-not get subnets from CNI err: %v", err)
	}

	for cidr, s := range sn {
		envs.Get().GetState().Networks().SetSubnet(cidr, s)
	}

	return nil
}

func NetworkManage(ctx context.Context, cidr string, sn *types.SubnetManifest) error {

	subnets := envs.Get().GetState().Networks().GetSubnets()
	if state, ok := subnets[cidr]; ok {

		if sn.State == types.StateDestroy {
			envs.Get().GetCNI().Destroy(ctx, &state)
			envs.Get().GetState().Networks().DelSubnet(cidr)
			return nil
		}

		// TODO: check if network manifest changes
		// if changes then update routes and interfaces
		return nil
	}

	if sn.State == types.StateDestroy {
		return nil
	}

	state, err := envs.Get().GetCNI().Create(ctx, sn)
	if err != nil {
		log.Errorf("Can not create network subnet: %s", err.Error())
		return err
	}

	envs.Get().GetState().Networks().AddSubnet(cidr, state)
	return nil
}

func NetworkDestroy(ctx context.Context, cidr string) error {

	sn := envs.Get().GetState().Networks().GetSubnet(cidr)

	if err := envs.Get().GetCNI().Destroy(ctx, sn); err != nil {
		log.Errorf("Can not destroy network subnet: %s", err.Error())
		return err
	}

	envs.Get().GetState().Networks().DelSubnet(cidr)
	return nil
}
