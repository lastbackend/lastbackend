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
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

func (r Runtime) containerInspect(ctx context.Context, container *models.PodContainer) error {
	info, err := r.cri.Inspect(ctx, container.ID)
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

	container.Image = models.PodContainerImage{
		Name: info.Image,
	}

	return nil
}

func (r Runtime) containerSubscribe(ctx context.Context) error {

	state := r.state.Pods()

	cs := make(chan *models.Container)

	go r.cri.Subscribe(ctx, cs)

	for c := range cs {

		container := state.GetContainer(c.ID)
		if container == nil {
			continue
		}

		container.Pod = c.Pod

		switch c.State {
		case models.StateDestroyed:
			state.DelContainer(container)
			break
		case models.StateCreated:
			container.State = models.PodContainerState{
				Created: models.PodContainerStateCreated{
					Created: time.Now().UTC(),
				},
			}
		case models.StateStarted:
			if container.State.Started.Started {
				continue
			}
			container.State = models.PodContainerState{
				Started: models.PodContainerStateStarted{
					Started:   true,
					Timestamp: time.Now().UTC(),
				},
			}
			container.State.Stopped.Stopped = false
		case models.StatusStopped:
			if container.State.Stopped.Stopped {
				continue
			}
			container.State.Stopped.Stopped = true
			container.State.Stopped.Exit = models.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			container.State.Started.Started = false
		case models.StateError:
			if container.State.Error.Error {
				continue
			}
			container.State.Error.Error = true
			container.State.Error.Message = c.Status
			container.State.Error.Exit = models.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			container.State.Started.Started = false
			container.State.Stopped.Stopped = false
			container.State.Stopped.Exit = models.PodContainerStateExit{
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

func (r Runtime) containerManifestCreate(ctx context.Context, pod string, spec *models.SpecTemplateContainer) (*models.ContainerManifest, error) {

	mf := models.NewContainerManifest(spec)

	var (
		namespace, name string
	)

	parts := strings.Split(pod, ":")

	if len(parts) == 1 {
		namespace = models.SYSTEM_NAMESPACE
		name = parts[0]
	}

	if len(parts) >= 2 {
		namespace = parts[0]
		name = parts[len(parts)-1]
	}

	mf.Name = fmt.Sprintf("%s-%s", name, spec.Name)
	mf.Labels = make(map[string]string, 0)
	for n, v := range spec.Labels {
		mf.Labels[n] = v
	}

	mf.Labels[models.ContainerTypeLBC] = pod
	mf.Labels[models.ContainerTypeRuntime] = models.ContainerTypeRuntimeService

	for _, s := range spec.EnvVars {

		switch true {

		case s.Secret.Name != models.EmptyString && s.Secret.Key != models.EmptyString:

			secret, err := r.SecretGet(ctx, namespace, s.Secret.Name)
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

		case s.Config.Name != models.EmptyString && s.Config.Key != models.EmptyString:
			configSelfLink := fmt.Sprintf("%s:%s", namespace, s.Config.Name)
			config := r.state.Configs().GetConfig(configSelfLink)
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

		if v.Name == models.EmptyString || v.MountPath == models.EmptyString {
			continue
		}

		claim := r.state.Volumes().GetClaim(r.podVolumeClaimNameCreate(pod, v.Name))
		if claim == nil {
			log.Debugf("volume claim %s not found in volumes state", r.podVolumeClaimNameCreate(pod, v.Name))
			continue
		}

		vol := r.state.Volumes().GetVolume(claim.Volume)
		if vol == nil {
			log.Debugf("volume %s not found in volumes state", claim.Volume)
			continue
		}

		if v.Mode != "ro" {
			v.Mode = "rw"
		}

		hostPath := vol.Status.Path

		if v.SubPath != models.EmptyString {
			hostPath = path.Join(hostPath, v.SubPath)
		}

		log.Debugf("attach volume: %s to %s:%s", v.Name, hostPath, v.MountPath)

		mf.Binds = append(mf.Binds, fmt.Sprintf("%s:%s:%s", hostPath, v.MountPath, v.Mode))
	}

	// TODO: Add dns search option only for LB domains

	net := r.network
	if net != nil {
		if net.GetResolverIP() != models.EmptyString {
			mf.DNS.Server = append(mf.DNS.Server, net.GetResolverIP())
		}

		if len(net.GetExternalDNS()) != 0 {
			mf.DNS.Server = append(mf.DNS.Server, net.GetExternalDNS()...)
		}
	}

	return mf, nil
}
