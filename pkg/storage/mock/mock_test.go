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

package mock_test

import (
	"testing"

	"github.com/lastbackend/lastbackend/pkg/storage/mock"
	"github.com/lastbackend/lastbackend/pkg/storage/test"
	"github.com/stretchr/testify/assert"
)

func TestStorage_Get(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	test_test.StorageGetAssets(t, stg)
}

func TestStorage_List(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	test_test.StorageListAssets(t, stg)
}

func TestStorage_Map(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	test_test.StorageMapAssets(t, stg)
}

func TestStorage_Put(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	test_test.StoragePutAssets(t, stg)
}

func TestStorage_Set(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	test_test.StorageSetAssets(t, stg)
}

func TestStorage_Del(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	test_test.StorageDelAssets(t, stg)
}
