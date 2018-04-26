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

package etcd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/coreos/etcd/clientv3"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	s "github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/spf13/viper"
)

func TestStorage_Cluster(t *testing.T) {

	tests := []struct {
		name    string
		want    storage.Cluster
		wantErr bool
	}{
		{"Cluster storage",
			newClusterStorage(),
			false,
		},
		{"Cluster storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Cluster(), tt.want)) {
				t.Errorf("Storage.Cluster() = %v, want %v", got.Cluster(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Cluster(), tt.want) {
					t.Errorf("Storage.Cluster() = %v, want %v", got.Cluster(), tt.want)
					return
				}
			}
		})
	}
}

func TestStorage_Deployment(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Deployment
		wantErr bool
	}{
		{"Deployment storage",
			newDeploymentStorage(),
			false,
		},
		{"Deployment storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Deployment(), tt.want)) {
				t.Errorf("Storage.Deployment() = %v, want %v", got.Deployment(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Deployment(), tt.want) {
					t.Errorf("Storage.Deployment() = %v, want %v", got.Deployment(), tt.want)
				}
			}

		})
	}
}

func TestStorage_Trigger(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Trigger
		wantErr bool
	}{
		{"Trigger storage",
			newTriggerStorage(),
			false,
		},
		{"Trigger storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Trigger(), tt.want)) {
				t.Errorf("Storage.Trigger() = %v, want %v", got.Trigger(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Trigger(), tt.want) {
					t.Errorf("Storage.Trigger() = %v, want %v", got.Trigger(), tt.want)
				}
			}
		})
	}
}

func TestStorage_Node(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Node
		wantErr bool
	}{
		{"Node storage",
			newNodeStorage(),
			false,
		},
		{"Node storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Node(), tt.want)) {
				t.Errorf("Storage.Node() = %v, want %v", got.Node(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Node(), tt.want) {
					t.Errorf("Storage.Node() = %v, want %v", got.Node(), tt.want)
				}
			}
		})
	}
}

func TestStorage_Ingress(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Ingress
		wantErr bool
	}{
		{"Ingress storage",
			newIngressStorage(),
			false,
		},
		{"Ingress storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Ingress(), tt.want)) {
				t.Errorf("Storage.Ingress() = %v, want %v", got.Ingress(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Ingress(), tt.want) {
					t.Errorf("Storage.Ingress() = %v, want %v", got.Ingress(), tt.want)
				}
			}
		})
	}
}

func TestStorage_Secret(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Secret
		wantErr bool
	}{
		{"Secret storage",
			newSecretStorage(),
			false,
		},
		{"Secret storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Secret(), tt.want)) {
				t.Errorf("Storage.Secret() = %v, want %v", got.Secret(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Secret(), tt.want) {
					t.Errorf("Storage.Secret() = %v, want %v", got.Secret(), tt.want)
				}
			}
		})
	}
}

func TestStorage_Namespace(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Namespace
		wantErr bool
	}{
		{"Namespace storage",
			newNamespaceStorage(),
			false,
		},
		{"Namespace storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Namespace(), tt.want)) {
				t.Errorf("Storage.Namespace() = %v, want %v", got.Namespace(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Namespace(), tt.want) {
					t.Errorf("Storage.Namespace() = %v, want %v", got.Namespace(), tt.want)
				}
			}
		})
	}
}

func TestStorage_Route(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Route
		wantErr bool
	}{
		{"Route storage",
			newRouteStorage(),
			false,
		},
		{"Route storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Route(), tt.want)) {
				t.Errorf("Storage.Route() = %v, want %v", got.Route(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Route(), tt.want) {
					t.Errorf("Storage.Route() = %v, want %v", got.Route(), tt.want)
				}
			}
		})
	}
}

func TestStorage_Pod(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Pod
		wantErr bool
	}{
		{"Pod storage",
			newPodStorage(),
			false,
		},
		{"Pod storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Pod(), tt.want)) {
				t.Errorf("Storage.Pod() = %v, want %v", got.Pod(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Pod(), tt.want) {
					t.Errorf("Storage.Pod() = %v, want %v", got.Pod(), tt.want)
				}
			}
		})
	}
}

func TestStorage_Service(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.Service
		wantErr bool
	}{
		{"Service storage",
			newServiceStorage(),
			false,
		},
		{"Service storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Service(), tt.want)) {
				t.Errorf("Storage.Service() = %v, want %v", got.Service(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Service(), tt.want) {
					t.Errorf("Storage.Service() = %v, want %v", got.Service(), tt.want)
				}
			}
		})
	}
}

func TestStorage_Volume(t *testing.T) {

	tests := []struct {
		name    string
		want    storage.Volume
		wantErr bool
	}{
		{"Volume storage",
			newVolumeStorage(),
			false,
		},
		{"Volume storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.Volume(), tt.want)) {
				t.Errorf("Storage.Volume() = %v, want %v", got.Volume(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.Volume(), tt.want) {
					t.Errorf("Storage.Volume() = %v, want %v", got.Volume(), tt.want)
				}
			}
		})
	}
}

func TestStorage_System(t *testing.T) {
	tests := []struct {
		name    string
		want    storage.System
		wantErr bool
	}{
		{"System storage",
			newSystemStorage(),
			false,
		},
		{"System storage nil",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		//in order to prevent "no available endpoints" error
		initV3DummyConf()
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if err != nil || (!tt.wantErr && !reflect.DeepEqual(got.System(), tt.want)) {
				t.Errorf("Storage.System() = %v, want %v", got.System(), tt.want)
				return
			}

			if tt.wantErr {
				got = nil
				if !reflect.DeepEqual(got.System(), tt.want) {
					t.Errorf("Storage.System() = %v, want %v", got.System(), tt.want)
				}
			}
		})
	}
}

func Test_keyCreate(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"key test",
			args{[]string{"test", "test"}},
			"test/test",
		},
		{"key demo",
			args{[]string{"test", "demo"}},
			"test/demo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyCreate(tt.args.args...); got != tt.want {
				t.Errorf("keyCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClient(t *testing.T) {

	initStorage()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test get client dummy",
			args{context.Background()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getClient(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("getClient() got = %v, want store", got)
			}

			if got1 == nil {
				t.Errorf("getClient() got1 = %v, want cancel func", got1)
			}
		})
	}
}

//Test connect to etcd
func Test_clientV3Connect(t *testing.T) {

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{ //should be first!
			"test client v3 connect fail",
			args{context.Background()},
			true,
			clientv3.ErrNoAvailableEndpoints,
		},
		{
			"test client v3 connect successfull",
			args{context.Background()},
			false,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initV3DummyConf()
			} else {
				clearV3Dummyconf()
			}

			_, _, err := v3.GetClient(context.Background())

			if err != nil {
				if !tt.wantErr {
					t.Errorf("clientV3Connect() got error %v, but shouldn't", err)
					return
				}
				if tt.wantErr && (err != tt.err) {
					t.Errorf("clientV3Connect() = %v, want = %v", err, tt.err)
				}
			}
		})
	}

}

func initStorage() {

	var (
		err error
	)
	//move dummy conf init to func
	initV3DummyConf()

	if c.store, c.dfunc, err = v3.GetClient(context.Background()); err != nil {
		log.Errorf("etcd: store initialize err: %s", err)
		return
	}

}

func initV3DummyConf() {
	cfg := v3.Config{}
	cfg.Prefix = "lstbknd"
	cfg.Endpoints = []string{"127.0.0.1:2379"}
	viper.Set("etcd", cfg)
}
func clearV3Dummyconf() {
	cfg := v3.Config{}
	viper.Set("etcd", cfg)
}

func compareMeta(got, want types.Meta) bool {
	result := false
	if (got.Name == want.Name) &&
		(got.Description == want.Description) &&
		(got.SelfLink == want.SelfLink) &&
		reflect.DeepEqual(got.Labels, want.Labels) {
		result = true
	}

	return result
}

//find etcdctrl adsolute path
func getEtcdctrl() string {
	path, lookErr := exec.LookPath("etcdctl")
	if lookErr != nil {
		return ""
	}
	return path
}

func runEtcdPut(path, key, value string) error {

	//ETCDCTL_API=3 etcdctl --endpoints=127.0.0.1:2379 put key value

	var conf = v3.Config{}
	if err := viper.UnmarshalKey("etcd", &conf); err != nil {
		return err
	}
	endpoint := conf.Endpoints
	//fmt.Println("endpoint=", endpoint[0]) //"127.0.0.1:2379"
	//have to use V3 API
	os.Setenv("ETCDCTL_API", "3")
	out, err := exec.Command(path, "--endpoints", endpoint[0], "put", key, value).Output()
	if err != nil {
		return errors.New("etcdctl not found")
	}
	if string(out) != "OK\n" {
		return errors.New("etcdctl put failed")
	}
	return nil
}

//special init storage for watch tests in order to get client
func initStorageWatch() (s.Store, s.DestroyFunc, error) {
	var (
		err error
	)
	initV3DummyConf()

	c.store, c.dfunc, err = v3.GetClient(context.Background())
	if err != nil {
		log.Errorf("etcd: store initialize err: %s", err)
		return nil, nil, err
	}
	return c.store, c.dfunc, nil
}
