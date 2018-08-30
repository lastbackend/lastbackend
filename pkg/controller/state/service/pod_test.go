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

package service

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/ipam"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testPodObserver(t *testing.T, name, werr string, wst *ServiceState, state *ServiceState, p *types.Pod) {

	var (
		ctx = context.Background()
		err error
	)

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ipm, _ := ipam.New("")
	envs.Get().SetIPAM(ipm)

	err = stg.Del(ctx, stg.Collection().Deployment(), "")
	if !assert.NoError(t, err) {
		return
	}

	err = stg.Del(ctx, stg.Collection().Endpoint(), "")
	if !assert.NoError(t, err) {
		return
	}

	t.Run(name, func(t *testing.T) {

		err := PodObserve(state, p)
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
				"provision deployment status state is different") {
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
				"active deployment is different") {
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

		// check deployment state
		for _, d := range wst.deployment.list {

			if _, ok := state.deployment.list[d.SelfLink()]; !ok {
				t.Errorf("deployment not found %s", d.SelfLink())
				return
			}

			if !assert.Equal(t,
				d.Spec.Replicas,
				state.deployment.list[d.SelfLink()].Spec.Replicas,
				"deployment replicas not match") {
				return
			}

			if !assert.Equal(t,
				d.Status.State,
				state.deployment.list[d.SelfLink()].Status.State,
				"deployment status state not match") {
				return
			}

			if !assert.Equal(t,
				d.Status.Message,
				state.deployment.list[d.SelfLink()].Status.Message,
				"deployment status message not match") {
				return
			}
		}

		// check pods count
		if !assert.Equal(t,
			len(wst.pod.list),
			len(state.pod.list),
			"pod deployment groups count is different") {
			return
		}

		// check pod states
		for d, pl := range wst.pod.list {

			if !assert.NotNil(t, state.pod.list[d], "pod group should exist") {
				return
			}

			for l, p := range pl {

				if !assert.NotNil(t, state.pod.list[d][l], "pod should exist") {
					return
				}

				switch state.pod.list[d][l].Status.State {
				case types.StateProvision:
					if !assert.NotEqual(t, state.pod.list[d][l].Meta.Node, types.EmptyString,
						"node value should not be empty") {
						return
					}
				case types.StateDestroyed:
					if !assert.Equal(t, state.pod.list[d][l].Meta.Node, types.EmptyString,
						"node value should be empty") {
						return
					}
				}

				if !assert.Equal(t,
					p.Status.State,
					state.pod.list[d][l].Status.State,
					"pod status state not match") {
					return
				}

				if !assert.Equal(t,
					p.Status.State,
					state.pod.list[d][l].Status.State,
					"pod status message not match") {
					return
				}

			}
		}

	})
}

func TestHandlePodStateCreated(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			pod   *types.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateCreated, types.EmptyString)
		pod := getPodAsset(dp, types.StateCreated, types.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[pod.DeploymentLink()][pod.SelfLink()].Status.State = types.StateProvision

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.pod)
	}
}

func TestHandlePodStateProvision(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			pod   *types.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateReady, types.EmptyString)
		pod := getPodAsset(dp, types.StateProvision, types.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateProvision

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, set to destroy"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		pod := getPodAsset(dp, types.StateProvision, types.EmptyString)
		pod.Meta.Node = "demo"
		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[pod.DeploymentLink()][pod.SelfLink()].Status.State = types.StateDestroy

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, set to destroyed"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		pod := getPodAsset(dp, types.StateProvision, types.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod, 0)

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.pod)
	}
}

func TestHandlePodStateReady(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			pod   *types.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with change deployment state"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		pod := getPodAsset(dp, types.StateReady, types.EmptyString)
		pod.Meta.Node = "node"

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateReady

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with change deployment state to ready"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		p1 := getPodAsset(dp, types.StateReady, types.EmptyString)
		p2 := getPodAsset(dp, types.StateReady, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateReady

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, stay deployment state"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		p1 := getPodAsset(dp, types.StateProvision, types.EmptyString)
		p1.Meta.Node = "node"

		p2 := getPodAsset(dp, types.StateReady, types.EmptyString)

		s.args.pod = p2

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateProvision

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.pod)
	}
}

func TestHandlePodStateError(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			pod   *types.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to error"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		pod := getPodAsset(dp, types.StateError, types.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateError

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to degradation"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		p1 := getPodAsset(dp, types.StateError, types.EmptyString)
		p2 := getPodAsset(dp, types.StateReady, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateDegradation

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to error with many pods"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		p1 := getPodAsset(dp, types.StateError, types.EmptyString)
		p2 := getPodAsset(dp, types.StateError, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateError

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.pod)
	}
}

func TestHandlePodStateDegradation(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			pod   *types.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to degradation"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		pod := getPodAsset(dp, types.StateDegradation, types.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateDegradation

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to degradation"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		p1 := getPodAsset(dp, types.StateDegradation, types.EmptyString)
		p2 := getPodAsset(dp, types.StateReady, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateDegradation

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to degradation"}

		svc := getServiceAsset(types.StateProvision, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateProvision, types.EmptyString)
		p1 := getPodAsset(dp, types.StateDegradation, types.EmptyString)
		p2 := getPodAsset(dp, types.StateDegradation, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateDegradation

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.pod)
	}
}

func TestHandlePodStateDestroy(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			pod   *types.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pod with node"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		pod := getPodAsset(dp, types.StateDestroy, types.EmptyString)
		pod.Meta.Node = "node"

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pod without node"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		pod := getPodAsset(dp, types.StateDestroy, types.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = types.StateDestroy
		s.want.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod, 0)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with multiple different pods states"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		p1 := getPodAsset(dp, types.StateDestroy, types.EmptyString)
		p2 := getPodAsset(dp, types.StateReady, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with multiple equal pods states"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		p1 := getPodAsset(dp, types.StateDestroy, types.EmptyString)
		p2 := getPodAsset(dp, types.StateDestroy, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.pod)
	}
}

func TestHandlePodStateDestroyed(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			pod   *types.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with one pod"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		pod := getPodAsset(dp, types.StateDestroyed, types.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[pod.DeploymentLink()][pod.SelfLink()] = pod

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[pod.DeploymentLink()] = make(map[string]*types.Pod)
		s.want.state.deployment.provision.Status.State = types.StateDestroyed

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with multiple equal pods states"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		p1 := getPodAsset(dp, types.StateDestroyed, types.EmptyString)
		p2 := getPodAsset(dp, types.StateDestroyed, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.want.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2
		s.want.state.deployment.provision.Status.State = types.StateDestroy

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with multiple different pods states"}

		svc := getServiceAsset(types.StateDestroy, types.EmptyString)
		dp := getDeploymentAsset(svc, types.StateDestroy, types.EmptyString)
		p1 := getPodAsset(dp, types.StateDestroyed, types.EmptyString)
		p2 := getPodAsset(dp, types.StateReady, types.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink()] = dp
		s.args.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.args.state.pod.list[p1.DeploymentLink()][p1.SelfLink()] = p1
		s.args.state.pod.list[p2.DeploymentLink()][p2.SelfLink()] = p2

		s.want.err = types.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[p1.DeploymentLink()] = make(map[string]*types.Pod)
		s.want.state.pod.list[p1.DeploymentLink()][p2.SelfLink()] = p2

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.pod)
	}
}
