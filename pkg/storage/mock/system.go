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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

const systemStorage = "system"
const systemLeadTTL = 11

// App Service type for interface in interfaces folder
type SystemStorage struct {
	storage.System
}

func (s *SystemStorage) ProcessSet(ctx context.Context, process *types.Process) error {
	return nil
}

func (s *SystemStorage) Elect(ctx context.Context, process *types.Process) (bool, error) {
	return false, nil
}

func (s *SystemStorage) ElectUpdate(ctx context.Context, process *types.Process) error {
	return nil
}

func (s *SystemStorage) ElectWait(ctx context.Context, process *types.Process, lead chan bool) error {
	return nil
}

func newSystemStorage() *SystemStorage {
	s := new(SystemStorage)
	return s
}
