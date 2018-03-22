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

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
)

const BUFFER_SIZE = 1024

func Manage(ctx context.Context, key string, spec *types.PodSpec) error {
	log.Debugf("Provision pod: %s", key)

	//==========================================================================
	// Destroy pod =============================================================
	//==========================================================================

	// Call destroy pod
	if spec.State.Destroy {

		if task := envs.Get().GetState().Tasks().GetTask(key); task != nil {
			log.Debugf("Cancel pod creating: %s", key)
			task.Cancel()
		}

		log.Debugf("Pod found and in destroy state > destroy it: %s", key)

		p := envs.Get().GetState().Pods().GetPod(key)
		if p == nil {
			return errors.New(errors.PodNotFound)
		}

		Destroy(ctx, key, p)

		p.SetDestroyed()
		events.NewPodStatusEvent(ctx, key, p)
		return nil
	}

	//==========================================================================
	// Check containers pod status =============================================
	//==========================================================================

	// Get pod list from current state
	p := envs.Get().GetState().Pods().GetPod(key)
	if p != nil {
		// TODO: check pod status
		// if not match - recreate pod or try to fix
		return nil
	}

	log.Debugf("Pod not found > create it: %s", key)

	ctx, cancel := context.WithCancel(context.Background())
	envs.Get().GetState().Tasks().AddTask(key, &types.NodeTask{Cancel: cancel})

	status, err := Create(ctx, key, spec)
	if err != nil {
		log.Errorf("Can not create pod: %s err: %s", key, err.Error())
		status.SetError(err)
	}

	events.NewPodStatusEvent(ctx, key, p)
	return nil
}

func Create(ctx context.Context, key string, spec *types.PodSpec) (*types.PodStatus, error) {

	var (
		err    error
		status = new(types.PodStatus)
	)

	log.Debugf("Create pod: %s", key)

	//==========================================================================
	// Pull image ==============================================================
	//==========================================================================

	status.SetPull()

	envs.Get().GetState().Pods().AddPod(key, status)
	events.NewPodStatusEvent(ctx, key, status)

	log.Debugf("Have %d containers", len(spec.Template.Containers))
	for _, c := range spec.Template.Containers {

		log.Debug("Pull images for pod if needed")
		r, err := envs.Get().GetCri().ImagePull(ctx, &c.Image)
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
	events.NewPodStatusEvent(ctx, key, status)

	for _, s := range spec.Template.Containers {

		//==========================================================================
		// Create container ========================================================
		//==========================================================================

		var c = new(types.PodContainer)
		c.ID, err = envs.Get().GetCri().ContainerCreate(ctx, &s)
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
		log.Debugf("Container created: %#v", c)

		if err := envs.Get().GetCri().ContainerStart(ctx, c.ID); err != nil {
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
		log.Debugf("Remove unnecessary container: %s", c.ID)
		if err := envs.Get().GetCri().ContainerRemove(ctx, c.ID, true, true); err != nil {
			log.Warnf("Can-not remove unnecessary container %s: %s", c.ID, err)
		}
	}

	for _, c := range status.Containers {
		log.Debugf("Try to clean image: %s", c.Image.Name)
		if err := envs.Get().GetCri().ImageRemove(ctx, c.Image.Name); err != nil {
			log.Warnf("Can-not remove unnecessary image %s: %s", c.Image.Name, err)
		}
	}
}

func Destroy(ctx context.Context, pod string, status *types.PodStatus) {
	log.Debugf("Try to remove pod: %s", pod)
	Clean(ctx, status)
	envs.Get().GetState().Pods().DelPod(pod)
}

func Restore(ctx context.Context) error {

	log.Debug("Runtime restore state")
	cl, err := envs.Get().GetCri().ContainerList(ctx, true)
	if err != nil {
		log.Errorf("Pods restore error: %s", err)
		return err
	}

	for _, c := range cl {
		log.Debugf("Container restore %s", c.ID)
		envs.Get().GetState().Pods().AddContainer(c)

		status := envs.Get().GetState().Pods().GetPod(c.Pod)

		if status == nil {
			status = new(types.PodStatus)
		}

		key := c.Pod

		cs := &types.PodContainer{
			ID: c.ID,
			Image: types.PodContainerImage{
				Name: c.Image,
			},
		}

		if c.Status == types.StateStopped {
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
		envs.Get().GetState().Pods().SetPod(key, status)
	}

	return nil
}

func Logs(ctx context.Context, id string, follow bool, s io.Writer) error {

	log.Debugf("Get container [%s] logs streaming", id)

	var (
		cri    = envs.Get().GetCri()
		buffer = make([]byte, BUFFER_SIZE)
		done   = make(chan bool, 1)
	)

	req, err := cri.ContainerLogs(ctx, id, true, true, follow)
	if err != nil {
		log.Errorf("Error get logs stream %s", err)
		return err
	}
	defer func() {
		log.Debugf("Stop container [%s] logs streaming", id)
		ctx.Done()
		close(done)
		req.Close()
	}()

	go func() {
		for {
			n, err := req.Read(buffer)
			if err != nil {

				if err == context.Canceled {
					log.Debug("Stream is canceled")
					return
				}

				log.Errorf("Error read bytes from stream %s", err)
				done <- true
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
	}()

	<-done
	return nil
}

func containerInspect(ctx context.Context, status *types.PodStatus, container *types.PodContainer) error {
	info, err := envs.Get().GetCri().ContainerInspect(ctx, container.ID)
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
		if info.Status == types.StateStopped {
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
