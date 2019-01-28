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
	"strings"
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

	cs := make(chan *types.Container)

	go envs.Get().GetCRI().Subscribe(ctx, cs)

	for c := range cs {

		container := state.GetContainer(c.ID)
		if container == nil {
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

		pod := state.GetPod(c.Pod)
		if pod != nil {
			pod.Containers[c.ID] = container
			state.SetPod(c.Pod, pod)
		}
	}

	return nil
}

func containerManifestCreate(ctx context.Context, pod string, spec *types.SpecTemplateContainer) (*types.ContainerManifest, error) {

	mf := types.NewContainerManifest(spec)

	name := strings.Split(pod, ":")
	mf.Name = fmt.Sprintf("%s-%s", name[len(name)-1], spec.Name)

	mf.Labels = make(map[string]string, 0)
	for n, v := range spec.Labels {
		mf.Labels[n] = v
	}

	mf.Labels[types.ContainerTypeLBC] = pod

	for _, s := range spec.EnvVars {

		switch true {

		case s.Secret.Name != types.EmptyString && s.Secret.Key != types.EmptyString:

			secretSelfLink := fmt.Sprintf("%s:%s", name[0], s.Secret.Name)

			secret, err := SecretGet(ctx, secretSelfLink)
			if err != nil {
				log.Errorf("Can not get secret for container: %s", err.Error())
				return nil, err
			}

			if secret == nil {
				continue
			}

			if _, ok := secret.Spec.Data[s.Secret.Key]; !ok {
				continue
			}

			val, err := secret.DecodeSecretTextData(s.Secret.Key)
			if err != nil {
				continue
			}

			env := fmt.Sprintf("%s=%s", s.Name, val)
			mf.Envs = append(mf.Envs, env)
			break

		case s.Config.Name != types.EmptyString && s.Config.Key != types.EmptyString:
			configSelfLink := fmt.Sprintf("%s:%s", name[0], s.Config.Name)
			config := envs.Get().GetState().Configs().GetConfig(configSelfLink)
			if config == nil {
				log.Errorf("Can not get config for container: %s", configSelfLink)
				continue
			}

			value, ok := config.Data[s.Config.Key]
			if !ok {
				continue
			}

			env := fmt.Sprintf("%s=%s", s.Name, value)
			mf.Envs = append(mf.Envs, env)
			break
		default:
			continue
		}

	}

	for _, v := range spec.Volumes {

		if v.Name == types.EmptyString || v.Path == types.EmptyString {
			continue
		}

		claim := envs.Get().GetState().Volumes().GetClaim(podVolumeClaimNameCreate(pod, v.Name))
		if claim == nil {
			log.Debugf("volume claim %s not found in volumes state", podVolumeClaimNameCreate(pod, v.Name))
			continue
		}

		vol := envs.Get().GetState().Volumes().GetVolume(claim.Volume)
		if vol == nil {
			log.Debugf("volume %s not found in volumes state", claim.Volume)
			continue
		}

		if v.Mode != "ro" {
			v.Mode = "rw"
		}

		log.Debugf("attach volume: %s to %s:%s", v.Name, vol.Status.Path, v.Path)

		mf.Binds = append(mf.Binds, fmt.Sprintf("%s:%s:%s", vol.Status.Path, v.Path, v.Mode))
	}

	// TODO: Add dns search option only for LB domains

	net := envs.Get().GetNet()
	if net != nil {
		if net.GetResolverIP() != types.EmptyString {
			mf.DNS.Server = append(mf.DNS.Server, net.GetResolverIP())
		}

		if len(net.GetExternalDNS()) != 0 {
			mf.DNS.Server = append(mf.DNS.Server, net.GetExternalDNS()...)
		}
	}

	return mf, nil
}
