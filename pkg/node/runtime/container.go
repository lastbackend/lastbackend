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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
	"time"
)

func containerInspect(ctx context.Context, status *types.PodStatus, container *types.PodContainer) error {
	info, err := envs.Get().GetCRI().Inspect(ctx, container.ID)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Errorf("Stop inspect container err: %v", err)
			return nil
		}
		log.Errorf("Can-not inspect container: %v", err)
		return err
	} else {
		container.Image = types.PodContainerImage{
			Name: info.Image,
		}
		if info.Status == types.StatusStopped {
			container.State.Stopped = types.PodContainerStateStopped{
				Stopped: true,
				Exit: types.PodContainerStateExit{
					Code:      info.ExitCode,
					Timestamp: time.Now().UTC(),
				},
			}
		}
	}

	if status.Network.PodIP == "" {
		status.Network.PodIP = info.Network.IPAddress
	}

	return nil
}

func containerSubscribe(ctx context.Context) error {

	state := envs.Get().GetState().Pods()

	cs, err := envs.Get().GetCRI().Subscribe(ctx)
	if err != nil {
		log.Errorf("container subscribe error: %s", err)
	}

	go func() {
		for c := range cs {

			if c.Pod != types.ContainerTypeLBC {
				continue
			}

			container := state.GetContainer(c.ID)
			if container == nil {
				log.V(logLevel).Debugf("Container not found")
				continue
			}

			container.Pod = c.Pod

			switch c.State {
			case types.StateDestroyed:
				state.DelContainer(container)
				break
			case types.StateCreated:
				container.State = types.PodContainerState{
					Created: types.PodContainerStateCreated{
						Created: time.Now().UTC(),
					},
				}
			case types.StateStarted:
				if container.State.Started.Started {
					continue
				}
				container.State = types.PodContainerState{
					Started: types.PodContainerStateStarted{
						Started:   true,
						Timestamp: time.Now().UTC(),
					},
				}
				container.State.Stopped.Stopped = false
			case types.StatusStopped:
				if container.State.Stopped.Stopped {
					continue
				}
				container.State.Stopped.Stopped = true
				container.State.Stopped.Exit = types.PodContainerStateExit{
					Code:      c.ExitCode,
					Timestamp: time.Now().UTC(),
				}
				container.State.Started.Started = false
			case types.StateError:
				if container.State.Error.Error {
					continue
				}
				container.State.Error.Error = true
				container.State.Error.Message = c.Status
				container.State.Error.Exit = types.PodContainerStateExit{
					Code:      c.ExitCode,
					Timestamp: time.Now().UTC(),
				}
				container.State.Started.Started = false
				container.State.Stopped.Stopped = false
				container.State.Stopped.Exit = types.PodContainerStateExit{
					Code:      c.ExitCode,
					Timestamp: time.Now().UTC(),
				}
				container.State.Started.Started = false
			}

			if err != nil {
				log.Errorf("Container: can-not inspect")
				break
			}

			if c.State != types.StateDestroyed {
				state.SetContainer(container)
			}

			events.NewPodStatusEvent(ctx, c.Pod)
		}
	}()

	return nil
}

func containerManifestCreate(ctx context.Context, pod string, spec *types.SpecTemplateContainer) (*types.ContainerManifest, error) {

	mf := types.NewContainerManifest(spec)

	for _, s := range spec.EnvVars {

		if s.From.Name == types.EmptyString || s.From.Key == types.EmptyString {
			continue
		}

		secret, err := SecretGet(ctx, s.From.Name)
		if err != nil {
			log.Errorf("Can not get secret for container: %s", err.Error())
			return nil, err
		}

		if secret == nil {
			continue
		}

		if _, ok := secret.Data[s.From.Key]; !ok {
			continue
		}

		val, err := secret.DecodeSecretTextData(s.From.Key)
		if err != nil {
			continue
		}

		env := fmt.Sprintf("%s=%s", s.Name, val)
		mf.Envs = append(mf.Envs, env)

	}

	for _, v := range spec.Volumes {

		log.Debugf("try to attach volume: %s", v.Name)

		if v.Name == types.EmptyString || v.Path == types.EmptyString {
			continue
		}

		vol := envs.Get().GetState().Volumes().GetVolume(podVolumeKeyCreate(pod, v.Name))
		if vol == nil {
			log.Debugf("volume %s not found in volumes state", v.Name)
			continue
		}

		if v.Mode != "rw" {
			v.Mode = "ro"
		}

		log.Debugf("attach volume: %s to %s:%s", v.Name, vol.Path, v.Path)

		mf.Binds = append(mf.Binds, fmt.Sprintf("%s:%s:%s", vol.Path, v.Path, v.Mode))
	}

	return mf, nil
}
