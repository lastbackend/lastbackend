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

package pod

import (
	"context"
	"io"
	"os"
	"time"

	"net/http"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
	"github.com/lastbackend/lastbackend/pkg/util/cleaner"
)

const (
	BUFFER_SIZE = 1024
	logLevel = 3
)

func Manage(ctx context.Context, key string, manifest *types.PodManifest) error {
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
			events.NewPodStatusEvent(ctx, key)

			return nil
		}

		log.V(logLevel).Debugf("Pod found > destroy it: %s", key)

		Destroy(ctx, key, p)

		p.SetDestroyed()
		envs.Get().GetState().Pods().SetPod(key, p)
		events.NewPodStatusEvent(ctx, key)
		return nil
	}

	//==========================================================================
	// Check containers pod status =============================================
	//==========================================================================

	// Get pod list from current state
	p := envs.Get().GetState().Pods().GetPod(key)
	if p != nil {
		if p.State != types.StateWarning  {
			events.NewPodStatusEvent(ctx, key)
			return nil
		}

		if p.State == types.StateError {
			Destroy(ctx, key, p)
		}
	}

	log.V(logLevel).Debugf("Pod not found > create it: %s", key)

	ctx, cancel := context.WithCancel(context.Background())
	envs.Get().GetState().Tasks().AddTask(key, &types.NodeTask{Cancel: cancel})

	status, err := Create(ctx, key, manifest)
	if err != nil {
		log.Errorf("Can not create pod: %s err: %s", key, err.Error())
		status.SetError(err)
	}

	envs.Get().GetState().Pods().SetPod(key, status)
	events.NewPodStatusEvent(ctx, key)
	return nil
}

func Create(ctx context.Context, key string, manifest *types.PodManifest) (*types.PodStatus, error) {

	var (
		err    error
		status = types.NewPodStatus()
	)

	log.V(logLevel).Debugf("Create pod: %s", key)

	//==========================================================================
	// Pull image ==============================================================
	//==========================================================================

	status.SetPull()

	envs.Get().GetState().Pods().AddPod(key, status)
	events.NewPodStatusEvent(ctx, key)

	log.V(logLevel).Debugf("Have %d containers", len(manifest.Template.Containers))
	for _, c := range manifest.Template.Containers {

		log.V(logLevel).Debug("Pull images for pod if needed")

		var secret = new(types.Secret)

		if c.Image.Secret != types.EmptyString {
			log.V(logLevel).Debug("Get secret info from api")
			vs, err := envs.Get().GetRestClient().Secret(c.Image.Secret).Get(ctx)
			if err != nil {
				log.Errorf("Can-not get secret from api: %s", err)
				status.SetError(err)
				Clean(context.Background(), status)
				return status, err
			}

			secret = vs.Decode()
		}

		r, err := envs.Get().GetCRI().ImagePull(ctx, &c.Image, secret)
		if err != nil {
			log.Errorf("Can-not pull image: %s", err)
			status.SetError(err)
			Clean(context.Background(), status)
			return status, err
		}

		io.Copy(os.Stdout, r)
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
	events.NewPodStatusEvent(ctx, key)

	var secrets = make(map[string]*types.Secret)
	for _, s := range manifest.Template.Containers {
		for _, e := range s.EnvVars {
			if e.From.Name != types.EmptyString {
				log.V(logLevel).Debug("Get secret info from api")

				vs, err := envs.Get().GetRestClient().Secret(e.From.Name).Get(ctx)
				if err != nil {
					log.Errorf("Can-not get secret from api: %s", err)
					status.SetError(err)
					Clean(context.Background(), status)
					return status, err
				}

				secret := vs.Decode()
				secrets[secret.Meta.Name] = secret
			}
		}
	}

	for _, s := range manifest.Template.Containers {

		//==========================================================================
		// Create container ========================================================
		//==========================================================================

		var c = new(types.PodContainer)
		c.ID, err = envs.Get().GetCRI().ContainerCreate(ctx, &s, secrets)
		if err != nil {
			switch err {
			case context.Canceled:
				log.Errorf("Stop creating container: %s", err.Error())
				Clean(context.Background(), status)
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
			Clean(context.Background(), status)
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

		if err := envs.Get().GetCRI().ContainerStart(ctx, c.ID); err != nil {
			switch err {
			case context.Canceled:
				log.Errorf("Stop starting container err: %s", err.Error())
				Clean(context.Background(), status)
				return status, nil
			}

			log.Errorf("Can-not start container: %s", err)
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

	status.Network.HostIP = envs.Get().GetCNI().Info(ctx).Addr
	envs.Get().GetState().Pods().SetPod(key, status)
	return status, nil
}

func Clean(ctx context.Context, status *types.PodStatus) {

	for _, c := range status.Containers {
		log.V(logLevel).Debugf("Remove unnecessary container: %s", c.ID)
		if err := envs.Get().GetCRI().ContainerRemove(ctx, c.ID, true, true); err != nil {
			log.Warnf("Can-not remove unnecessary container %s: %s", c.ID, err)
		}
	}

	for _, c := range status.Containers {
		log.V(logLevel).Debugf("Try to clean image: %s", c.Image.Name)
		if err := envs.Get().GetCRI().ImageRemove(ctx, c.Image.Name); err != nil {
			log.Warnf("Can-not remove unnecessary image %s: %s", c.Image.Name, err)
		}
	}
}

func Destroy(ctx context.Context, pod string, status *types.PodStatus) {
	log.V(logLevel).Debugf("Try to remove pod: %s", pod)
	Clean(ctx, status)
	envs.Get().GetState().Pods().DelPod(pod)
}

func Restore(ctx context.Context) error {

	log.V(logLevel).Debug("Runtime restore state")

	cl, err := envs.Get().GetCRI().ContainerList(ctx, true)
	if err != nil {
		log.Errorf("Pods restore error: %s", err)
		return err
	}

	for _, c := range cl {

		log.V(logLevel).Debugf("Pod [%s] > container restore %s", c.Pod, c.ID)

		status := envs.Get().GetState().Pods().GetPod(c.Pod)
		if status == nil {
			status = types.NewPodStatus()
		}

		key := c.Pod

		cs := &types.PodContainer{
			ID: c.ID,
			Image: types.PodContainerImage{
				Name: c.Image,
			},
		}

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

		log.V(logLevel).Debugf("Container restored %s", c.ID)
		envs.Get().GetState().Pods().SetPod(key, status)
		log.V(logLevel).Debugf("Pod restored %#v", status)
	}

	return nil
}

func Logs(ctx context.Context, id string, follow bool, s io.Writer, doneChan chan bool) error {

	log.V(logLevel).Debugf("Get container [%s] logs streaming", id)

	var (
		cri    = envs.Get().GetCRI()
		buffer = make([]byte, BUFFER_SIZE)
		done   = make(chan bool, 1)
	)

	req, err := cri.ContainerLogs(ctx, id, true, true, follow)
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

func containerInspect(ctx context.Context, status *types.PodStatus, container *types.PodContainer) error {
	info, err := envs.Get().GetCRI().ContainerInspect(ctx, container.ID)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Errorf("Stop inspect container %s: %s", err.Error())
			return nil
		}
		log.Errorf("Can-not inspect container: %s", err)
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
