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

package etcd_test

import (
	"testing"

	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestStorage_Get(t *testing.T) {
	stg, err := etcd.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageGetAssets(t, stg)
}

func TestStorage_List(t *testing.T) {
	stg, err := etcd.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageListAssets(t, stg)
}

func TestStorage_Map(t *testing.T) {
	stg, err := etcd.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageMapAssets(t, stg)
}

func TestStorage_Put(t *testing.T) {
	stg, err := etcd.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StoragePutAssets(t, stg)
}

func TestStorage_Set(t *testing.T) {
	stg, err := etcd.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageSetAssets(t, stg)
}

func TestStorage_Del(t *testing.T) {
	stg, err := etcd.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageDelAssets(t, stg)
}

func init() {

	cfg := v3.Config{}
	cfg.Prefix = "lstbknd"
	cfg.Endpoints = []string{"127.0.0.1:2379"}
	viper.Set("etcd", cfg)

}
