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

package etcd

import (
	"path/filepath"
	"github.com/spf13/viper"
	"strings"
	"github.com/lastbackend/lastbackend/pkg/log"
	"testing"
	"reflect"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"context"
)

func TestStorage_Cluster(t *testing.T) {

	tests := []struct {
		name string
		want storage.Cluster
	}{
		{"cluster storage",
			newClusterStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Cluster(), tt.want) {
				t.Errorf("Storage.Cluster() = %v, want %v", got.Cluster(), tt.want)
			}
		})
	}
}

func TestStorage_Deployment(t *testing.T) {
	tests := []struct {
		name string
		want storage.Deployment
	}{
		{"Deployment storage",
			newDeploymentStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Deployment(), tt.want) {
				t.Errorf("Storage.Deployment() = %v, want %v", got.Deployment(), tt.want)
			}
		})
	}
}

func TestStorage_Trigger(t *testing.T) {
	tests := []struct {
		name string
		want storage.Trigger
	}{
		{"cluster storage",
			newTriggerStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Trigger(), tt.want) {
				t.Errorf("Storage.Trigger() = %v, want %v", got.Trigger(), tt.want)
			}
		})
	}
}

func TestStorage_Node(t *testing.T) {
	tests := []struct {
		name string
		want storage.Node
	}{
		{"Node storage",
			newNodeStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Node(), tt.want) {
				t.Errorf("Storage.Node() = %v, want %v", got.Node(), tt.want)
			}
		})
	}
}

func TestStorage_Namespace(t *testing.T) {
	tests := []struct {
		name string
		want storage.Namespace
	}{
		{"Namespace storage",
			newNamespaceStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Namespace(), tt.want) {
				t.Errorf("Storage.Namespace() = %v, want %v", got.Namespace(), tt.want)
			}
		})
	}
}

func TestStorage_Route(t *testing.T) {
	tests := []struct {
		name string
		want storage.Route
	}{
		{"Route storage",
			newRouteStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Route(), tt.want) {
				t.Errorf("Storage.Route() = %v, want %v", got.Route(), tt.want)
			}
		})
	}
}

func TestStorage_Pod(t *testing.T) {
	tests := []struct {
		name string
		want storage.Pod
	}{
		{"Pod storage",
			newPodStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Pod(), tt.want) {
				t.Errorf("Storage.Pod() = %v, want %v", got.Pod(), tt.want)
			}
		})
	}
}

func TestStorage_Service(t *testing.T) {
	tests := []struct {
		name string
		want storage.Service
	}{
		{"Service storage",
			newServiceStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Service(), tt.want) {
				t.Errorf("Storage.Service() = %v, want %v", got.Service(), tt.want)
			}
		})
	}
}

func TestStorage_Volume(t *testing.T) {

	tests := []struct {
		name string
		want storage.Volume
	}{
		{"Volume storage",
			newVolumeStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.Volume(), tt.want) {
				t.Errorf("Storage.Volume() = %v, want %v", got.Volume(), tt.want)
			}
		})
	}
}

func TestStorage_System(t *testing.T) {
	tests := []struct {
		name string
		want storage.System
	}{
		{"System storage",
			newSystemStorage(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(); (err != nil) || !reflect.DeepEqual(got.System(), tt.want) {
				t.Errorf("Storage.System() = %v, want %v", got.System(), tt.want)
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
		{"test get client dummy",
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


func initStorage () {
	c := "../../../contrib/test/config.yml"
	abs, err := filepath.Abs(c)
	if err != nil {
		log.Errorf("Error reading filepath: %s \n", err)
	}

	// get the config name
	base := filepath.Base(abs)

	// get the path
	path := filepath.Dir(abs)

	viper.SetConfigName(strings.Split(base, ".")[0])
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %s\n", err)
	}
}
