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
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
	"github.com/lastbackend/lastbackend/pkg/util/stream"
)

const BUFFER_SIZE = 1024

func Manage(ctx context.Context, pod *types.Pod) error {
	log.Debugf("Provision pod: %s", pod.Meta.Name)

	//==========================================================================
	// Destroy pod =============================================================
	//==========================================================================

	// Call destroy pod
	if pod.State.Destroy {

		if task := envs.Get().GetState().Tasks().GetTask(pod); task != nil {
			log.Debugf("Cancel pod creating: %s", pod.Meta.Name)
			task.Cancel()
		}

		log.Debugf("Pod found and in destroy state > destroy it: %s", pod.Meta.Name)

		pods := envs.Get().GetState().Pods().GetPods()
		if p, ok := pods[pod.Meta.Name]; ok {
			Destroy(ctx, &p)
		}

		pod.MarkAsDestroyed()
		envs.Get().GetState().Pods().SetPod(pod)
		events.NewPodStateEvent(ctx, pod)

		return nil
	}

	//==========================================================================
	// Create pod ==============================================================
	//==========================================================================

	// Get pod list from current state
	pods := envs.Get().GetState().Pods().GetPods()

	if _, ok := pods[pod.Meta.Name]; !ok {
		log.Debugf("Pod not found > create it: %s", pod.Meta.Name)

		ctx, cancel := context.WithCancel(context.Background())
		envs.Get().GetState().Tasks().AddTask(pod, &types.NodeTask{Cancel: cancel})

		if err := Create(ctx, pod); err != nil {
			log.Errorf("Can not be create pod: %s", pod.Meta.Name)
			return err
		}

		return nil
	}

	//==========================================================================
	// Scale pod ===============================================================
	//==========================================================================

	if len(pod.Spec.Containers) != len(pods[pod.Meta.Name].Status.Containers) {

		log.Debugf("Pod containers not match: %d != %d",
			len(pod.Spec.Containers), len(pods[pod.Meta.Name].Spec.Containers))

		if task := envs.Get().GetState().Tasks().GetTask(pod); task != nil {
			log.Debugf("Cancel pod creating: %s", pod.Meta.Name)
			task.Cancel()
		}

		pod.MarkAsDestroyed()
		envs.Get().GetState().Pods().SetPod(pod)
		events.NewPodStateEvent(ctx, pod)
		Destroy(ctx, pod)

		ctx, cancel := context.WithCancel(context.Background())
		envs.Get().GetState().Tasks().AddTask(pod, &types.NodeTask{Cancel: cancel})

		if err := Create(ctx, pod); err != nil {
			log.Errorf("Can not be create pod: %s", pod.Meta.Name)
			return err
		}

		return nil
	}

	return nil
}

func Create(ctx context.Context, pod *types.Pod) error {

	var (
		err     error
		inspect = func(ctx context.Context, container *types.PodContainer) error {
			info, err := envs.Get().GetCri().ContainerInspect(ctx, container.ID)
			if err != nil {
				switch err {
				case context.Canceled:
					log.Errorf("Stop inspect container %s: %s", err.Error())
					pod.MarkAsDestroyed()
					envs.Get().GetState().Pods().SetPod(pod)
					events.NewPodStateEvent(ctx, pod)
					Clean(context.Background(), pod)
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
			if pod.Status.Network.PodIP == "" {
				pod.Status.Network.PodIP = info.Network.IPAddress
			}
			return nil
		}
	)

	log.Debugf("Create pod: %s", pod.Meta.Name)

	//==========================================================================
	// Pull image ==============================================================
	//==========================================================================

	pod.MarkAsPull()
	envs.Get().GetState().Pods().AddPod(pod)
	events.NewPodStateEvent(ctx, pod)

	log.Debugf("Have %d containers", len(pod.Spec.Containers))
	for _, c := range pod.Spec.Containers {
		log.Debug("Pull images for pod if needed")
		r, err := envs.Get().GetCri().ImagePull(ctx, &c.Image)
		if err != nil {
			log.Errorf("Can-not pull image: %s", err)
			pod.MarkAsError(err)
			envs.Get().GetState().Pods().SetPod(pod)
			events.NewPodStateEvent(ctx, pod)
			Clean(context.Background(), pod)
			return err
		}

		io.Copy(os.Stdout, r)
	}

	//==========================================================================
	// Run container ===========================================================
	//==========================================================================

	pod.MarkAsStarting()
	pod.Status.Steps[types.PodStepPull] = types.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}
	envs.Get().GetState().Pods().SetPod(pod)
	events.NewPodStateEvent(ctx, pod)

	for _, s := range pod.Spec.Containers {

		if s.Labels == nil {
			s.Labels = make(map[string]string)
		}

		s.Labels["LB"] = fmt.Sprintf("%s/%s/%s", pod.Meta.Namespace, pod.Meta.Deployment, pod.Meta.Name)

		//==========================================================================
		// Create container ========================================================
		//==========================================================================

		var c = new(types.PodContainer)
		c.ID, err = envs.Get().GetCri().ContainerCreate(ctx, &s)
		if err != nil {
			switch err {
			case context.Canceled:
				log.Errorf("Stop creating container: %s", err.Error())

				if err := inspect(context.Background(), c); err != nil {
					log.Errorf("Inspect container after create: err %s", err.Error())
				}

				pod.MarkAsDestroyed()
				pod.Status.Containers[c.ID] = c
				envs.Get().GetState().Pods().SetPod(pod)
				events.NewPodStateEvent(ctx, pod)
				Clean(context.Background(), pod)

				return nil
			}

			log.Errorf("Can-not create container: %s", err)
			c.State.Error = types.PodContainerStateError{
				Error:   true,
				Message: err.Error(),
				Exit: types.PodContainerStateExit{
					Timestamp: time.Now().UTC(),
				},
			}
			pod.Status.Containers[c.ID] = c

			pod.MarkAsError(err)
			envs.Get().GetState().Pods().SetPod(pod)
			events.NewPodStateEvent(ctx, pod)
			Clean(context.Background(), pod)
			return err
		}

		if err := inspect(context.Background(), c); err != nil {
			log.Errorf("Inspect container after create: err %s", err.Error())
			pod.MarkAsError(err)
			envs.Get().GetState().Pods().SetPod(pod)
			events.NewPodStateEvent(ctx, pod)
			Clean(context.Background(), pod)
			return err
		}

		//==========================================================================
		// Start container =========================================================
		//==========================================================================

		c.State.Created = types.PodContainerStateCreated{
			Created: time.Now().UTC(),
		}
		pod.Status.Containers[c.ID] = c
		envs.Get().GetState().Pods().SetPod(pod)
		events.NewPodStateEvent(ctx, pod)

		log.Debugf("Container created: %#v", c)

		if err := envs.Get().GetCri().ContainerStart(ctx, c.ID); err != nil {
			switch err {
			case context.Canceled:
				log.Errorf("Stop starting container err: %s", err.Error())

				if err := inspect(context.Background(), c); err != nil {
					log.Errorf("Inspect container after create: err %s", err.Error())
				}

				pod.Status.Containers[c.ID] = c
				pod.MarkAsDestroyed()
				envs.Get().GetState().Pods().SetPod(pod)
				events.NewPodStateEvent(ctx, pod)
				Clean(context.Background(), pod)
				return nil
			}
			log.Errorf("Can-not start container: %s", err)
			c.State.Error = types.PodContainerStateError{
				Error:   true,
				Message: err.Error(),
				Exit: types.PodContainerStateExit{
					Timestamp: time.Now().UTC(),
				},
			}

			pod.Status.Containers[c.ID] = c
			pod.MarkAsError(err)
			envs.Get().GetState().Pods().SetPod(pod)
			events.NewPodStateEvent(ctx, pod)
			Clean(context.Background(), pod)
			return err
		}

		if err := inspect(context.Background(), c); err != nil {
			log.Errorf("Inspect container after create: err %s", err.Error())
			pod.MarkAsError(err)
			envs.Get().GetState().Pods().SetPod(pod)
			events.NewPodStateEvent(ctx, pod)
			Clean(context.Background(), pod)
			return err
		}

		c.Ready = true
		c.State.Started = types.PodContainerStateStarted{
			Started:   true,
			Timestamp: time.Now().UTC(),
		}
		pod.Status.Containers[c.ID] = c
		envs.Get().GetState().Pods().SetPod(pod)
		events.NewPodStateEvent(ctx, pod)
	}

	pod.MarkAsRunning()
	pod.Status.Steps[types.PodStepReady] = types.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}

	pod.Status.Network.HostIP = envs.Get().GetCNI().Info(ctx).Addr
	envs.Get().GetState().Pods().SetPod(pod)
	events.NewPodStateEvent(ctx, pod)

	return nil
}

func Clean(ctx context.Context, pod *types.Pod) {

	for _, c := range pod.Status.Containers {
		log.Debugf("Remove unnecessary container: %s", c.ID)
		if err := envs.Get().GetCri().ContainerRemove(ctx, c.ID, true, true); err != nil {
			log.Warnf("Can-not remove unnecessary container %s: %s", c.ID, err)
		}
	}

	for _, c := range pod.Spec.Containers {
		log.Debugf("Try to clean image: %s", c.Image.Name)
		if err := envs.Get().GetCri().ImageRemove(ctx, c.Image.Name); err != nil {
			log.Warnf("Can-not remove unnecessary image %s: %s", c.Image.Name, err)
		}
	}
}

func Destroy(ctx context.Context, pod *types.Pod) {
	log.Debugf("Try to remove pod: %s", pod.Meta.Name)
	Clean(ctx, pod)
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

		pod := envs.Get().GetState().Pods().GetPod(c.Pod)

		if pod == nil {
			pod = types.NewPod()
		}

		pod.Meta.Name = c.Pod
		pod.Meta.Deployment = c.Deployment
		pod.Meta.Namespace = c.Namespace

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

		pod.Status.Containers[cs.ID] = cs
		envs.Get().GetState().Pods().SetPod(pod)
	}

	return nil
}

func Logs(ctx context.Context, id string, follow bool, s *stream.Stream) error {

	log.Debugf("Get container [%s] logs streaming", id)

	var (
		cri    = envs.Get().GetCri()
		buffer = make([]byte, BUFFER_SIZE)
		done   = make(chan bool, 1)
	)

	ctx, cfunc := context.WithCancel(ctx)

	req, err := cri.ContainerLogs(ctx, id, true, true, follow)
	if err != nil {
		log.Errorf("Error get logs stream %s", err)
		return err
	}
	defer func() {
		log.Debugf("Stop container [%s] logs streaming", id)
		s.Close()
		close(done)
	}()

	go func() {
		s.Done()
		cfunc()
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
				s.Flush()
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
