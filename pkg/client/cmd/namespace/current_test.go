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

package namespace_test

import (
	"github.com/lastbackend/lastbackend/pkg/client/cmd/namespace"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	s "github.com/lastbackend/lastbackend/pkg/client/storage"
	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrent(t *testing.T) {

	const (
		tName       = "test name"
		tDesc       = "test description"
		storageName = "test"
	)

	var (
		err error
		ctx = context.Mock()

		data = n.Namespace{
			Meta: n.NamespaceMeta{
				Name:        tName,
				Description: tDesc,
			},
		}
	)

	storage, err := s.Init()
	assert.NoError(t, err)
	ctx.SetStorage(storage)
	defer (func() {
		storage.Clear()
	})()

	err = storage.Set(storageName, data)
	assert.NoError(t, err)

	nspace, err := namespace.Current()
	assert.NoError(t, err)
	assert.Equal(t, tName, nspace.Meta.Name)
	assert.Equal(t, tDesc, nspace.Meta.Description)
}
