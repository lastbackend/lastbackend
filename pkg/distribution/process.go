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

package distribution

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logProcessPrefix = "distribution:process"
)

type IProcess interface {
	ProcessSet(process *types.Process) error
	Elect(p *types.Process) (bool, error)
	ElectWait(p *types.Process, l chan bool) error
	ElectUpdate(p *types.Process) error
}

type Process struct {
	context context.Context
	storage.Storage
}

func (p *Process) ProcessSet(process *types.Process) error {
	return p.Storage.System().ProcessSet(p.context, process)
}

func (p *Process) Elect(process *types.Process) (bool, error) {
	return p.Storage.System().Elect(p.context, process)
}

func (p *Process) ElectWait(process *types.Process, event chan *types.Event) error {
	return p.Storage.System().ElectWait(p.context, process, event)
}

func (p *Process) ElectUpdate(process *types.Process) error {
	return p.Storage.System().ElectUpdate(p.context, process)
}
