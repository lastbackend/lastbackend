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
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)

func containerInspect(ctx context.Context, container *types.PodContainer) error {
	info, err := envs.Get().GetCRI().Inspect(ctx, container.ID)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Errorf("Stop inspect container err: %v", err)
			return nil
		}
		log.Errorf("Can-not inspect container: %v", err)
		return err
	}

	container.Pod = info.Pod
	container.Name = info.Name
	container.Exec = info.Exec
	container.Envs = info.Envs
	container.Binds = info.Binds

	container.Image = types.PodContainerImage{
		Name: info.Image,
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
			pod.Runtime.Services[c.ID] = container
			state.SetPod(c.Pod, pod)
		}
	}

	return nil
}

func containerManifestCreate(ctx context.Context, pod string, spec *types.SpecTemplateContainer) (*types.ContainerManifest, error) {

	mf := types.NewContainerManifest(spec)

	var namespace, name string

	parts := strings.Split(pod, ":")

	if len(parts) == 1 {
		namespace = types.SYSTEM_NAMESPACE
		name = parts[0]
	}

	if len(parts) >= 2 {
		namespace = parts[0]
		name = parts[1]
	}

	mf.Name = fmt.Sprintf("%s-%s", name, spec.Name)

	mf.Labels = make(map[string]string, 0)
	for n, v := range spec.Labels {
		mf.Labels[n] = v
	}

	mf.Labels[types.ContainerTypeLBC] = pod
	mf.Labels[types.ContainerTypeRuntime] = types.ContainerTypeRuntimeService

	for _, s := range spec.EnvVars {

		switch true {

		case s.Secret.Name != types.EmptyString && s.Secret.Key != types.EmptyString:

			secret, err := SecretGet(ctx, namespace, s.Secret.Name)
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

		case s.Config.Name != types.EmptyString && s.Config.Key != types.EmptyString:
			configSelfLink := fmt.Sprintf("%s:%s", namespace, s.Config.Name)
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
		default:
			continue
		}

	}

	for _, v := range spec.Volumes {

		if v.Name == types.EmptyString || v.MountPath == types.EmptyString {
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

		hostPath := vol.Status.Path

		if v.SubPath != types.EmptyString {
			hostPath = path.Join(hostPath, v.SubPath)
		}

		log.Debugf("attach volume: %s to %s:%s", v.Name, hostPath, v.MountPath)

		mf.Binds = append(mf.Binds, fmt.Sprintf("%s:%s:%s", hostPath, v.MountPath, v.Mode))
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
