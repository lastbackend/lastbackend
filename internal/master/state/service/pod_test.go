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

package service

import (
	"context"
	envs2 "github.com/lastbackend/lastbackend/internal/master/envs"
	ipam2 "github.com/lastbackend/lastbackend/internal/master/ipam"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testPodObserver(t *testing.T, name, werr string, wst *ServiceState, state *ServiceState, p *models.Pod) {

	var (
		ctx = context.Background()
		err error
	)

	stg := envs2.Get().GetStorage()

	ipm, _ := ipam2.New("")
	envs2.Get().SetIPAM(ipm)

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
		if werr != models.EmptyString {

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

			if _, ok := state.deployment.list[d.SelfLink().String()]; !ok {
				t.Errorf("deployment not found %s", d.SelfLink())
				return
			}

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
				case models.StateProvision:
					if !assert.NotEqual(t, state.pod.list[d][l].Meta.Node, models.EmptyString,
						"node value should not be empty") {
						return
					}
				case models.StateDestroyed:
					if !assert.Equal(t, state.pod.list[d][l].Meta.Node, models.EmptyString,
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
			pod   *models.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateCreated, models.EmptyString)
		pod := getPodAsset(dp, models.StateCreated, models.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()].Status.State = models.StateProvision

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
			pod   *models.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateReady, models.EmptyString)
		pod := getPodAsset(dp, models.StateProvision, models.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateProvision

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, set to destroy"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		pod := getPodAsset(dp, models.StateProvision, models.EmptyString)
		pod.Meta.Node = "demo"
		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()].Status.State = models.StateDestroy

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, set to destroyed"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		pod := getPodAsset(dp, models.StateProvision, models.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod, 0)

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
			pod   *models.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with change deployment state"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		pod := getPodAsset(dp, models.StateReady, models.EmptyString)
		pod.Meta.Node = "node"

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateReady

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with change deployment state to ready"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateReady, models.EmptyString)
		p2 := getPodAsset(dp, models.StateReady, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateReady

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, stay deployment state"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateProvision, models.EmptyString)
		p1.Meta.Node = "node"

		p2 := getPodAsset(dp, models.StateReady, models.EmptyString)

		s.args.pod = p2

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateProvision

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
			pod   *models.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to error"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		pod := getPodAsset(dp, models.StateError, models.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateError

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to degradation"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateError, models.EmptyString)
		p2 := getPodAsset(dp, models.StateReady, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateDegradation

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to error with many pods"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateError, models.EmptyString)
		p2 := getPodAsset(dp, models.StateError, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateError

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
			pod   *models.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to degradation"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		pod := getPodAsset(dp, models.StateDegradation, models.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateDegradation

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to degradation"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateDegradation, models.EmptyString)
		p2 := getPodAsset(dp, models.StateReady, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateDegradation

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle, change deployment state to degradation"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateDegradation, models.EmptyString)
		p2 := getPodAsset(dp, models.StateDegradation, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateDegradation

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
			pod   *models.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pod with node"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		pod := getPodAsset(dp, models.StateDestroy, models.EmptyString)
		pod.Meta.Node = "node"

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pod without node"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		pod := getPodAsset(dp, models.StateDestroy, models.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateDestroy
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod, 0)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with multiple different pods states"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		p1 := getPodAsset(dp, models.StateDestroy, models.EmptyString)
		p2 := getPodAsset(dp, models.StateReady, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with multiple equal pods states"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		p1 := getPodAsset(dp, models.StateDestroy, models.EmptyString)
		p2 := getPodAsset(dp, models.StateDestroy, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
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
			pod   *models.Pod
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with one pod"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		pod := getPodAsset(dp, models.StateDestroyed, models.EmptyString)

		s.args.pod = pod

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.deployment.provision.Status.State = models.StateDestroyed

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with multiple equal pods states"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		p1 := getPodAsset(dp, models.StateDestroyed, models.EmptyString)
		p2 := getPodAsset(dp, models.StateDestroyed, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2
		s.want.state.deployment.provision.Status.State = models.StateDestroy

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with multiple different pods states"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		p1 := getPodAsset(dp, models.StateDestroyed, models.EmptyString)
		p2 := getPodAsset(dp, models.StateReady, models.EmptyString)

		s.args.pod = p1

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		return s
	}())

	for _, tt := range tests {
		testPodObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.pod)
	}
}
