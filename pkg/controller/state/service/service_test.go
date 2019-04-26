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

package service

import (
	"context"
	"github.com/spf13/viper"
	"testing"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/ipam"
	"github.com/lastbackend/lastbackend/pkg/controller/state/cluster"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"github.com/stretchr/testify/assert"
	"time"
)

func init() {
	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ipm, _ := ipam.New("")
	envs.Get().SetIPAM(ipm)
}

func testServiceObserver(t *testing.T, name, werr string, wst *ServiceState, state *ServiceState, svc *types.Service) {

	var (
		ctx = context.Background()
		err error
	)

	stg := envs.Get().GetStorage()

	err = stg.Del(ctx, stg.Collection().Deployment(), "")
	if !assert.NoError(t, err) {
		return
	}

	err = stg.Del(ctx, stg.Collection().Pod(), "")
	if !assert.NoError(t, err) {
		return
	}

	t.Run(name, func(t *testing.T) {

		err := serviceObserve(state, svc)
		if werr != types.EmptyString {

			if assert.NoError(t, err, "error should be presented") {
				return
			}

			if !assert.Equal(t, werr, err.Error(), "err message different") {
				return
			}

			return
		}

		if wst.service == nil {
			if !assert.Nil(t, state.service, "service should be nil") {
				return
			}

		}

		if wst.service != nil {

			// check service status state is equal
			if !assert.Equal(t, wst.service.Status.State, state.service.Status.State,
				"status state is different") {
				return
			}

			// check service status message is equal
			if !assert.Equal(t, wst.service.Status.Message, state.service.Status.Message,
				"status message is different") {
				return
			}

		}

		// check endpoint
		if wst.endpoint.endpoint != nil {
			if !assert.NotNil(t, state.endpoint.endpoint, "endpoint should be not nil") {
				return
			}
			if !assert.Equal(t, wst.endpoint.endpoint.Meta.Name, state.endpoint.endpoint.Meta.Name,
				"endpoint is different") {
				return
			}

			if !assert.Equal(t, wst.endpoint.endpoint.Spec.PortMap, state.endpoint.endpoint.Spec.PortMap,
				"endpoint portmap is different") {
				return
			}

		}

		if wst.endpoint.endpoint == nil {
			if !assert.Nil(t, state.endpoint.endpoint, "endpoint should be nil") {
				return
			}
		}

		// check provision deployment
		if wst.deployment.provision != nil {
			if !assert.NotNil(t, state.deployment.provision, "provision deployment should be not nil") {
				return
			}

			if !assert.Equal(t,
				wst.deployment.provision.Spec.Template.Updated,
				state.deployment.provision.Spec.Template.Updated,
				"provision deployment is different") {
				return
			}
		}

		if wst.deployment.provision == nil {
			if !assert.Nil(t, state.deployment.provision, "provision deployment should be nil") {
				return
			}
		}

		if wst.deployment.active != nil {

			if !assert.NotNil(t, wst.deployment.active, "active deployment should be not nil") {
				return
			}
			// check active deployment
			if !assert.Equal(t,
				wst.deployment.active.Spec.Template.Updated,
				state.deployment.active.Spec.Template.Updated,
				"provision deployment is different") {
				return
			}
		}

		if wst.deployment.active == nil {
			if !assert.Nil(t, state.deployment.active, "active deployment should be nil") {
				return
			}
		}

		// check deployments count
		if !assert.Equal(t,
			len(wst.deployment.list),
			len(state.deployment.list),
			"deployment count is different") {
			return
		}

		for _, d := range wst.deployment.list {

			if _, ok := state.deployment.list[d.SelfLink().String()]; ok {

				if !assert.Equal(t,
					d.Spec.Replicas,
					state.deployment.list[d.SelfLink().String()].Spec.Replicas,
					"deployment replicas not match") {
					return
				}

				if !assert.Equal(t,
					d.Status.State,
					state.deployment.list[d.SelfLink().String()].Status.State,
					"deployment status state not match") {
					return
				}

				if !assert.Equal(t,
					d.Status.Message,
					state.deployment.list[d.SelfLink().String()].Status.Message,
					"deployment status message not match") {
					return
				}
			}
		}

	})
}

func testStatusState(t *testing.T, fn func(*ServiceState) error, name string, wst, state *ServiceState) {

	t.Run(name, func(t *testing.T) {

		fn(state)

		if wst.service == nil {
			if !assert.Nil(t, state.service, "service should be nil") {
				return
			}

		}

		if wst.service != nil {

			// check service status state is equal
			if !assert.Equal(t, wst.service.Status.State, state.service.Status.State,
				"status state is different") {
				return
			}

			// check service status message is equal
			if !assert.Equal(t, wst.service.Status.Message, state.service.Status.Message,
				"status message is different") {
				return
			}

		}

		// check endpoint
		if wst.endpoint.endpoint != nil {
			if !assert.NotNil(t, state.endpoint, "endpoint should be not nil") {
				return
			}
			if !assert.Equal(t, wst.endpoint.endpoint.Meta.Name, state.endpoint.endpoint.Meta.Name,
				"endpoint is different") {
				return
			}

			if !assert.Equal(t, wst.endpoint.endpoint.Spec.PortMap, state.endpoint.endpoint.Spec.PortMap,
				"endpoint portmap is different") {
				return
			}

		}

		if wst.endpoint.endpoint == nil {
			if !assert.Nil(t, state.endpoint.endpoint, "endpoint should be nil") {
				return
			}
		}

		// check provision deployment
		if wst.deployment.provision != nil {
			if !assert.NotNil(t, state.deployment.provision, "provision deployment should be not nil") {
				return
			}

			if !assert.Equal(t,
				wst.deployment.provision.Spec.Template.Updated,
				state.deployment.provision.Spec.Template.Updated,
				"provision deployment is different") {
				return
			}
		}

		if wst.deployment.provision == nil {
			if !assert.Nil(t, state.deployment.provision, "provision deployment should be nil") {
				return
			}
		}

		if wst.deployment.active != nil {

			if !assert.NotNil(t, wst.deployment.active, "active deployment should be not nil") {
				return
			}
			// check active deployment
			if !assert.Equal(t,
				wst.deployment.active.Spec.Template.Updated,
				state.deployment.active.Spec.Template.Updated,
				"provision deployment is different") {
				return
			}
		}

		if wst.deployment.active == nil {
			if !assert.Nil(t, state.deployment.active, "active deployment should be nil") {
				return
			}
		}

		// check deployments count
		if !assert.Equal(t,
			len(wst.deployment.list),
			len(state.deployment.list),
			"deployment count is different") {
			return
		}

		for _, d := range wst.deployment.list {

			if _, ok := state.deployment.list[d.SelfLink().String()]; ok {

				if !assert.Equal(t,
					d.Spec.Replicas,
					state.deployment.list[d.SelfLink().String()].Spec.Replicas,
					"deployment replicas not match") {
					return
				}

				if !assert.Equal(t,
					d.Status.State,
					state.deployment.list[d.SelfLink().String()].Status.State,
					"deployment status state not match") {
					return
				}

				if !assert.Equal(t,
					d.Status.Message,
					state.deployment.list[d.SelfLink().String()].Status.Message,
					"deployment status message not match") {
					return
				}
			}
		}

	})
}

func TestHandleServiceStateCreated(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			svc   *types.Service
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle without endpoint create"}

		s.args.svc = getServiceAsset(types.StateCreated, types.EmptyString)
		s.args.state = getServiceStateAsset(s.args.svc)
		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		dp := getDeploymentAsset(s.want.state.service, types.StateCreated, types.EmptyString)
		s.want.state.deployment.provision = dp
		s.want.state.deployment.list[dp.SelfLink().String()] = dp

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with endpoint create"}

		s.args.svc = getServiceAsset(types.StateCreated, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)
		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		dp := getDeploymentAsset(s.want.state.service, types.StateCreated, types.EmptyString)

		s.want.state.deployment.provision = dp
		s.want.state.deployment.list[dp.SelfLink().String()] = dp
		s.want.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with endpoint update"}

		s.args.svc = getServiceAsset(types.StateCreated, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.endpoint.endpoint = getEndpointAsset(s.args.svc)
		s.args.state.endpoint.endpoint.Spec.PortMap[9000] = "9000/udp"

		s.want.err = types.EmptyString
		s.want.state = getServiceStateAsset(s.args.svc)
		dp := getDeploymentAsset(s.want.state.service, types.StateCreated, types.EmptyString)

		s.want.state.deployment.provision = dp
		s.want.state.deployment.list[dp.SelfLink().String()] = dp

		s.args.svc.Spec.Network.Ports[8080] = "8080/tcp"
		s.want.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with endpoint remove"}

		s.args.svc = getServiceAsset(types.StateCreated, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		s.args.svc.Spec.Network.Ports = make(map[uint16]string)

		s.want.err = types.EmptyString
		s.want.state = getServiceStateAsset(s.args.svc)
		dp := getDeploymentAsset(s.want.state.service, types.StateCreated, types.EmptyString)

		s.want.state.deployment.provision = dp
		s.want.state.deployment.list[dp.SelfLink().String()] = dp

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with deployment provision update"}

		s.args.svc = getServiceAsset(types.StateCreated, types.EmptyString)
		s.args.svc.Spec.Template.Updated.Add(-5 * time.Second)

		dp1 := getDeploymentAsset(s.args.svc, types.StateCreated, types.EmptyString)

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.svc.Spec.Template.Updated = time.Now()

		s.args.state.deployment.provision = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		dp2 := getDeploymentAsset(s.want.state.service, types.StateProvision, types.EmptyString)

		s.want.state.service.Status.State = types.StateCreated
		s.want.state.deployment.provision = dp2
		s.want.state.deployment.list[dp2.SelfLink().String()] = dp2
		s.want.state.deployment.list[dp1.SelfLink().String()].Status.State = types.StateDestroyed

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with deployment active scale"}

		s.args.svc = getServiceAsset(types.StateCreated, types.EmptyString)
		dp1 := getDeploymentAsset(s.args.svc, types.StateReady, types.EmptyString)
		cdp := *dp1

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1

		s.want.err = types.EmptyString
		s.args.svc.Spec.Replicas++

		cdp.Spec.Replicas = s.args.svc.Spec.Replicas
		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.deployment.list[cdp.SelfLink().String()] = &cdp
		s.want.state.deployment.active = &cdp
		s.want.state.deployment.list[cdp.SelfLink().String()].Status.State = types.StateProvision

		return s
	}())

	for _, tt := range tests {
		testServiceObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.svc)
	}

}

func TestHandleServiceStateProvision(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			svc   *types.Service
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with endpoint create"}

		s.args.svc = getServiceAsset(types.StateProvision, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)
		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		dp := getDeploymentAsset(s.want.state.service, types.StateCreated, types.EmptyString)

		s.want.state.deployment.provision = dp
		s.want.state.deployment.list[dp.SelfLink().String()] = dp
		s.want.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with endpoint update"}

		s.args.svc = getServiceAsset(types.StateProvision, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.endpoint.endpoint = getEndpointAsset(s.args.svc)
		s.args.state.endpoint.endpoint.Spec.PortMap[9000] = "9000/udp"

		s.want.err = types.EmptyString
		s.want.state = getServiceStateAsset(s.args.svc)
		dp := getDeploymentAsset(s.want.state.service, types.StateCreated, types.EmptyString)

		s.want.state.deployment.provision = dp
		s.want.state.deployment.list[dp.SelfLink().String()] = dp

		s.args.svc.Spec.Network.Ports[80] = "8000/tcp"
		s.args.svc.Spec.Network.Ports[8080] = "8080/tcp"
		s.want.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with endpoint remove"}

		s.args.svc = getServiceAsset(types.StateProvision, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		s.args.svc.Spec.Network.Ports = make(map[uint16]string)

		s.want.err = types.EmptyString
		s.want.state = getServiceStateAsset(s.args.svc)
		dp := getDeploymentAsset(s.want.state.service, types.StateCreated, types.EmptyString)

		s.want.state.deployment.provision = dp
		s.want.state.deployment.list[dp.SelfLink().String()] = dp

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with active deployment scale"}

		s.args.svc = getServiceAsset(types.StateProvision, types.EmptyString)
		dp1 := getDeploymentAsset(s.args.svc, types.StateReady, types.EmptyString)
		cdp := *dp1

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1

		s.want.err = types.EmptyString
		s.args.svc.Spec.Replicas++

		cdp.Spec.Replicas = s.args.svc.Spec.Replicas
		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.deployment.list[cdp.SelfLink().String()] = &cdp
		s.want.state.deployment.active = &cdp
		s.want.state.deployment.list[cdp.SelfLink().String()].Status.State = types.StateProvision

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with provision deployment scale"}

		s.args.svc = getServiceAsset(types.StateProvision, types.EmptyString)
		dp1 := getDeploymentAsset(s.args.svc, types.StateError, types.EmptyString)
		dp2 := getDeploymentAsset(s.args.svc, types.StateCreated, types.EmptyString)

		cdp1 := *dp1
		cdp2 := *dp2

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.want.err = types.EmptyString
		s.args.svc.Spec.Replicas++

		cdp2.Spec.Replicas = s.args.svc.Spec.Replicas
		cdp2.Status.State = types.StateProvision

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.deployment.list[cdp1.SelfLink().String()] = &cdp1
		s.want.state.deployment.list[cdp2.SelfLink().String()] = &cdp2
		s.want.state.deployment.active = &cdp1
		s.want.state.deployment.provision = &cdp2

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with new deployment create with provision deployment"}

		s.args.svc = getServiceAsset(types.StateProvision, types.EmptyString)
		s.args.svc.Spec.Template.Updated.Add(-5 * time.Second)
		dp1 := getDeploymentAsset(s.args.svc, types.StateReady, types.EmptyString)
		dp2 := getDeploymentAsset(s.args.svc, types.StateError, types.EmptyString)

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.args.svc.Spec.Template.Updated = time.Now()
		dp3 := getDeploymentAsset(s.args.svc, types.StateProvision, types.EmptyString)

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.list[dp1.SelfLink().String()].Status.State = types.StateReady
		s.want.state.deployment.list[dp2.SelfLink().String()].Status.State = types.StateDestroyed
		s.want.state.deployment.list[dp3.SelfLink().String()] = dp3
		s.want.state.deployment.provision = dp3

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with new deployment create without provision deployment"}

		s.args.svc = getServiceAsset(types.StateProvision, types.EmptyString)
		s.args.svc.Spec.Template.Updated.Add(-5 * time.Second)
		dp1 := getDeploymentAsset(s.args.svc, types.StateError, types.EmptyString)

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1

		s.args.svc.Spec.Template.Updated = time.Now()
		dp2 := getDeploymentAsset(s.args.svc, types.StateProvision, types.EmptyString)

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.list[dp1.SelfLink().String()].Status.State = types.StateDestroyed
		s.want.state.deployment.list[dp2.SelfLink().String()] = dp2
		s.want.state.deployment.provision = dp2

		return s
	}())

	for _, tt := range tests {
		testServiceObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.svc)
	}
}

func TestHandleServiceStateReady(t *testing.T) {
	type suit struct {
		name string
		args struct {
			state *ServiceState
			svc   *types.Service
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {
		s := suit{name: "check ready state stay ready"}

		svc := getServiceAsset(types.StateReady, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.svc = svc

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateReady

		return s
	}())

	for _, tt := range tests {
		testServiceObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.svc)
	}
}

func TestHandleServiceStateError(t *testing.T) {
	type suit struct {
		name string
		args struct {
			state *ServiceState
			svc   *types.Service
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {
		s := suit{name: "check error state stay error"}

		svc := getServiceAsset(types.StateError, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.svc = svc

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateError

		return s
	}())

	for _, tt := range tests {
		testServiceObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.svc)
	}
}

func TestHandleServiceStateDegradation(t *testing.T) {
	type suit struct {
		name string
		args struct {
			state *ServiceState
			svc   *types.Service
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {
		s := suit{name: "check degradation state stay degradation"}

		svc := getServiceAsset(types.StateDegradation, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.svc = svc

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateDegradation

		return s
	}())

	for _, tt := range tests {
		testServiceObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.svc)
	}
}

func TestHandleServiceStateDestroy(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			svc   *types.Service
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle without deployments and endpoint"}

		s.args.svc = getServiceAsset(types.StateDestroy, types.EmptyString)
		s.args.state = getServiceStateAsset(s.args.svc)
		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.service = nil

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle without deployments"}

		s.args.svc = getServiceAsset(types.StateDestroy, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.service = nil
		s.want.state.endpoint.endpoint = nil
		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with deployments and without endpoint"}

		s.args.svc = getServiceAsset(types.StateDestroy, types.EmptyString)
		s.args.state = getServiceStateAsset(s.args.svc)

		dp1 := getDeploymentAsset(s.args.svc, types.StateReady, types.EmptyString)
		dp2 := getDeploymentAsset(s.args.svc, types.StateCreated, types.EmptyString)

		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		cdp1 := *dp1
		cdp2 := *dp2

		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.deployment.list[cdp1.SelfLink().String()] = &cdp1
		s.want.state.deployment.list[cdp2.SelfLink().String()] = &cdp2

		s.want.state.deployment.active = &cdp1
		s.want.state.deployment.provision = &cdp2

		cdp1.Status.State = types.StateDestroy
		cdp2.Status.State = types.StateDestroy

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with deployments and with endpoint"}

		s.args.svc = getServiceAsset(types.StateDestroy, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)

		dp1 := getDeploymentAsset(s.args.svc, types.StateReady, types.EmptyString)
		dp2 := getDeploymentAsset(s.args.svc, types.StateCreated, types.EmptyString)

		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.args.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		cdp1 := *dp1
		cdp2 := *dp2

		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.deployment.list[cdp1.SelfLink().String()] = &cdp1
		s.want.state.deployment.list[cdp2.SelfLink().String()] = &cdp2

		s.want.state.deployment.active = &cdp1
		s.want.state.deployment.provision = &cdp2

		cdp1.Status.State = types.StateDestroy
		cdp2.Status.State = types.StateDestroy

		s.want.state.endpoint.endpoint = nil

		return s
	}())

	for _, tt := range tests {
		testServiceObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.svc)
	}

}

func TestHandleServiceStateDestroyed(t *testing.T) {
	type suit struct {
		name string
		args struct {
			state *ServiceState
			svc   *types.Service
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle without deployments and endpoint"}

		s.args.svc = getServiceAsset(types.StateDestroyed, types.EmptyString)
		s.args.state = getServiceStateAsset(s.args.svc)
		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.service = nil

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle without deployments"}

		s.args.svc = getServiceAsset(types.StateDestroyed, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)
		s.args.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.service = nil
		s.want.state.endpoint.endpoint = nil

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with deployments and without endpoint"}

		s.args.svc = getServiceAsset(types.StateDestroyed, types.EmptyString)
		s.args.state = getServiceStateAsset(s.args.svc)

		dp1 := getDeploymentAsset(s.args.svc, types.StateReady, types.EmptyString)
		dp2 := getDeploymentAsset(s.args.svc, types.StateCreated, types.EmptyString)

		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		cdp1 := *dp1
		cdp2 := *dp2

		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.deployment.list[cdp1.SelfLink().String()] = &cdp1
		s.want.state.deployment.list[cdp2.SelfLink().String()] = &cdp2

		s.want.state.deployment.active = &cdp1
		s.want.state.deployment.provision = &cdp2

		cdp1.Status.State = types.StateDestroy
		cdp2.Status.State = types.StateDestroy

		s.want.state.service.Status.State = types.StateDestroy

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "successful state handle with deployments and with endpoint"}

		s.args.svc = getServiceAsset(types.StateDestroyed, types.EmptyString)
		s.args.svc.Spec.Network.Ports = make(map[uint16]string)
		s.args.svc.Spec.Network.Ports[80] = "80/tcp"

		s.args.state = getServiceStateAsset(s.args.svc)

		dp1 := getDeploymentAsset(s.args.svc, types.StateReady, types.EmptyString)
		dp2 := getDeploymentAsset(s.args.svc, types.StateCreated, types.EmptyString)

		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.args.state.endpoint.endpoint = getEndpointAsset(s.args.svc)

		cdp1 := *dp1
		cdp2 := *dp2

		s.want.err = types.EmptyString

		s.want.state = getServiceStateAsset(s.args.svc)
		s.want.state.deployment.list[cdp1.SelfLink().String()] = &cdp1
		s.want.state.deployment.list[cdp2.SelfLink().String()] = &cdp2

		s.want.state.deployment.active = &cdp1
		s.want.state.deployment.provision = &cdp2

		cdp1.Status.State = types.StateDestroy
		cdp2.Status.State = types.StateDestroy

		s.want.state.endpoint.endpoint = nil
		s.want.state.service.Status.State = types.StateDestroy

		return s
	}())

	for _, tt := range tests {
		testServiceObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.svc)
	}
}

func TestServiceStatusState(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
		}
		want struct {
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {
		s := suit{name: "check created state stay created"}

		svc := getServiceAsset(types.StateCreated, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateCreated

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "check created state to provision without active deployment"}

		svc := getServiceAsset(types.StateCreated, types.EmptyString)
		dp1 := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		dp2 := getDeploymentAsset(svc, types.StateError, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateProvision

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "check created state to provision with active provision deployment"}

		svc := getServiceAsset(types.StateCreated, types.EmptyString)
		dp1 := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		dp2 := getDeploymentAsset(svc, types.StateCancel, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateProvision

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "check created state to provision with active created deployment"}

		svc := getServiceAsset(types.StateCreated, types.EmptyString)
		dp1 := getDeploymentAsset(svc, types.StateCreated, types.EmptyString)
		dp2 := getDeploymentAsset(svc, types.StateCancel, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateCreated

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "check provision state to ready with active deployment"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)

		dp1 := getDeploymentAsset(svc, types.StateReady, types.EmptyString)
		dp2 := getDeploymentAsset(svc, types.StateError, types.EmptyString)
		dp2.Spec.Template.Containers[0].Name = "changed"
		dp2.Spec.Template.Updated.Add(3 * time.Second)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateReady

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "check provision state to ready with active deployment and provision error"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)

		dp1 := getDeploymentAsset(svc, types.StateReady, types.EmptyString)
		dp2 := getDeploymentAsset(svc, types.StateError, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateReady

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "check error state to ready with active deployment"}

		svc := getServiceAsset(types.StateError, types.EmptyString)

		dp1 := getDeploymentAsset(svc, types.StateReady, types.EmptyString)
		dp2 := getDeploymentAsset(svc, types.StateError, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateReady

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "check provision state to error with active deployment"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)

		dp1 := getDeploymentAsset(svc, types.StateError, types.EmptyString)
		dp2 := getDeploymentAsset(svc, types.StateError, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateError

		return s
	}())

	tests = append(tests, func() suit {
		s := suit{name: "check provision state to error without active deployment"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)

		dp1 := getDeploymentAsset(svc, types.StateError, types.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1

		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = types.StateError

		return s
	}())

	for _, tt := range tests {
		testStatusState(t, serviceStatusState, tt.name, tt.want.state, tt.args.state)
	}
}

func getServiceAsset(state, message string) *types.Service {
	s := new(types.Service)

	s.Meta.Namespace = "test"
	s.Meta.Name = "service"
	s.Meta.SelfLink = *types.NewServiceSelfLink(s.Meta.Namespace, s.Meta.Name)

	s.Status.State = state
	s.Status.Message = message

	s.Spec.Replicas = 1
	s.Spec.Template.Updated = time.Now()

	return s
}

func getEndpointAsset(svc *types.Service) *types.Endpoint {
	e := new(types.Endpoint)

	e.Meta.Namespace = svc.Meta.Namespace
	e.Meta.Name = svc.Meta.Name
	e.Meta.SelfLink = *types.NewEndpointSelfLink(e.Meta.Namespace, e.Meta.Name)

	e.Spec.PortMap = make(map[uint16]string)
	for k, v := range svc.Spec.Network.Ports {
		e.Spec.PortMap[k] = v
	}

	return e
}

func getDeploymentAsset(svc *types.Service, state, message string) *types.Deployment {

	d := new(types.Deployment)

	d.Meta.Namespace = svc.Meta.Namespace
	d.Meta.Service = svc.Meta.Name
	d.Meta.Name = generator.GetUUIDV4()
	d.Meta.SelfLink = *types.NewDeploymentSelfLink(d.Meta.Namespace, d.Meta.Service, d.Meta.Name)

	d.Status.State = state
	d.Status.Message = message
	d.Status.Dependencies.Volumes = make(map[string]types.StatusDependency, 0)
	d.Status.Dependencies.Secrets = make(map[string]types.StatusDependency, 0)
	d.Status.Dependencies.Configs = make(map[string]types.StatusDependency, 0)

	d.Spec.State = svc.Spec.State
	d.Spec.Template = svc.Spec.Template
	d.Spec.Replicas = svc.Spec.Replicas

	d.Spec.Template.Containers = types.SpecTemplateContainers{}
	d.Spec.Template.Containers = append(d.Spec.Template.Containers, &types.SpecTemplateContainer{
		Name: "demo",
	})

	return d
}

func getPodAsset(d *types.Deployment, state, message string) *types.Pod {

	p := new(types.Pod)

	p.Meta.Namespace = d.Meta.Namespace
	p.Meta.Name = generator.GetUUIDV4()

	sl, _ := types.NewPodSelfLink(types.KindDeployment, d.SelfLink().String(), p.Meta.Name)
	p.Meta.SelfLink = *sl

	p.Status.State = state
	p.Status.Message = message

	if state == types.StateReady {
		p.Status.Running = true
	}

	p.Spec.State = d.Spec.State
	p.Spec.Template = d.Spec.Template

	return p
}

func getServiceStateAsset(svc *types.Service) *ServiceState {

	n := new(types.Node)

	n.Meta.Name = "node"
	n.Meta.Hostname = "node.local"
	n.Status.Capacity = types.NodeResources{
		Containers: 10,
		Pods:       10,
		RAM:        1000,
		CPU:        1,
		Storage:    1000,
	}
	n.Meta.SelfLink = *types.NewNodeSelfLink(n.Meta.Hostname)

	cs := cluster.NewClusterState()
	cs.SetNode(n)
	s := NewServiceState(cs, svc)
	return s
}

func getServiceStateCopy(ss *ServiceState) *ServiceState {

	svc := *ss.service

	s := NewServiceState(ss.cluster, &svc)

	for k, d := range ss.deployment.list {
		cd := *d
		s.deployment.list[k] = &cd
	}

	if ss.deployment.active != nil {
		s.deployment.active = s.deployment.list[ss.deployment.active.SelfLink().String()]
	}

	if ss.deployment.provision != nil {
		s.deployment.provision = s.deployment.list[ss.deployment.provision.SelfLink().String()]
	}

	if len(ss.pod.list) > 0 {
		for k, pl := range ss.pod.list {

			if _, ok := s.pod.list[k]; !ok {
				s.pod.list[k] = make(map[string]*types.Pod)
			}

			for l, p := range pl {
				s.pod.list[k][l] = p
			}

		}
	}

	return s
}
