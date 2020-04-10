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

package bbolt_test

import (
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/storage/bbolt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Get(t *testing.T) {
	stg, err := bbolt.New(bbolt.Options{Path: path.Join(os.TempDir(), ".lb_test")})
	assert.NoError(t, err, "storage initialize err")
	storage.StorageGetAssets(t, stg)
	stg.Close()
}

func TestStorage_List(t *testing.T) {
	stg, err := bbolt.New(bbolt.Options{Path: path.Join(os.TempDir(), ".lb_test")})
	assert.NoError(t, err, "storage initialize err")
	storage.StorageListAssets(t, stg)
	stg.Close()
}

func TestStorage_Put(t *testing.T) {
	stg, err := bbolt.New(bbolt.Options{Path: path.Join(os.TempDir(), ".lb_test")})
	assert.NoError(t, err, "storage initialize err")
	storage.StoragePutAssets(t, stg)
	stg.Close()
}

func TestStorage_Set(t *testing.T) {
	stg, err := bbolt.New(bbolt.Options{Path: path.Join(os.TempDir(), ".lb_test")})
	assert.NoError(t, err, "storage initialize err")
	storage.StorageSetAssets(t, stg)
	stg.Close()
}

func TestStorage_Del(t *testing.T) {
	stg, err := bbolt.New(bbolt.Options{Path: path.Join(os.TempDir(), ".lb_test")})
	assert.NoError(t, err, "storage initialize err")
	storage.StorageDelAssets(t, stg)
	stg.Close()
}
