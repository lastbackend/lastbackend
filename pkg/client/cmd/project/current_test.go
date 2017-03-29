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

package project_test

import (
	"github.com/lastbackend/lastbackend/pkg/client/storage"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"testing"
	"time"
	"github.com/lastbackend/lastbackend/pkg/api/types"
)

func TestCurrent(t *testing.T) {

	const token = "mocktoken"

	var (
		err  error
		ctx  = context.Mock()
		data = types.Project{
			Name:        "mock_name",
			Created:     time.Now(),
			Updated:     time.Now(),
			User:        "mock_demo",
			Description: "sample description",
		}
	)

	ctx.Storage, err = storage.Init()
	if err != nil {
		t.Error(err)
		return
	}
	defer (func() {
		err = ctx.Storage.Clear()
		if err != nil {
			t.Error(err)
			return
		}
		err = ctx.Storage.Close()
		if err != nil {
			t.Error(err)
			return
		}
	})()

	ctx.Token = token

	if err != nil {
		t.Error(err)
		return
	}
	defer ctx.Storage.Close()

	err = ctx.Storage.Set("project", data)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = project.Current()
	if err != nil {
		t.Error(err)
		return
	}
}
