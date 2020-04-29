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
	"bytes"
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/tools/logger"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/cleaner"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
)

const (
	logPodPrefix               = "node:runtime:pod:>"
	defaultRootLocalStorgePath = "/var/lib/lastbackend/runtime/"

	BufferSize = 1024
)

// tplScript is a helper script this is added to the template the commands.
const logScript = `
echo ""
echo "[task: %s]"
echo ""
set -eux
%s
`

// logScript is a helper script that is added to
// the build script to logging run a command.
const tplScript = `
%s
`

func (r Runtime) PodManage(ctx context.Context, key string, manifest *models.PodManifest) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s provision pod: %s", logPodPrefix, key)

	//==========================================================================
	// Destroy pod =============================================================
	//==========================================================================

	// Call destroy pod
	if manifest.State.Destroy {

		if task := r.state.Tasks().GetTask(key); task != nil {
			log.Debugf("%s cancel pod creating: %s", logPodPrefix, key)
			task.Cancel()
		}

		p := r.state.Pods().GetPod(key)
		if p == nil {

			ps := models.NewPodStatus()
			ps.SetDestroyed()
			r.state.Pods().AddPod(key, ps)

			return nil
		}

		log.Debugf("%s pod found > destroy it: %s", logPodPrefix, key)

		r.PodDestroy(ctx, key, p)

		p.SetDestroyed()
		r.state.Pods().SetPod(key, p)
		return nil
	}

	//==========================================================================
	// Check containers pod status =============================================
	//==========================================================================

	// Get pod list from current state
	p := r.state.Pods().GetPod(key)
	if p != nil {

		// restore pov volume claims
		r.podVolumeClaimRestore(key, manifest)

		switch true {
		case !r.PodSpecCheck(ctx, key, manifest) || len(manifest.Runtime.Tasks) > 0:
			r.PodDestroy(ctx, key, p)
			break
		case !r.PodVolumesCheck(ctx, key, manifest.Template.Volumes):
			log.Debugf("%s volumes data changed: %s", logPodPrefix, key)
			for _, v := range manifest.Template.Volumes {

				if v.Volume.Name != models.EmptyString {

					log.Debugf("%s attach volume %s for pod %s", logPodPrefix, v.Name, key)

					pv, err := r.PodVolumeAttach(ctx, key, v)
					if err != nil {
						log.Errorf("%s can not attach volume for pod: %s", logPodPrefix, err.Error())
						return err
					}

					p.Volumes[v.Name] = pv

				} else {

					log.Debugf("%s create pod volume %s for pod %s", logPodPrefix, v.Name, key)

					var name = r.podVolumeKeyCreate(key, v.Name)

					vol := r.state.Volumes().GetVolume(name)

					if vol == nil {
						log.Debugf("%s update pod volume: volume not found: create %s: %s", logPodPrefix, key, v.Name)

						vs, err := r.PodVolumeCreate(ctx, key, v)
						if err != nil {
							log.Errorf("%s can not update volume data: %s", logPodPrefix, err.Error())
							return err
						}

						pv := &models.VolumeClaim{
							Name:   r.podVolumeClaimNameCreate(key, v.Name),
							Volume: name,
							Path:   vs.Status.Path,
						}

						r.state.Volumes().SetClaim(pv.Name, pv)
						p.Volumes[pv.Name] = pv

					} else {

						_, err := r.PodVolumeUpdate(ctx, key, v)
						if err != nil {
							log.Errorf("%s can not update volume data: %s", logPodPrefix, err.Error())
							return err
						}
					}

				}

			}
			return r.PodRestart(ctx, key)
		default:
			return nil
		}
	}

	log.Debugf("%s pod not found > create it: %s", logPodPrefix, key)

	ctx, cancel := context.WithCancel(context.Background())
	r.state.Tasks().AddTask(key, &models.NodeTask{Cancel: cancel})

	go func() {

		status, err := r.PodCreate(ctx, key, manifest)
		if err != nil {
			log.Errorf("%s can not create pod: %s err: %s", logPodPrefix, key, err.Error())
			status.SetError(err)
		}

		r.state.Pods().SetPod(key, status)
	}()

	return nil
}

func (r Runtime) PodRestart(ctx context.Context, key string) error {

	pod := r.state.Pods().GetPod(key)
	if pod == nil {
		return errors.New("pod not found")
	}

	cri := r.cri

	for _, c := range pod.Runtime.Services {
		if err := cri.Restart(ctx, c.ID, nil); err != nil {
			return err
		}
	}

	return nil
}

func (r Runtime) PodCreate(ctx context.Context, key string, manifest *models.PodManifest) (*models.PodStatus, error) {
	log := logger.WithContext(context.Background())
	var (
		status    = models.NewPodStatus()
		namespace = r.getPodNamespace(key)

		setError = func(err error) (*models.PodStatus, error) {
			log.Errorf("%s can not pull image: %s", logPodPrefix, err)
			status.SetError(err)
			r.state.Pods().SetPod(key, status)
			r.PodClean(ctx, status)
			return status, err
		}
	)

	log.Debugf("%s create pod: %s", logPodPrefix, key)

	//==========================================================================
	// Pull image ==============================================================
	//==========================================================================

	status.SetPull()

	r.state.Pods().AddPod(key, status)

	log.Debugf("%s have %d volumes", logPodPrefix, len(manifest.Template.Volumes))
	for _, v := range manifest.Template.Volumes {

		var name string
		if v.Volume.Name != models.EmptyString {
			name = fmt.Sprintf("%s:%s", r.getPodNamespace(key), v.Volume.Name)
		} else {
			name = r.podVolumeKeyCreate(key, v.Name)
		}

		vol := r.state.Volumes().GetVolume(name)
		if vol == nil {
			log.Debugf("%s update pod volume: volume not found: create %s: %s", logPodPrefix, key, v.Name)

			vs, err := r.PodVolumeCreate(ctx, key, v)
			if err != nil {
				log.Errorf("%s can not update volume data: %s", logPodPrefix, err.Error())
				return status, err
			}

			pv := &models.VolumeClaim{
				Name:   r.podVolumeClaimNameCreate(key, v.Name),
				Volume: name,
				Path:   vs.Status.Path,
			}

			r.state.Volumes().SetClaim(pv.Name, pv)
			status.Volumes[pv.Name] = pv

		} else {

			_, err := r.PodVolumeUpdate(ctx, key, v)
			if err != nil {
				log.Errorf("%s can not update volume data: %s", logPodPrefix, err.Error())
				return status, err
			}

			claim := r.state.Volumes().GetClaim(r.podVolumeClaimNameCreate(key, v.Name))
			if claim == nil {
				pv := &models.VolumeClaim{
					Name:   r.podVolumeClaimNameCreate(key, v.Name),
					Volume: name,
					Path:   vol.Status.Path,
				}

				r.state.Volumes().SetClaim(pv.Name, pv)
				status.Volumes[pv.Name] = pv
			}

		}

		r.state.Pods().SetPod(key, status)
	}

	if len(manifest.Runtime.Tasks) > 0 {
		for _, t := range manifest.Runtime.Tasks {
			pst := new(models.PodStatusPipelineStep)
			pst.Status = models.StateProvision
			pst.Error = false
			status.Runtime.Pipeline[t.Name] = pst
		}

		r.state.Pods().SetPod(key, status)
	}

	log.Debugf("%s have %d images", logPodPrefix, len(manifest.Template.Containers))

	for _, c := range manifest.Template.Containers {
		log.Debugf("%s pull image %s for pod if needed", logPodPrefix, c.Image.Name)
		if err := r.ImagePull(ctx, namespace, &c.Image); err != nil {
			log.Errorf("%s can not pull image: %s", logPodPrefix, err.Error())
			return setError(err)
		}
	}

	//==========================================================================
	// Run container ===========================================================
	//==========================================================================

	status.SetStarting()
	status.Steps[models.StepPull] = models.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}

	r.state.Pods().SetPod(key, status)

	for _, s := range manifest.Template.Containers {
		for _, e := range s.EnvVars {
			if e.Secret.Name != models.EmptyString {
				log.Debugf("%s get secret info from api", logPodPrefix)
				if err := r.SecretCreate(ctx, namespace, e.Secret.Name); err != nil {
					log.Errorf("%s can not fetch secret from api", logPodPrefix)
				}
			}
		}
	}

	var (
		primary  string
		services = make([]*models.ContainerManifest, 0)
	)

	if len(manifest.Runtime.Services) == 0 {

		if len(manifest.Runtime.Tasks) > 0 {
			tpl := models.GetPauseContainerTemplate()
			manifest.Runtime.Services = append(manifest.Runtime.Services, tpl.Name)
			manifest.Template.Containers = append(manifest.Template.Containers, tpl)
		} else {
			for _, s := range manifest.Template.Containers {
				m, err := r.containerManifestCreate(ctx, key, s)
				if err != nil {
					log.Errorf("%s can not create container manifest from spec: %s", logPodPrefix, err.Error())
					return setError(err)
				}

				services = append(services, m)
			}
		}
	}

	if len(manifest.Runtime.Services) != 0 {

		for _, name := range manifest.Runtime.Services {
			for _, s := range manifest.Template.Containers {

				if s.Name != name {
					continue
				}

				m, err := r.containerManifestCreate(ctx, key, s)
				if err != nil {
					log.Errorf("%s can not create container manifest from spec: %s", logPodPrefix, err.Error())
					return setError(err)
				}

				services = append(services, m)
			}
		}

	}

	// run services
	for _, svc := range services {

		if primary != models.EmptyString {
			svc.Network.Mode = fmt.Sprintf("container:%s", primary)
		} else {
			primary = svc.Name
			// TODO: Set default container extra hosts
			//svc.ExtraHosts = util.RemoveDuplicates(append(svc.ExtraHosts, envs.Get().GetConfig().Container.ExtraHosts...))
		}

		if err := r.serviceStart(ctx, key, svc, status); err != nil {
			log.Errorf("%s can not start service: %s", logPodPrefix, err.Error())
			return status, err
		}

	}

	status.SetRunning()
	status.Steps[models.StepReady] = models.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}

	r.state.Pods().SetPod(key, status)

	// run tasks
	for _, t := range manifest.Runtime.Tasks {

		log.Debugf("%s start task %s", logPodPrefix, t.Name)

		var f, e bool

		for _, s := range manifest.Template.Containers {

			if s.Name != t.Container {
				continue
			}

			f = true
			spec := *s

			if len(t.EnvVars) > 0 {
				for _, te := range t.EnvVars {
					var f = false
					for _, se := range spec.EnvVars {
						if te.Name == se.Name {
							se.Value = te.Value
							se.Secret = te.Secret
							se.Config = te.Config
							f = true
						}
					}
					if !f {
						spec.EnvVars = append(spec.EnvVars, te)
					}
				}
			}

			var buf bytes.Buffer
			for _, command := range t.Commands {
				buf.WriteString(fmt.Sprintf(tplScript, command))
			}

			escaped := fmt.Sprintf("%q", t.Name)
			escaped = strings.Replace(escaped, "$", `\$`, -1)
			script := fmt.Sprintf(logScript, escaped, buf.String())

			rootPath := defaultRootLocalStorgePath
			if len(r.config.WorkDir) != 0 {
				rootPath = r.config.WorkDir
			}

			filepath := path.Join(rootPath, strings.Replace(key, ":", "-", -1), "init")

			log.Debugf("%s runtime volume create: %s", logPodPrefix, filepath)

			err := r.podLocalFileCreate(filepath, script)
			if err != nil {
				log.Errorf("%s can not create runtime volume err: %s", logPodPrefix, err.Error())
				return setError(err)
			}

			log.Debugf("%s container manifest create", logPodPrefix)

			m, err := r.containerManifestCreate(ctx, key, &spec)
			if err != nil {
				log.Errorf("%s can not create container manifest from spec: %s", logPodPrefix, err.Error())
				return setError(err)
			}

			if primary != models.EmptyString {
				m.Network.Mode = fmt.Sprintf("container:%s", primary)
			}

			m.Name = ""

			m.Exec.Command = []string{"/usr/local/bin/lb_entrypoint"}
			m.Exec.Entrypoint = []string{"/bin/sh"}
			m.Binds = append(m.Binds, fmt.Sprintf("%s:%s:ro", filepath, "/usr/local/bin/lb_entrypoint"))

			if err := r.taskExecute(ctx, key, t, *m, status); err != nil {
				log.Errorf("%s can not execute task: %s", logPodPrefix, err.Error())
				return setError(err)
			}

			for _, s := range status.Runtime.Pipeline {
				if s.Error {
					e = true
					status.SetError(errors.New(s.Message))
					break
				}
			}

		}

		if e || !f {
			break
		}
	}

	if len(manifest.Runtime.Tasks) > 0 {
		r.PodExit(ctx, key, status, true)
	}

	return status, nil
}

func (r Runtime) PodClean(ctx context.Context, status *models.PodStatus) {
	log := logger.WithContext(context.Background())
	for _, c := range status.Runtime.Services {
		log.Debugf("%s remove unnecessary container: %s", logPodPrefix, c.ID)
		if err := r.cri.Remove(ctx, c.ID, true, true); err != nil {
			log.Warnf("%s can-not remove unnecessary container %s: %s", logPodPrefix, c.ID, err)
		}
	}

	for _, c := range status.Runtime.Services {
		log.Debugf("%s try to clean image: %s", logPodPrefix, c.Image.Name)
		//if err := ImageRemove(ctx, c.Image.Name); err != nil {
		//	log.Errorf("%s can not remove image: %s", logPodPrefix, err.Error())
		//	continue
		//}
	}

}

func (r Runtime) PodExit(ctx context.Context, pod string, status *models.PodStatus, clean bool) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s exit pod: %s", logPodPrefix, pod)

	timeout := time.Duration(3 * time.Second)
	for _, c := range status.Runtime.Services {

		var attempts = 5
		for i := 1; i <= attempts; i++ {
			if err := r.cri.Stop(ctx, c.ID, &timeout); err != nil {
				// TODO: check container not found error
				log.Warnf("%s can-not stop container %s: %s", logPodPrefix, c.ID, err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		c.State.Stopped = models.PodContainerStateStopped{
			Stopped: true,
			Exit: models.PodContainerStateExit{
				Code:      0,
				Timestamp: time.Now(),
			},
		}
	}

	status.Steps[models.StateExited] = models.PodStep{
		Ready:     true,
		Timestamp: time.Now(),
	}

	if status.Status != models.StateError {
		status.SetExited()
	}

	r.state.Pods().SetPod(pod, status)

	if clean {
		r.PodClean(ctx, status)
		return
	}
}

func (r Runtime) PodDestroy(ctx context.Context, pod string, status *models.PodStatus) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s try to remove pod: %s", logPodPrefix, pod)

	r.PodClean(ctx, status)
	r.state.Pods().DelPod(pod)
	for _, v := range status.Volumes {
		if err := r.PodVolumeDestroy(ctx, pod, v.Name); err != nil {
			log.Errorf("%s can not destroy pod: %s", logPodPrefix, err.Error())
		}
	}

	rootPath := defaultRootLocalStorgePath
	if len(r.config.WorkDir) != 0 {
		rootPath = r.config.WorkDir
	}

	dirPath := path.Join(rootPath, strings.Replace(pod, ":", "-", -1), "init")

	log.Debugf("%s runtime volume remove: %s", logPodPrefix, dirPath)

	if err := r.podLocalFileDestroy(dirPath); err != nil {
		log.Errorf("%s can not destroy runtime volume path: %s", logPodPrefix, err.Error())
	}
}

func (r Runtime) PodRestore(ctx context.Context) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s runtime restore state", logPodPrefix)

	cl, err := r.cri.List(ctx, true)
	if err != nil {
		log.Errorf("%s pods restore error: %s", logPodPrefix, err)
		return err
	}

	for _, c := range cl {

		if c.Pod == models.EmptyString {
			continue
		}

		log.Debugf("%s pod [%s] > container restore %s", logPodPrefix, c.Pod, c.ID)

		status := r.state.Pods().GetPod(c.Pod)
		if status == nil {
			status = models.NewPodStatus()
		}

		key := c.Pod

		cs := &models.PodContainer{
			ID:   c.ID,
			Name: c.Name,
			Image: models.PodContainerImage{
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
		case models.StateCreated:
			cs.State = models.PodContainerState{
				Created: models.PodContainerStateCreated{
					Created: time.Now().UTC(),
				},
			}
		case models.StateStarted:
			cs.State = models.PodContainerState{
				Started: models.PodContainerStateStarted{
					Started:   true,
					Timestamp: time.Now().UTC(),
				},
			}
			cs.State.Stopped.Stopped = false
		case models.StatusStopped:
			cs.State.Stopped.Stopped = true
			cs.State.Stopped.Exit = models.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
		case models.StateError:

			cs.State.Error.Error = true
			cs.State.Error.Message = c.Status
			cs.State.Error.Exit = models.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
			cs.State.Stopped.Stopped = false
			cs.State.Stopped.Exit = models.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
		}

		if c.Status == models.StatusStopped {
			cs.State.Stopped = models.PodContainerStateStopped{
				Stopped: true,
				Exit: models.PodContainerStateExit{
					Code:      c.ExitCode,
					Timestamp: time.Now().UTC(),
				},
			}
		}

		cs.Ready = true
		status.Runtime.Services[cs.ID] = cs
		status.Network.PodIP = c.Network.IPAddress

		log.Debugf("%s container restored %s", logPodPrefix, c.ID)
		r.state.Pods().SetPod(key, status)
		log.Debugf("%s Pod restored %s: %s", key, status.State)
	}

	return nil
}

func (r Runtime) PodLogs(ctx context.Context, id string, follow bool, s io.Writer, doneChan chan bool) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s get container [%s] logs streaming", logPodPrefix, id)

	var (
		cri    = r.cri
		buffer = make([]byte, BufferSize)
		done   = make(chan bool, 1)
	)

	req, err := cri.Logs(ctx, id, true, true, follow)
	if err != nil {
		log.Errorf("%s error get logs stream %s", logPodPrefix, err)
		return err
	}
	defer func() {
		log.Debugf("%s stop container [%s] logs streaming", logPodPrefix, id)
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
						log.Debugf("%s Stream is canceled", logPodPrefix)
						return
					}

					log.Errorf("%s read bytes from stream err %s", logPodPrefix, err)
					doneChan <- true
					return
				}

				_, err = func(p []byte) (n int, err error) {
					n, err = s.Write(p)
					if err != nil {
						log.Errorf("%s write bytes to stream err %s", logPodPrefix, err)
						return n, err
					}

					if f, ok := s.(http.Flusher); ok {
						f.Flush()
					}
					return n, nil
				}(buffer[0:n])

				if err != nil {
					log.Errorf("%s write to stream err: %s", logPodPrefix, err)
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

func (r Runtime) PodSpecCheck(ctx context.Context, key string, manifest *models.PodManifest) bool {
	log := logger.WithContext(context.Background())
	log.Infof("%s pod check spec pod: %s", logPodPrefix, key)

	state := r.state.Pods().GetPod(key)

	var statec = make(map[string]*models.ContainerManifest, 0)
	var specc = make(map[string]*models.ContainerManifest, 0)

	for _, c := range manifest.Template.Containers {
		mf, err := r.containerManifestCreate(ctx, key, c)
		if err != nil {
			return false
		}
		specc[mf.Name] = mf
	}

	for _, c := range state.Runtime.Services {
		statec[c.Name] = c.GetManifest()
	}

	if len(statec) != len(specc) {
		log.Debugf("%s container spec count not equal not exists: %d != %d", logPodPrefix, len(statec), len(specc))
		return false
	}

	for n, mf := range specc {

		if _, ok := statec[n]; !ok {
			log.Debugf("%s container spec not exists: %s", logPodPrefix, n)
			return false
		}

		// check image

		c := statec[n]

		if c.Image != mf.Image {
			log.Debugf("%s images not equal: %s != %s", logPodPrefix, c.Image, mf.Image)
			return false
		}

		img := r.state.Images().GetImage(c.Image)
		if img == nil {
			log.Debugf("%s image not found in state: %s", logPodPrefix, mf.Image)
			return false
		}

		if len(mf.Exec.Command) == 0 {
			if strings.Join(c.Exec.Command, " ") != strings.Join(img.Status.Container.Exec.Command, " ") {
				log.Debugf("%s cmd different with img cmd: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Command, " "),
					strings.Join(img.Status.Container.Exec.Command, " "))
				return false
			}
		} else {
			if strings.Join(c.Exec.Command, " ") != strings.Join(mf.Exec.Command, " ") {
				log.Debugf("%s cmd different with manifest cmd: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Command, " "),
					strings.Join(mf.Exec.Command, " "))
				return false
			}
		}

		if len(mf.Exec.Entrypoint) == 0 {
			if strings.Join(c.Exec.Entrypoint, " ") != strings.Join(img.Status.Container.Exec.Entrypoint, " ") {
				log.Debugf("%s entrypoint changed: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Entrypoint, " "),
					strings.Join(img.Status.Container.Exec.Entrypoint, " "))
				return false
			}
		} else {
			if strings.Join(c.Exec.Entrypoint, " ") != strings.Join(mf.Exec.Entrypoint, " ") {
				log.Debugf("%s entrypoint changed: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Entrypoint, " "),
					strings.Join(mf.Exec.Entrypoint, " "))
				return false
			}
		}

		if mf.Exec.Workdir == models.EmptyString {
			if c.Exec.Workdir != img.Status.Container.Exec.Workdir {
				log.Debugf("%s workdir changed: %s != %s", logPodPrefix, c.Exec.Workdir, img.Status.Container.Exec.Workdir)
				return false
			}
		} else {
			if c.Exec.Workdir != mf.Exec.Workdir {
				log.Debugf("%s workdir changed: %s != %s", logPodPrefix, c.Exec.Workdir, mf.Exec.Workdir)
				return false
			}
		}

		if len(mf.Exec.Args) != 0 {
			if strings.Join(c.Exec.Args, " ") != strings.Join(mf.Exec.Args, " ") {
				log.Debugf("%s args changed: %s != %s", logPodPrefix,
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
				log.Debugf("%s env not found: %s", logPodPrefix, e)
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
				log.Debugf("%s env is unnecessary: %s", logPodPrefix, e)
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
				log.Debugf("%s bind not found: %s", logPodPrefix, e)
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
				log.Debugf("%s bind is unnecessary: %s", logPodPrefix, e)
				return false
			}
		}

		// Check ports
		for _, e := range mf.Ports {
			var f = false
			for _, ie := range c.Ports {
				if e.HostIP != models.EmptyString {
					if e.HostIP != ie.HostIP {
						log.Debugf("%s port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s", logPodPrefix,
							e.HostIP, e.HostPort, e.ContainerPort, e.Protocol,
							ie.HostIP, ie.HostPort, ie.ContainerPort, ie.Protocol)
						return false
					}
				}

				if e.Protocol != models.EmptyString {
					if e.Protocol != ie.Protocol {
						log.Debugf("%s port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s", logPodPrefix,
							e.HostIP, e.HostPort, e.ContainerPort, e.Protocol,
							ie.HostIP, ie.HostPort, ie.ContainerPort, ie.Protocol)
						return false
					}
				} else {
					if ie.Protocol != "tcp" {
						log.Debugf("%s port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s", logPodPrefix,
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
				log.Debugf("%s port map not found: \t\t %s:%d:%d/%s ", logPodPrefix,
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
				log.Debugf("%s port map is unnecessary: %d", logPodPrefix, e.HostPort)
				return false
			}
		}

		if mf.RestartPolicy.Policy != c.RestartPolicy.Policy ||
			mf.RestartPolicy.Attempt != c.RestartPolicy.Attempt {

			log.Debugf("%s restart policy changed: %s:%d => %s:%d", logPodPrefix,
				c.RestartPolicy.Policy, c.RestartPolicy.Attempt,
				mf.RestartPolicy.Policy, mf.RestartPolicy.Attempt)
			return false
		}

	}

	return true

}

func (r Runtime) PodVolumesCheck(ctx context.Context, pod string, spec []*models.SpecTemplateVolume) bool {
	log := logger.WithContext(context.Background())
	log.Debugf("%s check pod volumes: %s: %d", logPodPrefix, pod, len(spec))

	for _, v := range spec {

		if v.Volume.Name != models.EmptyString {
			continue
		}

		name := r.podVolumeKeyCreate(pod, v.Name)

		if v.Config.Name != models.EmptyString && len(v.Config.Binds) > 0 {
			equal, err := r.VolumeCheckConfigData(ctx, name, v.Config.Name)
			if err != nil {
				return false
			}
			return equal
		}

		if v.Secret.Name != models.EmptyString && len(v.Secret.Binds) > 0 {
			equal, err := r.VolumeCheckSecretData(ctx, name, v.Config.Name)
			if err != nil {
				return false
			}
			return equal
		}
	}

	return true
}

func (r Runtime) PodVolumeUpdate(ctx context.Context, pod string, spec *models.SpecTemplateVolume) (*models.VolumeStatus, error) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s update pod volume: %s: %s", logPodPrefix, pod, spec.Name)

	path := strings.Replace(pod, ":", "-", -1)
	path = fmt.Sprintf("%s-%s", path, spec.Name)

	var (
		name = r.podVolumeKeyCreate(pod, spec.Name)
	)

	status := r.state.Volumes().GetVolume(name)

	if spec.Secret.Name != models.EmptyString && len(spec.Secret.Binds) > 0 {
		if err := r.VolumeSetSecretData(ctx, name, spec.Secret.Name); err != nil {
			log.Errorf("%s can not set config data to volume: %s", logPodPrefix, err.Error())
			return status, err
		}
	}

	if spec.Secret.Name == models.EmptyString && spec.Config.Name != models.EmptyString && len(spec.Config.Binds) > 0 {
		if err := r.VolumeSetConfigData(ctx, name, spec.Config.Name); err != nil {
			log.Errorf("%s can not set config data to volume: %s", logPodPrefix, err.Error())
			return status, err
		}
	}

	return status, nil
}

func (r Runtime) PodVolumeAttach(ctx context.Context, pod string, spec *models.SpecTemplateVolume) (*models.VolumeClaim, error) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s attach pod volume: %s: %s", logPodPrefix, pod, spec.Name)

	var name = fmt.Sprintf("%s:%s", r.getPodNamespace(pod), spec.Name)

	volume := r.state.Volumes().GetVolume(name)
	if volume == nil {
		return nil, errors.New("volume not found on node")
	}

	pv := &models.VolumeClaim{
		Name:   r.podVolumeClaimNameCreate(pod, spec.Name),
		Volume: name,
		Path:   volume.Status.Path,
	}

	r.state.Volumes().SetClaim(pv.Name, pv)

	return pv, nil
}

func (r Runtime) PodVolumeCreate(ctx context.Context, pod string, spec *models.SpecTemplateVolume) (*models.VolumeStatus, error) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s create pod volume: %s:%s", logPodPrefix, pod, spec.Name)

	hostPath := strings.Replace(pod, ":", "-", -1)
	hostPath = fmt.Sprintf("%s-%s", hostPath, spec.Name)

	var (
		name = r.podVolumeKeyCreate(pod, spec.Name)
		vm   = models.VolumeManifest{
			HostPath: hostPath,
			Type:     models.KindVolumeHostDir,
		}
	)

	st, err := r.VolumeCreate(ctx, name, &vm)
	if err != nil {
		log.Errorf("%s can not create pod volume: %s", logPodPrefix, err.Error())
		return nil, err
	}

	if spec.Secret.Name != models.EmptyString && len(spec.Secret.Binds) > 0 {
		if err := r.VolumeSetSecretData(ctx, name, spec.Secret.Name); err != nil {
			log.Errorf("%s can not set secret data to volume: %s", logPodPrefix, err.Error())
			return st, err
		}
	}

	if spec.Secret.Name == models.EmptyString && spec.Config.Name != models.EmptyString && len(spec.Config.Binds) > 0 {
		if err := r.VolumeSetConfigData(ctx, name, spec.Config.Name); err != nil {
			log.Errorf("%s can not set config data to volume: %s", logPodPrefix, err.Error())
			return st, err
		}
	}

	r.state.Volumes().SetLocal(name)

	return st, nil
}

func (r Runtime) PodVolumeDestroy(ctx context.Context, pod, volume string) error {
	r.state.Volumes().DelLocal(r.podVolumeKeyCreate(pod, volume))
	return r.VolumeDestroy(ctx, r.podVolumeKeyCreate(pod, volume))
}

func (r Runtime) podVolumeClaimRestore(key string, manifest *models.PodManifest) {

	pod := r.state.Pods().GetPod(key)
	if pod == nil {
		return
	}

	for _, v := range manifest.Template.Volumes {

		var name string
		if v.Volume.Name != models.EmptyString {
			name = fmt.Sprintf("%s:%s", r.getPodNamespace(key), v.Volume.Name)
		} else {
			name = r.podVolumeKeyCreate(key, v.Name)
		}

		vol := r.state.Volumes().GetVolume(name)
		if vol == nil {
			continue
		}

		claim := r.state.Volumes().GetClaim(r.podVolumeClaimNameCreate(key, v.Name))
		if claim == nil {
			pv := &models.VolumeClaim{
				Name:   r.podVolumeClaimNameCreate(key, v.Name),
				Volume: name,
				Path:   vol.Status.Path,
			}

			r.state.Volumes().SetClaim(pv.Name, pv)
			pod.Volumes[pv.Name] = pv
		} else {
			pod.Volumes[claim.Name] = claim
		}
	}
}

func (r Runtime) podVolumeKeyCreate(pod, volume string) string {
	return fmt.Sprintf("%s-%s", strings.Replace(pod, ":", "-", -1), volume)
}

func (r Runtime) podVolumeClaimNameCreate(pod, volume string) string {
	return fmt.Sprintf("%s:%s", pod, volume)
}

func (r Runtime) podLocalFileCreate(path, data string) error {
	return filesystem.WriteStrToFile(path, data, 0777)
}

func (r Runtime) podLocalFileDestroy(path string) error {
	return os.RemoveAll(path)
}

// TODO: move to distribution
func (r Runtime) getPodNamespace(key string) string {
	var namespace = models.DEFAULT_NAMESPACE

	parts := strings.Split(key, ":")

	if len(parts) == 4 {
		namespace = parts[0]
	}

	return namespace
}
