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
	"github.com/lastbackend/lastbackend/internal/master/envs"
	"github.com/lastbackend/lastbackend/internal/master/ipam"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testDeploymentObserver(t *testing.T, name, werr string, wst *ServiceState, state *ServiceState, d *models.Deployment) {
	var (
		ctx = context.Background()
		err error
	)

	stg := envs.Get().GetStorage()

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

		err := deploymentObserve(state, d)
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
				"service status state is different") {
				return
			}

			// check service status message is equal
			if !assert.Equal(t, wst.service.Status.Message, state.service.Status.Message,
				"service status message is different") {
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

			if !assert.NotNil(t, state.deployment.active, "active deployment should be not nil") {
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

		// check pods count

		if !assert.Equal(t,
			len(wst.pod.list),
			len(state.pod.list),
			"pod deployment groups count is different") {
			return
		}

		if d.Status.State != models.StateDestroyed {
			if _, ok := state.pod.list[d.SelfLink().String()]; !ok {
				t.Errorf("pod deployment group not exitst: %s", d.SelfLink())
				return
			}
		}

		if d.Status.State == models.StateDestroyed {
			if _, ok := state.pod.list[d.SelfLink().String()]; ok {
				t.Errorf("pod deployment group should not exitst: %s", d.SelfLink())
				return
			}
		}

		if !assert.Equal(t,
			len(wst.pod.list[d.SelfLink().String()]),
			len(state.pod.list[d.SelfLink().String()]),
			"state pods count not match") {
			return
		}

		if d.Status.State != models.StateDestroy && d.Status.State != models.StateDestroyed && d.Status.State != models.StateWaiting {

			var count = 0
			for _, p := range state.pod.list[d.SelfLink().String()] {
				if p.Status.State == models.StateDestroyed || p.Status.State == models.StateDestroy {
					continue
				}
				count++
			}

			if !assert.Equal(t,
				d.Spec.Replicas,
				count,
				"pods count not match with replicas") {
				return
			}

			return
		}

		if d.Status.State == models.StateWaiting {

			return
		}

		if d.Status.State == models.StateDestroyed {
			if !assert.Equal(t,
				0,
				len(state.pod.list[d.SelfLink().String()]),
				"pods count not match with replicas") {
				return
			}

			return
		}

	})
}

func TestHandleDeploymentStateCreated(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			d     *models.Deployment
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle without pods should create pod"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateCreated, models.EmptyString)
		pod := getPodAsset(dp, models.StateCreated, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateProvision
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pods should scale up pods count"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateCreated, models.EmptyString)
		p1 := getPodAsset(dp, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp, models.StateCreated, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateProvision
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pods should scale down pods count"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateCreated, models.EmptyString)
		p1 := getPodAsset(dp, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp, models.StateError, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2
		dp.Spec.Replicas = 1
		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Spec.Replicas = 1

		s.want.state.deployment.provision.Status.State = models.StateProvision
		s.want.state.deployment.provision.Spec.Replicas = 1

		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pods should create new"}

		svc := getServiceAsset(models.StateCreated, models.EmptyString)

		dp := getDeploymentAsset(svc, models.StateCreated, models.EmptyString)
		p1 := getPodAsset(dp, models.StateDestroy, models.EmptyString)
		p2 := getPodAsset(dp, models.StateProvision, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = models.StateProvision
		s.want.state.deployment.provision.Status.State = models.StateProvision
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "waiting state handle without volumes"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateCreated, models.EmptyString)
		dp.Spec.Template.Volumes = append(dp.Spec.Template.Volumes, &models.SpecTemplateVolume{
			Name: "demo",
			Volume: models.SpecTemplateVolumeClaim{
				Name: "test",
			},
		})

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = models.StateWaiting
		s.want.state.deployment.provision.Status.State = models.StateWaiting
		s.want.state.deployment.provision.Status.Dependencies.Volumes["test"] = models.StatusDependency{
			Name:   "test",
			Type:   models.KindVolume,
			Status: models.StateNotReady,
		}

		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)

		return s
	}())

	for _, tt := range tests {
		testDeploymentObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.d)
	}
}

func TestHandleDeploymentStateProvision(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			d     *models.Deployment
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle without pods should create pod"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		pod := getPodAsset(dp, models.StateCreated, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateProvision
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][pod.SelfLink().String()] = pod

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pods should scale up pods count"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp, models.StateCreated, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision.Status.State = models.StateProvision
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pods should scale down pods count"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp, models.StateError, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2
		dp.Spec.Replicas = 1
		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Spec.Replicas = 1
		s.want.state.deployment.provision.Spec.Replicas = 1
		s.want.state.deployment.provision.Status.State = models.StateProvision

		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pods should create new"}

		svc := getServiceAsset(models.StateCreated, models.EmptyString)

		dp := getDeploymentAsset(svc, models.StateProvision, models.EmptyString)
		p1 := getPodAsset(dp, models.StateDestroy, models.EmptyString)
		p2 := getPodAsset(dp, models.StateProvision, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = models.StateProvision
		s.want.state.deployment.provision.Status.State = models.StateProvision
		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		return s
	}())

	for _, tt := range tests {
		testDeploymentObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.d)
	}
}

func TestHandleDeploymentStateReady(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			d     *models.Deployment
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle without active deployment"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateReady, models.EmptyString)
		p1 := getPodAsset(dp, models.StateReady, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		s.want.state.deployment.provision = nil
		s.want.state.deployment.active = dp
		s.want.state.service.Status.State = models.StateReady

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with active deployment replacement without pods"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp1 := getDeploymentAsset(svc, models.StateReady, models.EmptyString)
		dp2 := getDeploymentAsset(svc, models.StateReady, models.EmptyString)
		p1 := getPodAsset(dp2, models.StateReady, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2
		s.args.state.pod.list[dp1.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp2.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp2.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.d = dp2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		s.want.state.deployment.provision = nil
		s.want.state.deployment.active = dp2

		s.want.state.deployment.list[dp1.SelfLink().String()].Status.State = models.StateDestroyed
		s.want.state.service.Status.State = models.StateReady

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with active deployment replacement with pods"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp1 := getDeploymentAsset(svc, models.StateReady, models.EmptyString)
		dp2 := getDeploymentAsset(svc, models.StateReady, models.EmptyString)

		p1 := getPodAsset(dp1, models.StateReady, models.EmptyString)
		p2 := getPodAsset(dp1, models.StateReady, models.EmptyString)
		p1.Meta.Node = "node"

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2
		s.args.state.pod.list[dp1.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp2.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp1.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp2.SelfLink().String()][p2.SelfLink().String()] = p2

		s.args.d = dp2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		s.want.state.deployment.provision = nil
		s.want.state.deployment.active = dp2

		s.want.state.deployment.list[dp1.SelfLink().String()].Status.State = models.StateDestroy
		s.want.state.service.Status.State = models.StateReady

		return s
	}())

	for _, tt := range tests {
		testDeploymentObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.d)
	}
}

func TestHandleDeploymentStateError(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			d     *models.Deployment
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with service state update"}

		svc := getServiceAsset(models.StateProvision, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateError, models.EmptyString)
		p1 := getPodAsset(dp, models.StateError, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)

		s.want.state.deployment.provision = nil
		s.want.state.deployment.active = dp

		s.want.state.service.Status.State = models.StateError

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle without service state update"}

		svc := getServiceAsset(models.StateReady, models.EmptyString)

		dp1 := getDeploymentAsset(svc, models.StateReady, models.EmptyString)
		dp2 := getDeploymentAsset(svc, models.StateError, models.EmptyString)

		p1 := getPodAsset(dp1, models.StateReady, models.EmptyString)
		p2 := getPodAsset(dp2, models.StateError, models.EmptyString)

		p1.Meta.Node = "node"
		p2.Meta.Node = "node"

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2
		s.args.state.pod.list[dp1.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp2.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp1.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp2.SelfLink().String()][p2.SelfLink().String()] = p2
		s.args.d = dp2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision = nil
		s.want.state.service.Status.State = models.StateReady
		s.want.state.deployment.list[dp2.SelfLink().String()].Status.State = models.StateError

		return s
	}())

	for _, tt := range tests {
		testDeploymentObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.d)
	}
}

func TestHandleDeploymentStateDegradation(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			d     *models.Deployment
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with service state update"}

		svc := getServiceAsset(models.StateReady, models.EmptyString)
		dp1 := getDeploymentAsset(svc, models.StateDegradation, models.EmptyString)
		p1 := getPodAsset(dp1, models.StateDegradation, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.pod.list[dp1.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp1.SelfLink().String()][p1.SelfLink().String()] = p1

		s.args.d = dp1

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = models.StateDegradation

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle without service state update"}

		svc := getServiceAsset(models.StateReady, models.EmptyString)
		dp1 := getDeploymentAsset(svc, models.StateReady, models.EmptyString)
		dp2 := getDeploymentAsset(svc, models.StateDegradation, models.EmptyString)
		p1 := getPodAsset(dp1, models.StateReady, models.EmptyString)
		p2 := getPodAsset(dp1, models.StateDegradation, models.EmptyString)

		p1.Meta.Node = "node"
		p2.Meta.Node = "node"

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2
		s.args.state.pod.list[dp1.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp2.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp1.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp2.SelfLink().String()][p2.SelfLink().String()] = p2
		s.args.d = dp2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = models.StateReady
		s.want.state.deployment.provision = nil
		s.want.state.deployment.list[dp2.SelfLink().String()].Status.State = models.StateDegradation

		return s
	}())

	for _, tt := range tests {
		testDeploymentObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.d)
	}
}

func TestHandleDeploymentStateDestroy(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			d     *models.Deployment
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle active deployment with pods with nodes"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		p1 := getPodAsset(dp, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp, models.StateError, models.EmptyString)

		p1.Meta.Node = "node"
		p2.Meta.Node = "node"

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.active.Status.State = models.StateDestroy

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle provision deployment with pods without nodes"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		p1 := getPodAsset(dp, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp, models.StateError, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = models.StateDestroyed
		s.want.state.deployment.provision = nil
		s.want.state.deployment.list = make(map[string]*models.Deployment, 0)
		s.want.state.pod.list = make(map[string]map[string]*models.Pod, 0)

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle active deployment with one pod without node"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		p1 := getPodAsset(dp, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp, models.StateError, models.EmptyString)
		p1.Meta.Node = "node"

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.active.Status.State = models.StateDestroy
		delete(s.want.state.pod.list[dp.SelfLink().String()], p2.SelfLink().String())

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle provision deployment without pods"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.list = make(map[string]*models.Deployment)
		s.want.state.deployment.provision = nil
		s.want.state.service.Status.State = models.StateDestroyed
		delete(s.want.state.pod.list, dp.SelfLink().String())

		return s
	}())

	for _, tt := range tests {
		testDeploymentObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.d)
	}
}

func TestHandleDeploymentStateDestroyed(t *testing.T) {

	type suit struct {
		name string
		args struct {
			state *ServiceState
			d     *models.Deployment
		}
		want struct {
			err   string
			state *ServiceState
		}
	}

	var tests []suit

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with pods"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		svc.Spec.Replicas = 2

		dp := getDeploymentAsset(svc, models.StateDestroyed, models.EmptyString)
		p1 := getPodAsset(dp, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp, models.StateError, models.EmptyString)

		p1.Meta.Node = "node"
		p2.Meta.Node = "node"

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.provision = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.service.Status.State = models.StateDestroy
		s.want.state.deployment.provision = nil
		s.want.state.deployment.list[dp.SelfLink().String()].Status.State = models.StateDestroy

		s.want.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)
		s.want.state.pod.list[dp.SelfLink().String()][p1.SelfLink().String()] = p1
		s.want.state.pod.list[dp.SelfLink().String()][p2.SelfLink().String()] = p2

		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle with service state change"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		dp := getDeploymentAsset(svc, models.StateDestroyed, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp
		s.args.state.deployment.list[dp.SelfLink().String()] = dp
		s.args.state.pod.list[dp.SelfLink().String()] = make(map[string]*models.Pod)

		s.args.d = dp

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.active = nil
		s.want.state.deployment.list = make(map[string]*models.Deployment)
		s.want.state.service.Status.State = models.StateDestroyed
		delete(s.want.state.deployment.list, dp.SelfLink().String())
		delete(s.want.state.pod.list, dp.SelfLink().String())
		return s
	}())

	tests = append(tests, func() suit {

		s := suit{name: "successful state handle without service state change"}

		svc := getServiceAsset(models.StateDestroy, models.EmptyString)
		svc.Spec.Replicas = 2

		dp1 := getDeploymentAsset(svc, models.StateDestroy, models.EmptyString)
		dp2 := getDeploymentAsset(svc, models.StateDestroyed, models.EmptyString)

		s.args.state = getServiceStateAsset(svc)
		s.args.state.deployment.active = dp1
		s.args.state.deployment.provision = dp2
		s.args.state.deployment.list[dp1.SelfLink().String()] = dp1
		s.args.state.deployment.list[dp2.SelfLink().String()] = dp2
		p1 := getPodAsset(dp1, models.StateCreated, models.EmptyString)
		p2 := getPodAsset(dp1, models.StateCreated, models.EmptyString)
		p1.Meta.Node = "node"
		p2.Meta.Node = "node"

		s.args.state.pod.list[dp1.SelfLink().String()] = make(map[string]*models.Pod)
		s.args.state.pod.list[dp1.SelfLink().String()][p1.SelfLink().String()] = p1
		s.args.state.pod.list[dp1.SelfLink().String()][p2.SelfLink().String()] = p2

		s.args.d = dp2

		s.want.err = models.EmptyString
		s.want.state = getServiceStateCopy(s.args.state)
		s.want.state.deployment.provision = nil
		s.want.state.deployment.list = make(map[string]*models.Deployment)
		s.want.state.deployment.list[dp1.SelfLink().String()] = dp1

		return s
	}())

	for _, tt := range tests {
		testDeploymentObserver(t, tt.name, tt.want.err, tt.want.state, tt.args.state, tt.args.d)
	}
}
