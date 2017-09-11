//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	ctx "github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/app"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	h "github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 3

type pod struct {
	Context   context.Context
	Namespace types.Meta
}

func (p *pod) Set(pod types.Pod) error {
	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Pod: set pod %#v", pod)

	svc, err := storage.Service().GetByPodName(p.Context, pod.Meta.Name)
	if err != nil {
		if err.Error() == store.ErrKeyNotFound {
			log.V(logLevel).Warnf("Pod: not found pod %s", pod.Meta.Name)
			return nil
		}
		log.V(logLevel).Errorf("Pod: get service by pod `%s` err: %s", pod.Meta.Name, err.Error())
		return err
	}

	ns := app.New(p.Context)
	item, err := ns.Get(svc.Meta.App)
	if err != nil {
		log.V(logLevel).Errorf("Pod: get app `%s` err: %s", pod.Meta.Name, err.Error())
		return err
	}
	if item == nil {
		log.V(logLevel).Warnf("Pod: app `%s` not found", pod.Meta.Name, err.Error())
		return nil
	}

	log.V(logLevel).Debugf("Pod: find pod `%s` in service %s", pod.Meta.Name, svc.Meta.App)

	// If service has not this pod then skip it
	if _, ok := svc.Pods[pod.Meta.Name]; ok {

		log.V(logLevel).Debugf("Pod: update pod %s in service %s", pod.Meta.Name, svc.Meta.App)

		svc.Pods[pod.Meta.Name].Containers = pod.Containers
		svc.Pods[pod.Meta.Name].Meta = pod.Meta
		svc.Pods[pod.Meta.Name].State = pod.State
		svc.Pods[pod.Meta.Name].Node = pod.Node
		pd := svc.Pods[pod.Meta.Name]

		log.V(logLevel).Debugf("Pod: pod `%s` has state `%s`", pd.Meta.Name, pd.State.State)

		if pd.State.State == types.StateDestroyed {
			log.V(logLevel).Debugf("Pod: remove pod `%s` with `%s` state", pd.Meta.Name, types.StateDestroyed)
			if err := storage.Pod().Remove(p.Context, item.Meta.Name, pd); err != nil {
				log.V(logLevel).Errorf("Pod: remove pod `%s` with `%s` state err: %s", pd.Meta.Name, types.StateDestroyed, err.Error())
				return err
			}
			delete(svc.Pods, pd.Meta.Name)

		} else {

			log.V(logLevel).Debugf("Pod: update pod `%s` -> %#v", pd.Meta.Name, pd)

			if err := storage.Pod().Update(p.Context, item.Meta.Name, pd); err != nil {
				log.V(logLevel).Errorf("Pod: update pod err: %s", err.Error())
				return err
			}

			log.V(logLevel).Debugf("Pod: update service `%s` -> %#v", svc.Meta.Name, svc)

			// Need update data info (state and resources) for this service after update pod info
			if err := storage.Service().Update(p.Context, svc); err != nil {
				log.V(logLevel).Errorf("Pod: update service `%s` err: %s", svc.Meta.Name, err)
				return err
			}
		}
	} else {
		log.V(logLevel).Warnf("Pod: pod %s in service %s not found", pod.Meta.Name, svc.Meta.App)
	}

	log.V(logLevel).Debugf("Pod: сheck the possibility of removing the service %s", svc.Meta.Name)

	// Remove service if the state set as destroyed and pods count is zero
	if len(svc.Pods) == 0 && svc.State.State == types.StateDestroyed {

		log.V(logLevel).Debugf("Pod: remove service %s", svc.Meta.Name)

		if err = storage.Hook().Remove(p.Context, svc.Meta.Hook); err != nil && err.Error() != store.ErrKeyNotFound {
			log.V(logLevel).Errorf("Pod: remove service hook err: %s", err.Error())
			return err
		}
		if err = storage.Service().Remove(p.Context, svc); err != nil {
			if err.Error() == store.ErrKeyNotFound {
				log.V(logLevel).Warnf("Pod: service `%s` not found", svc.Meta.Name)
				return nil
			}
			log.V(logLevel).Debugf("Pod: remove service `%s` err: %s", svc.Meta.Name, err.Error())
			return err
		}
	}

	return nil
}

func Logs(c context.Context, ns, pod, container string, stream io.Writer, done chan bool) error {

	const buffer_size = 1024

	var (
		storage  = ctx.Get().GetStorage()
		buffer   = make([]byte, buffer_size)
		doneChan = make(chan bool, 1)
	)

	log.V(logLevel).Debugf("Pod: get container `%s` logs for pod `%s` in app `%s`", container, pod, ns)

	svc, err := storage.Service().GetByPodName(c, pod)
	if err != nil {
		if err.Error() == store.ErrKeyNotFound {
			log.V(logLevel).Debugf("Pod: pod `%s` with container `%s` not found", pod, c)
			return nil
		}
		log.V(logLevel).Errorf("Pod: get container `%s` logs for pod `%s` in app `%s` err: %s", container, pod, ns, err.Error())
		return err
	}

	_ns := app.New(c)
	item, err := _ns.Get(svc.Meta.App)
	if err != nil {
		log.V(logLevel).Errorf("Pod: get app `%s` err: %s", svc.Meta.App, err.Error())
		return err
	}
	if item == nil {
		log.V(logLevel).Warnf("Pod: app `%s` not found", svc.Meta.App, err.Error())
		return err
	}

	if ns != item.Meta.Name {
		log.V(logLevel).Errorf("Pod: app `%s` not found", ns)
		return err
	}

	p, e := storage.Pod().GetByName(c, item.Meta.Name, pod)
	if e != nil {
		log.V(logLevel).Errorf("Pod: get pod `%s` err: %s", item.Meta.Name, err.Error())
		return err
	}

	var cnt string
	for c := range p.Containers {
		if c == container {
			cnt = container
			break
		}
	}
	if cnt == "" {
		log.V(logLevel).Errorf("Pod: container `%s` not found", container)
		return errors.New("access denied")
	}

	n, err := storage.Node().Get(c, p.Node.ID)
	if err != nil {
		log.V(logLevel).Errorf("Pod: get node by id `%s` err: %s", p.Node.ID, err.Error())
		return err
	}

	uri := fmt.Sprintf("%s:%d", n.Meta.IP, n.Meta.Port)
	client, err := h.New(uri, &h.ReqOpts{TLS: false})
	if err != nil {
		log.V(logLevel).Errorf("Pod: create http client err: %s", err.Error())
		return err
	}

	_, res, err := client.
		GET(fmt.Sprintf("/container/%s/logs", cnt)).Do()
	if err != nil {
		log.V(logLevel).Errorf("Pod: http request err: %s", err.Error())
		return err
	}

	go func() {
		for {
			select {
			case <-doneChan:
				res.Body.Close()
				return
			default:
				n, err := res.Body.Read(buffer)
				if err != nil {
					log.V(logLevel).Errorf("Pod: read bytes from stream err: %s", err)
					res.Body.Close()
					return
				}

				_, err = func(p []byte) (n int, err error) {
					n, err = stream.Write(p)
					if err != nil {
						log.V(logLevel).Errorf("Pod: write bytes from stream err: %s", err)
						return n, err
					}
					if f, ok := stream.(http.Flusher); ok {
						f.Flush()
					}
					return n, nil
				}(buffer[0:n])
				if err != nil {
					log.V(logLevel).Errorf("Pod: written to stream err: %s", err)
					return
				}

				for i := 0; i < n; i++ {
					buffer[i] = 0
				}
			}
		}
	}()

	<-done

	close(doneChan)

	return nil
}

func New(ctx context.Context) *pod {
	return &pod{
		Context: ctx,
	}
}
