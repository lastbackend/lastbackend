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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"io"
	"strings"
	"time"

	"net/http"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/util/cleaner"
)

const (
	BUFFER_SIZE = 1024
)

func PodManage(ctx context.Context, key string, manifest *types.PodManifest) error {
	log.V(logLevel).Debugf("Provision pod: %s", key)

	//==========================================================================
	// Destroy pod =============================================================
	//==========================================================================

	// Call destroy pod
	if manifest.State.Destroy {

		if task := envs.Get().GetState().Tasks().GetTask(key); task != nil {
			log.V(logLevel).Debugf("Cancel pod creating: %s", key)
			task.Cancel()
		}

		p := envs.Get().GetState().Pods().GetPod(key)
		if p == nil {

			ps := types.NewPodStatus()
			ps.SetDestroyed()
			envs.Get().GetState().Pods().AddPod(key, ps)

			return nil
		}

		log.V(logLevel).Debugf("Pod found > destroy it: %s", key)

		PodDestroy(ctx, key, p)

		p.SetDestroyed()
		envs.Get().GetState().Pods().SetPod(key, p)
		return nil
	}

	//==========================================================================
	// Check containers pod status =============================================
	//==========================================================================

	// Get pod list from current state
	p := envs.Get().GetState().Pods().GetPod(key)
	if p != nil {

		switch true {
		case !PodSpecCheck(ctx, key, manifest):
			PodDestroy(ctx, key, p)
			break
		case !PodVolumesCheck(ctx, key, manifest.Template.Volumes):
			log.Debugf("Volumes data changed: %s", key)
			for _, v := range manifest.Template.Volumes {

				if v.Volume.Name != types.EmptyString {

					pv, err := PodVolumeAttach(ctx, key, v)
					if err != nil {
						log.Errorf("can not attach volume for pod: %s", err.Error())
						return err
					}

					p.Volumes[v.Name] = pv

				} else {

					var name string
					if v.Volume.Name != types.EmptyString {

						name = fmt.Sprintf("%s:%s", getPodNamespace(key), v.Name)
					} else {
						name = podVolumeKeyCreate(key, v.Name)
					}

					vol := envs.Get().GetState().Volumes().GetVolume(name)

					if vol == nil {
						log.V(logLevel).Debugf("Update pod volume: volume not found: create %s: %s", key, v.Name)

						vs, err := PodVolumeCreate(ctx, key, v)
						if err != nil {
							log.Errorf("can not update volume data: %s", err.Error())
							return err
						}

						pv := &types.VolumeClaim{
							Name:   podVolumeClaimNameCreate(key, v.Name),
							Volume: name,
							Path:   vs.Status.Path,
						}

						envs.Get().GetState().Volumes().SetClaim(pv.Name, pv)
						p.Volumes[pv.Name] = pv

					} else {

						_, err := PodVolumeUpdate(ctx, key, v)
						if err != nil {
							log.Errorf("can not update volume data: %s", err.Error())
							return err
						}
					}

				}

			}
			return PodRestart(ctx, key)
		default:
			return nil
		}
	}

	log.V(logLevel).Debugf("Pod not found > create it: %s", key)

	ctx, cancel := context.WithCancel(context.Background())
	envs.Get().GetState().Tasks().AddTask(key, &types.NodeTask{Cancel: cancel})

	status, err := PodCreate(ctx, key, manifest)
	if err != nil {
		log.Errorf("Can not create pod: %s err: %s", key, err.Error())
		status.SetError(err)
	}

	envs.Get().GetState().Pods().SetPod(key, status)
	return nil
}

func PodRestart(ctx context.Context, key string) error {

	pod := envs.Get().GetState().Pods().GetPod(key)
	if pod == nil {
		return errors.New("pod not found")
	}

	cri := envs.Get().GetCRI()

	for _, c := range pod.Containers {
		if err := cri.Restart(ctx, c.ID, nil); err != nil {
			return err
		}
	}

	return nil
}

func PodCreate(ctx context.Context, key string, manifest *types.PodManifest) (*types.PodStatus, error) {

	var (
		status    = types.NewPodStatus()
		namespace = getPodNamespace(key)
		setError  = func(err error) (*types.PodStatus, error) {
			log.Errorf("Can-not pull image: %s", err)
			status.SetError(err)
			PodClean(ctx, status)
			return status, err
		}
	)

	log.V(logLevel).Debugf("Create pod: %s", key)

	//==========================================================================
	// Pull image ==============================================================
	//==========================================================================

	status.SetPull()

	envs.Get().GetState().Pods().AddPod(key, status)

	log.V(logLevel).Debugf("Have %d volumes", len(manifest.Template.Volumes))
	for _, v := range manifest.Template.Volumes {

		var name string
		if v.Volume.Name != types.EmptyString {
			name = fmt.Sprintf("%s:%s", getPodNamespace(key), v.Name)
		} else {
			name = podVolumeKeyCreate(key, v.Name)
		}

		vol := envs.Get().GetState().Volumes().GetVolume(name)
		if vol == nil {
			log.V(logLevel).Debugf("Update pod volume: volume not found: create %s: %s", key, v.Name)

			vs, err := PodVolumeCreate(ctx, key, v)
			if err != nil {
				log.Errorf("can not update volume data: %s", err.Error())
				return status, err
			}

			pv := &types.VolumeClaim{
				Name:   podVolumeClaimNameCreate(key, v.Name),
				Volume: name,
				Path:   vs.Status.Path,
			}

			envs.Get().GetState().Volumes().SetClaim(pv.Name, pv)
			status.Volumes[pv.Name] = pv

		} else {
			_, err := PodVolumeUpdate(ctx, key, v)
			if err != nil {
				log.Errorf("can not update volume data: %s", err.Error())
				return status, err
			}
		}

		envs.Get().GetState().Pods().SetPod(key, status)
	}

	log.V(logLevel).Debugf("Have %d containers", len(manifest.Template.Containers))
	for _, c := range manifest.Template.Containers {
		log.V(logLevel).Debug("Pull images for pod if needed")
		if err := ImagePull(ctx, namespace, &c.Image); err != nil {
			log.Errorf("can not pull image: %s", err.Error())
			return setError(err)
		}
	}

	//==========================================================================
	// Run container ===========================================================
	//==========================================================================

	status.SetStarting()
	status.Steps[types.StepPull] = types.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}

	envs.Get().GetState().Pods().SetPod(key, status)

	for _, s := range manifest.Template.Containers {
		for _, e := range s.EnvVars {
			if e.Secret.Name != types.EmptyString {
				log.V(logLevel).Debug("Get secret info from api")
				if err := SecretCreate(ctx, fmt.Sprintf("%s:%s", namespace, e.Secret.Name)); err != nil {
					log.Errorf("can not fetch secret from api")
				}
			}
		}
	}

	for _, s := range manifest.Template.Containers {

		//==========================================================================
		// Create container ========================================================
		//==========================================================================

		var (
			c = new(types.PodContainer)
		)

		m, err := containerManifestCreate(ctx, key, s)
		if err != nil {
			log.Errorf("can not create container manifest from spec: %s", err.Error())
			return setError(err)
		}

		c.ID, err = envs.Get().GetCRI().Create(ctx, m)
		if err != nil {
			switch err {
			case context.Canceled:
				log.Errorf("Stop creating container: %s", err.Error())
				PodClean(context.Background(), status)
				return status, nil
			}

			log.Errorf("Can-not create container: %s", err)
			c.State.Error = types.PodContainerStateError{
				Error:   true,
				Message: err.Error(),
				Exit: types.PodContainerStateExit{
					Timestamp: time.Now().UTC(),
				},
			}
			return status, err
		}

		if err := containerInspect(context.Background(), status, c); err != nil {
			log.Errorf("Inspect container after create: err %s", err.Error())
			PodClean(context.Background(), status)
			return status, err
		}

		//==========================================================================
		// Start container =========================================================
		//==========================================================================

		c.State.Created = types.PodContainerStateCreated{
			Created: time.Now().UTC(),
		}
		status.Containers[c.ID] = c
		envs.Get().GetState().Pods().SetPod(key, status)
		log.V(logLevel).Debugf("Container created: %#v", c)

		if err := envs.Get().GetCRI().Start(ctx, c.ID); err != nil {

			log.Errorf("Can-not start container: %s", err)

			switch err {
			case context.Canceled:
				log.Errorf("Stop starting container err: %s", err.Error())
				PodClean(context.Background(), status)
				return status, nil
			}

			c.State.Error = types.PodContainerStateError{
				Error:   true,
				Message: err.Error(),
				Exit: types.PodContainerStateExit{
					Timestamp: time.Now().UTC(),
				},
			}

			status.Containers[c.ID] = c
			return status, err
		}

		log.V(logLevel).Debugf("Container started: %#v", c)

		if err := containerInspect(context.Background(), status, c); err != nil {
			log.Errorf("Inspect container after create: err %s", err.Error())
			return status, err
		}

		c.Ready = true
		c.State.Started = types.PodContainerStateStarted{
			Started:   true,
			Timestamp: time.Now().UTC(),
		}
		status.Containers[c.ID] = c
		envs.Get().GetState().Pods().SetPod(key, status)
	}

	status.SetRunning()
	status.Steps[types.StepReady] = types.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}

	envs.Get().GetState().Pods().SetPod(key, status)
	return status, nil
}

func PodClean(ctx context.Context, status *types.PodStatus) {

	for _, c := range status.Containers {
		log.V(logLevel).Debugf("Remove unnecessary container: %s", c.ID)
		if err := envs.Get().GetCRI().Remove(ctx, c.ID, true, true); err != nil {
			log.Warnf("Can-not remove unnecessary container %s: %s", c.ID, err)
		}
	}

	for _, c := range status.Containers {
		log.V(logLevel).Debugf("Try to clean image: %s", c.Image.Name)
		if err := ImageRemove(ctx, c.Image.Name); err != nil {
			log.Errorf("can not remove image: %s", err.Error())
			continue
		}
	}
}

func PodDestroy(ctx context.Context, pod string, status *types.PodStatus) {
	log.V(logLevel).Debugf("Try to remove pod: %s", pod)
	PodClean(ctx, status)
	envs.Get().GetState().Pods().DelPod(pod)
	for _, v := range status.Volumes {
		PodVolumeDestroy(ctx, pod, v.Name)
	}
}

func PodRestore(ctx context.Context) error {

	log.V(logLevel).Debug("Runtime restore state")

	cl, err := envs.Get().GetCRI().List(ctx, true)
	if err != nil {
		log.Errorf("Pods restore error: %s", err)
		return err
	}

	for _, c := range cl {

		if c.Pod == types.EmptyString {
			continue
		}

		log.V(logLevel).Debugf("Pod [%s] > container restore %s", c.Pod, c.ID)

		status := envs.Get().GetState().Pods().GetPod(c.Pod)
		if status == nil {
			status = types.NewPodStatus()
		}

		key := c.Pod

		cs := &types.PodContainer{
			ID:   c.ID,
			Name: c.Name,
			Image: types.PodContainerImage{
				Name: c.Image,
			},
			Envs:  c.Envs,
			Ports: c.Network.Ports,
			Binds: c.Binds,
		}

		cs.Restart.Policy = c.Restart.Policy
		cs.Restart.Attempt = c.Restart.Retry
		cs.Exec = c.Exec

		switch c.State {
		case types.StateCreated:
			cs.State = types.PodContainerState{
				Created: types.PodContainerStateCreated{
					Created: time.Now().UTC(),
				},
			}
		case types.StateStarted:
			cs.State = types.PodContainerState{
				Started: types.PodContainerStateStarted{
					Started:   true,
					Timestamp: time.Now().UTC(),
				},
			}
			cs.State.Stopped.Stopped = false
		case types.StatusStopped:
			cs.State.Stopped.Stopped = true
			cs.State.Stopped.Exit = types.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
		case types.StateError:

			cs.State.Error.Error = true
			cs.State.Error.Message = c.Status
			cs.State.Error.Exit = types.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
			cs.State.Stopped.Stopped = false
			cs.State.Stopped.Exit = types.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
		}

		if c.Status == types.StatusStopped {
			cs.State.Stopped = types.PodContainerStateStopped{
				Stopped: true,
				Exit: types.PodContainerStateExit{
					Code:      c.ExitCode,
					Timestamp: time.Now().UTC(),
				},
			}
		}

		cs.Ready = true
		status.Containers[cs.ID] = cs
		status.Network.PodIP = c.Network.IPAddress

		log.V(logLevel).Debugf("Container restored %s", c.ID)
		envs.Get().GetState().Pods().SetPod(key, status)
		log.V(logLevel).Debugf("Pod restored %s: %#v", key, status)
	}

	return nil
}

func PodLogs(ctx context.Context, id string, follow bool, s io.Writer, doneChan chan bool) error {

	log.V(logLevel).Debugf("Get container [%s] logs streaming", id)

	var (
		cri    = envs.Get().GetCRI()
		buffer = make([]byte, BUFFER_SIZE)
		done   = make(chan bool, 1)
	)

	req, err := cri.Logs(ctx, id, true, true, follow)
	if err != nil {
		log.Errorf("Error get logs stream %s", err)
		return err
	}
	defer func() {
		log.V(logLevel).Debugf("Stop container [%s] logs streaming", id)
		ctx.Done()
		close(done)
		req.Close()
	}()

	go func() {
		for {
			select {
			case <-done:
				req.Close()
				return
			default:

				n, err := cleaner.NewReader(req).Read(buffer)
				if err != nil {

					if err == context.Canceled {
						log.V(logLevel).Debug("Stream is canceled")
						return
					}

					log.Errorf("Error read bytes from stream %s", err)
					doneChan <- true
					return
				}

				_, err = func(p []byte) (n int, err error) {
					n, err = s.Write(p)
					if err != nil {
						log.Errorf("Error write bytes to stream %s", err)
						return n, err
					}

					if f, ok := s.(http.Flusher); ok {
						f.Flush()
					}
					return n, nil
				}(buffer[0:n])

				if err != nil {
					log.Errorf("Error written to stream %s", err)
					done <- true
					return
				}

				for i := 0; i < n; i++ {
					buffer[i] = 0
				}
			}
		}
	}()

	<-doneChan

	return nil
}

func PodSpecCheck(ctx context.Context, key string, manifest *types.PodManifest) bool {

	log.V(logLevel).Infof("Pod check spec pod: %s", key)

	state := envs.Get().GetState().Pods().GetPod(key)

	var statec = make(map[string]*types.ContainerManifest, 0)
	var specc = make(map[string]*types.ContainerManifest, 0)

	for _, c := range manifest.Template.Containers {
		mf, err := containerManifestCreate(ctx, key, c)
		if err != nil {
			return false
		}
		specc[mf.Name] = mf
	}

	for _, c := range state.Containers {
		statec[c.Name] = c.GetManifest()
	}

	if len(statec) != len(specc) {
		log.Debugf("container spec count not equal not exists: %d != %d", len(statec), len(specc))
		return false
	}

	for n, mf := range specc {

		if _, ok := statec[n]; !ok {
			log.Debugf("container spec not exists: %s", n)
			return false
		}

		// check image

		c := statec[n]

		if c.Image != mf.Image {
			log.Debugf("images not equal: %s != %s", c.Image, mf.Image)
			return false
		}

		img := envs.Get().GetState().Images().GetImage(c.Image)
		if img == nil {
			log.Debugf("image not found in state: %s", mf.Image)
			return false
		}

		if len(mf.Exec.Command) == 0 {
			if strings.Join(c.Exec.Command, " ") != strings.Join(img.Status.Container.Exec.Command, " ") {
				log.Debugf("cmd different with img cmd: %s != %s",
					strings.Join(c.Exec.Command, " "),
					strings.Join(img.Status.Container.Exec.Command, " "))
				return false
			}
		} else {
			if strings.Join(c.Exec.Command, " ") != strings.Join(mf.Exec.Command, " ") {
				log.Debugf("cmd different with manifest cmd: %s != %s",
					strings.Join(c.Exec.Command, " "),
					strings.Join(mf.Exec.Command, " "))
				return false
			}
		}

		if len(mf.Exec.Entrypoint) == 0 {
			if strings.Join(c.Exec.Entrypoint, " ") != strings.Join(img.Status.Container.Exec.Entrypoint, " ") {
				log.Debugf("entrypoint changed: %s != %s",
					strings.Join(c.Exec.Entrypoint, " "),
					strings.Join(img.Status.Container.Exec.Entrypoint, " "))
				return false
			}
		} else {
			if strings.Join(c.Exec.Entrypoint, " ") != strings.Join(mf.Exec.Entrypoint, " ") {
				log.Debugf("entrypoint changed: %s != %s",
					strings.Join(c.Exec.Entrypoint, " "),
					strings.Join(mf.Exec.Entrypoint, " "))
				return false
			}
		}

		if mf.Exec.Workdir == types.EmptyString {
			if c.Exec.Workdir != img.Status.Container.Exec.Workdir {
				log.Debugf("workdir changed: %s != %s", c.Exec.Workdir, img.Status.Container.Exec.Workdir)
				return false
			}
		} else {
			if c.Exec.Workdir != mf.Exec.Workdir {
				log.Debugf("workdir changed: %s != %s", c.Exec.Workdir, mf.Exec.Workdir)
				return false
			}
		}

		if len(mf.Exec.Args) != 0 {
			if strings.Join(c.Exec.Args, " ") != strings.Join(mf.Exec.Args, " ") {
				log.Debugf("args changed: %s != %s",
					strings.Join(c.Exec.Args, " "),
					strings.Join(mf.Exec.Args, " "))
				return false
			}
		}

		// Check environments
		for _, e := range mf.Envs {
			var f = false
			for _, ie := range c.Envs {

				if ie == e {
					f = true
					break
				}
			}
			if !f {
				log.Debugf("Env not found:%s", e)
				return false
			}
		}

		for _, e := range c.Envs {
			var f = false
			for _, ie := range mf.Envs {
				if ie == e {
					f = true
					break
				}
			}

			if !f {
				for _, ie := range img.Status.Container.Envs {
					if ie == e {
						f = true
						break
					}
				}
			}

			if !f {
				log.Debugf("\tEnv is unnecessary:%s", e)
				return false
			}
		}

		// Check binds
		for _, e := range mf.Binds {
			var f = false
			for _, ie := range c.Binds {
				if ie == e {
					f = true
					break
				}
			}
			if !f {
				log.Debugf("Bind not found:%s", e)
				return false
			}
		}

		for _, e := range c.Binds {
			var f = false
			for _, ie := range mf.Binds {
				if ie == e {
					f = true
					break
				}
			}
			if !f {
				log.Debugf("Bind is unnecessary:%s", e)
				return false
			}
		}

		// Check ports
		for _, e := range mf.Ports {
			var f = false
			for _, ie := range c.Ports {
				if e.HostIP != types.EmptyString {
					if e.HostIP != ie.HostIP {
						log.Debugf("\t Port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s",
							e.HostIP, e.HostPort, e.ContainerPort, e.Protocol,
							ie.HostIP, ie.HostPort, ie.ContainerPort, ie.Protocol)
						return false
					}
				}

				if e.Protocol != types.EmptyString {
					if e.Protocol != ie.Protocol {
						log.Debugf("\t Port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s",
							e.HostIP, e.HostPort, e.ContainerPort, e.Protocol,
							ie.HostIP, ie.HostPort, ie.ContainerPort, ie.Protocol)
						return false
					}
				} else {
					if ie.Protocol != "tcp" {
						log.Debugf("\t Port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s",
							e.HostIP, e.HostPort, e.ContainerPort, e.Protocol,
							ie.HostIP, ie.HostPort, ie.ContainerPort, ie.Protocol)
						return false
					}
				}

				if ie.ContainerPort == e.ContainerPort &&
					ie.HostPort == ie.HostPort {
					f = true
					break
				}
			}

			if !f {
				log.Debugf("\t Port map not found: \t\t %s:%d:%d/%s ",
					e.HostIP, e.HostPort, e.ContainerPort, e.Protocol)
				return false
			}
		}

		for _, e := range c.Ports {
			var f = false
			for _, ie := range mf.Ports {
				if ie.ContainerPort == e.ContainerPort &&
					ie.HostPort == ie.HostPort &&
					ie.Protocol == ie.Protocol &&
					ie.HostIP == ie.HostIP {
					f = true
					break
				}
			}

			if !f {
				log.Debugf("Port map is unnecessary: %#v", e)
				return false
			}
		}

		if mf.RestartPolicy.Policy != c.RestartPolicy.Policy ||
			mf.RestartPolicy.Attempt != c.RestartPolicy.Attempt {

			log.Debugf("Restart policy changed: %s:%d => %s:%d",
				c.RestartPolicy.Policy, c.RestartPolicy.Attempt,
				mf.RestartPolicy.Policy, mf.RestartPolicy.Attempt)
			return false
		}

	}

	return true

}

func PodVolumesCheck(ctx context.Context, pod string, spec []*types.SpecTemplateVolume) bool {

	log.V(logLevel).Debugf("Check pod volumes: %s: %d", pod, len(spec))

	for _, v := range spec {
		name := podVolumeKeyCreate(pod, v.Name)

		if v.Config.Name != types.EmptyString && len(v.Config.Files) > 0 {
			equal, err := VolumeCheckConfigData(ctx, name, v.Config.Name)
			if err != nil {
				return false
			}
			return equal
		}

		if v.Secret.Name != types.EmptyString && len(v.Secret.Files) > 0 {
			equal, err := VolumeCheckSecretData(ctx, name, v.Config.Name)
			if err != nil {
				return false
			}
			return equal
		}
	}

	return true
}

func PodVolumeUpdate(ctx context.Context, pod string, spec *types.SpecTemplateVolume) (*types.VolumeStatus, error) {

	log.V(logLevel).Debugf("Update pod volume: %s: %s", pod, spec.Name)

	path := strings.Replace(pod, ":", "-", -1)
	path = fmt.Sprintf("%s-%s", path, spec.Name)

	var (
		name = podVolumeKeyCreate(pod, spec.Name)
	)

	status := envs.Get().GetState().Volumes().GetVolume(name)

	if spec.Secret.Name != types.EmptyString && len(spec.Secret.Files) > 0 {
		if err := VolumeSetSecretData(ctx, name, spec.Secret.Name); err != nil {
			log.Errorf("can not set config data to volume: %s", err.Error())
			return status, err
		}
	}

	if spec.Secret.Name == types.EmptyString && spec.Config.Name != types.EmptyString && len(spec.Config.Files) > 0 {
		if err := VolumeSetConfigData(ctx, name, spec.Config.Name); err != nil {
			log.Errorf("can not set config data to volume: %s", err.Error())
			return status, err
		}
	}

	return status, nil
}

func PodVolumeAttach(ctx context.Context, pod string, spec *types.SpecTemplateVolume) (*types.VolumeClaim, error) {

	log.V(logLevel).Debugf("Attach pod volume: %s: %s", pod, spec.Name)

	var name = fmt.Sprintf("%s:%s", getPodNamespace(pod), spec.Name)

	volume := envs.Get().GetState().Volumes().GetVolume(name)
	if volume == nil {
		return nil, errors.New("volume not found on node")
	}

	pv := &types.VolumeClaim{
		Name:   podVolumeClaimNameCreate(pod, spec.Name),
		Volume: name,
		Path:   volume.Status.Path,
	}

	envs.Get().GetState().Volumes().SetClaim(pv.Name, pv)

	return pv, nil
}

func PodVolumeCreate(ctx context.Context, pod string, spec *types.SpecTemplateVolume) (*types.VolumeStatus, error) {

	log.V(logLevel).Debugf("Create pod volume: %s:%s", pod, spec.Name)

	path := strings.Replace(pod, ":", "-", -1)
	path = fmt.Sprintf("%s-%s", path, spec.Name)

	var (
		name = podVolumeKeyCreate(pod, spec.Name)
		vm   = types.VolumeManifest{
			HostPath: path,
			Type:     types.KindVolumeHostDir,
		}
	)

	st, err := VolumeCreate(ctx, name, &vm)
	if err != nil {
		log.Errorf("can not create pod volume: %s", err.Error())
		return nil, err
	}

	if spec.Secret.Name != types.EmptyString && len(spec.Secret.Files) > 0 {
		if err := VolumeSetSecretData(ctx, name, spec.Secret.Name); err != nil {
			log.Errorf("can not set secret data to volume: %s", err.Error())
			return st, err
		}
	}

	if spec.Secret.Name == types.EmptyString && spec.Config.Name != types.EmptyString && len(spec.Config.Files) > 0 {
		if err := VolumeSetConfigData(ctx, name, spec.Config.Name); err != nil {
			log.Errorf("can not set config data to volume: %s", err.Error())
			return st, err
		}
	}

	envs.Get().GetState().Volumes().SetLocal(name)

	return st, nil
}

func PodVolumeDestroy(ctx context.Context, pod, volume string) error {
	envs.Get().GetState().Volumes().DelLocal(podVolumeKeyCreate(pod, volume))
	return VolumeDestroy(ctx, podVolumeKeyCreate(pod, volume))
}

func podVolumeKeyCreate(pod, volume string) string {
	return fmt.Sprintf("%s-%s", strings.Replace(pod, ":", "-", -1), volume)
}

func podVolumeClaimNameCreate(pod, volume string) string {
	return fmt.Sprintf("%s:%s", pod, volume)
}

// TODO: move to distribution
func getPodNamespace(key string) string {
	var namespace = types.DEFAULT_NAMESPACE

	parts := strings.Split(key, ":")

	if len(parts) == 4 {
		namespace = parts[0]
	}

	return namespace
}
