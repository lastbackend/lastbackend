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

package events

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
)

const (
	logLevel  = 3
	logPrefix = "node:events"
)

// NewConnectEventt - send node info event after
// node is successful accepted and each hour
func NewConnectEvent(ctx context.Context) error {

	var (
		c = envs.Get().GetNodeClient()
	)

	opts := v1.Request().Node().NodeConnectOptions()
	opts.Info = envs.Get().GetState().Node().Info
	opts.Status = envs.Get().GetState().Node().Status
	opts.Network = *envs.Get().GetCNI().Info(ctx)

	if viper.IsSet("node.tls") {
		opts.TLS = !viper.GetBool("node.tls.insecure")

		if opts.TLS {
			caData, err := ioutil.ReadFile(viper.GetString("node.tls.ca"))
			if err != nil {
				log.Errorf("%s:connect_event:> read ca cert file err: %v", logPrefix, err)
				return err
			}

			certData, err := ioutil.ReadFile(viper.GetString("node.tls.client_cert"))
			if err != nil {
				log.Errorf("%s:connect_event:> read client cert file err: %v", logPrefix, err)
				return err
			}

			keyData, err := ioutil.ReadFile(viper.GetString("node.tls.client_key"))
			if err != nil {
				log.Errorf("%s:connect_event:> read client key file err: %v", logPrefix, err)
				return err
			}

			opts.SSL = new(request.SSL)
			opts.SSL.CA = caData
			opts.SSL.Key = keyData
			opts.SSL.Cert = certData
		}
	}

	return c.Connect(ctx, opts)

}

// NewStatusEvent - send node state event after
// node is successful accepted and each hour
func NewStatusEvent(ctx context.Context) error {
	var (
		e = envs.Get().GetExporter()
	)

	e.Resources(envs.Get().GetState().Node().Status)
	return nil
}

// NewPodStatusEvent - send pod state event after
// node is successful accepted and each hour
func NewPodStatusEvent(ctx context.Context, pod string) error {

	var (
		p = envs.Get().GetState().Pods().GetPod(pod)
		e = envs.Get().GetExporter()
	)

	if pod == "" {
		log.Errorf("%s:pod_status_event:> pod state event: pod is empty", logPrefix)
		return errors.New("pod state event: pod is empty")
	}

	if p == nil {
		return nil
	}

	e.PodStatus(pod, p)

	return nil
}

// NewRouteStatusEvent - send pod state event after
// node is successful accepted and each hour
func NewVolumeStatusEvent(ctx context.Context, volume string) error {

	var (
		c = envs.Get().GetNodeClient()
	)

	if volume == "" {
		log.Errorf("%s:volume_status_event:> volume state event: volume is empty", logPrefix)
		return errors.New("Event: volume state event: volume is empty")
	}

	log.V(logLevel).Debugf("%s:volume_status_event:> volume state event state: %s", logPrefix, volume)

	opts := v1.Request().Node().NodeVolumeStatusOptions()
	return c.SetVolumeStatus(ctx, volume, opts)
}
