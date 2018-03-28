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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

// App Service type for interface in interfaces folder
type SystemStorage struct {
	storage.System
	data map[string]struct {
		lead    string
		process string
	}
}

func (s *SystemStorage) ProcessSet(ctx context.Context, process *types.Process) error {
	s.data[process.Meta.Kind] = struct{ lead, process string }{process: process.Meta.ID}
	return nil
}

func (s *SystemStorage) Elect(ctx context.Context, process *types.Process) (bool, error) {
	s.data[process.Meta.Kind] = struct{ lead, process string }{lead: process.Meta.ID}
	return true, nil
}

func (s *SystemStorage) ElectUpdate(ctx context.Context, process *types.Process) error {
	s.data[process.Meta.Kind] = struct{ lead, process string }{lead: process.Meta.ID}
	return nil
}

func (s *SystemStorage) ElectWait(ctx context.Context, process *types.Process, lead chan bool) error {
	return nil
}

// Clear system storage
func (s *SystemStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]struct {
		lead    string
		process string
	})
	return nil
}

func newSystemStorage() *SystemStorage {
	s := new(SystemStorage)
	s.data = make(map[string]struct {
		lead    string
		process string
	})
	return s
}
