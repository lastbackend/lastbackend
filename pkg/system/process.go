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

package system

import (
	"encoding/base64"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/common/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"strconv"
	"strings"
	"time"
)

// HeartBeat Interval
const heartBeatInterval = 10 // in seconds

type Process struct {
	// Process operations context
	ctx context.IContext

	// Managed process
	process *types.Process
}

// Process register function
// The main purpose is to register process in the system
// If we need to distribution and need master/replicas, use WaitElected function
func (c *Process) Register(ctx context.IContext, kind string) (*types.Process, error) {

	var (
		err  error

		log  = ctx.GetLogger()
		stg  = ctx.GetStorage()

		item = new(types.Process)
	)

	log.Debugf("System: Process: Register: %s", kind)
	item.Meta.SetDefault()
	item.Meta.Kind = kind

	if item.Meta.Hostname, err = system.GetHostname(); err != nil {
		log.Errorf("System: Process: Register: get hostname: %s", err.Error())
		return item, err
	}

	item.Meta.PID = system.GetPid()
	item.Meta.ID = encodeID(item)

	c.ctx = ctx
	c.process = item

	if err := stg.System().ProcessSet(c.ctx.Background(), c.process); err != nil {
		log.Errorf("System: Process: Register: %s", err.Error())
		return item, err
	}

	go c.HeartBeat()
	return item, nil
}

// HeartBeat function - updates current process state
// and master election ttl option
func (c *Process) HeartBeat() {

	var (
		log = c.ctx.GetLogger()
		stg = c.ctx.GetStorage()
	)

	log.Debugf("System: Process: Start HeartBeat for: %s", c.process.Meta.Kind)
	ticker := time.NewTicker(heartBeatInterval * time.Second)
	for range ticker.C {
		// Update process state
		log.Debug("System: Process: Beat")
		if err := stg.System().ProcessSet(c.ctx.Background(), c.process); err != nil {

		}
		// Check election
		if c.process.Meta.Lead {
			log.Debug("System: Process: Beat: Lead TTL update")
			stg.System().ElectUpdate(c.ctx.Background(), c.process)
		}

	}
}

// WaitElected function used for election waiting if
// master/replicas type of process used
func (c *Process) WaitElected(lead chan bool) error {
	var (
		log = c.ctx.GetLogger()
		stg = c.ctx.GetStorage()
		ld  = make(chan bool)
	)

	log.Debug("System: Process: Wait for election")
	l, err := stg.System().Elect(c.ctx.Background(), c.process)
	if err != nil {
		return err
	}

	if l {
		log.Debug("System: Process: Set as Lead")
		c.process.Meta.Lead = true
		c.process.Meta.Slave = false
		lead <- true
	}

	go func() {
		for {
			select {
			case l := <-ld:
				c.process.Meta.Lead = l
				c.process.Meta.Slave = !l
				lead <- l
			}
		}
	}()

	return stg.System().ElectWait(c.ctx.Background(), c.process, ld)
}

// Encode unique ID from pid and process hostname
func encodeID(c *types.Process) string {
	key := fmt.Sprintf("%s|%d", c.Meta.Hostname, c.Meta.PID)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

// Decode ID into hostname and pid
func decodeID(id string) (int, string, error) {

	var (
		key      []byte
		pid      int
		err      error
		hostname string
	)

	key, err = base64.StdEncoding.DecodeString(id)
	if err != nil {
		return pid, hostname, err
	}

	parts := strings.Split(string(key), "|")
	if len(parts) == 2 {
		hostname = parts[0]
		pid, _ = strconv.Atoi(parts[1])
	}

	return pid, hostname, nil
}
