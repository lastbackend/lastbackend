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

package docker

import (
	"sync"

	"github.com/pkg/errors"
)

// StatusTracker to track status of operations
type StatusTracker interface {
	GetStatus(string) (Status, error)
	SetStatus(string, Status)
}

type memoryStatusTracker struct {
	statuses map[string]Status
	m        sync.Mutex
}

// NewInMemoryTracker returns a StatusTracker that tracks content status in-memory
func NewInMemoryTracker() StatusTracker {
	return &memoryStatusTracker{
		statuses: map[string]Status{},
	}
}

func (t *memoryStatusTracker) GetStatus(ref string) (Status, error) {
	t.m.Lock()
	defer t.m.Unlock()
	status, ok := t.statuses[ref]
	if !ok {
		return Status{}, errors.Wrapf(ErrNotFound, "status for ref %v", ref)
	}
	return status, nil
}

func (t *memoryStatusTracker) SetStatus(ref string, status Status) {
	t.m.Lock()
	t.statuses[ref] = status
	t.m.Unlock()
}
