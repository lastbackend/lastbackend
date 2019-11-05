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
	"github.com/lastbackend/lastbackend/internal/agent/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
	"time"
)

const (
	logServicePrefix = "node:runtime:service"
)

func serviceStart(ctx context.Context, pod string, m *types.ContainerManifest, status *types.PodStatus) error {

	var (
		err error
		c   = new(types.PodContainer)
	)

	c.ID, err = envs.Get().GetCRI().Create(ctx, m)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Errorf("%s stop creating container: %s", logServicePrefix, err.Error())
			return nil
		}

		log.Errorf("%s can-not create container: %s", logServicePrefix, err)

		c.State.Error = types.PodContainerStateError{
			Error:   true,
			Message: err.Error(),
			Exit: types.PodContainerStateExit{
				Timestamp: time.Now().UTC(),
			},
		}
		return err
	}

	status.Runtime.Services[c.ID] = c

	if err := containerInspect(context.Background(), c); err != nil {
		log.Errorf("%s inspect container after create: err %s", logServicePrefix, err.Error())
		PodClean(context.Background(), status)
		return err
	}

	//==========================================================================
	// Start container =========================================================
	//==========================================================================

	c.State.Created = types.PodContainerStateCreated{
		Created: time.Now().UTC(),
	}

	envs.Get().GetState().Pods().SetPod(pod, status)
	log.V(logLevel).Debugf("%s container created: %s", logServicePrefix, c.ID)

	if err := envs.Get().GetCRI().Start(ctx, c.ID); err != nil {

		log.Errorf("%s can-not start container: %s", logServicePrefix, err)

		switch err {
		case context.Canceled:
			log.Errorf("%s stop starting container err: %s", logServicePrefix, err.Error())
			return nil
		}

		c.State.Error = types.PodContainerStateError{
			Error:   true,
			Message: err.Error(),
			Exit: types.PodContainerStateExit{
				Timestamp: time.Now().UTC(),
			},
		}

		return err
	}

	log.V(logLevel).Debugf("%s container started: %s", logServicePrefix, c.ID)

	if err := containerInspect(context.Background(), c); err != nil {
		log.Errorf("%s inspect container after create: err %s", logServicePrefix, err.Error())
		return err
	}

	c.Ready = true
	c.State.Started = types.PodContainerStateStarted{
		Started:   true,
		Timestamp: time.Now().UTC(),
	}

	info, err := envs.Get().GetCRI().Inspect(ctx, c.ID)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Errorf("Stop inspect container err: %v", err)
			return nil
		}
		log.Errorf("Can-not inspect container: %v", err)
		return err
	}

	if status.Network.PodIP == "" {
		status.Network.PodIP = info.Network.IPAddress
	}

	envs.Get().GetState().Pods().SetPod(pod, status)
	return nil
}

func serviceStop() error {

	return nil
}

func serviceRestart() error {

	return nil
}

func serviceRemove() error {
	return nil
}
