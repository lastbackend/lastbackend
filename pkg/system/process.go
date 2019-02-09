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

package system

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/system"
)

// HeartBeat Interval
const heartBeatInterval = 5 // in seconds
const systemLeadTTL = 15
const logLevel = 7

type Process struct {
	// Process storage
	storage storage.Storage
	// Managed process
	process *types.Process
}

// Process register function
// The main purpose is to register process in the system
// If we need to distribution and need master/replicas, use WaitElected function
func (c *Process) Register(ctx context.Context, kind string, stg storage.Storage) (*types.Process, error) {

	var (
		err  error
		item = new(types.Process)
	)

	log.V(logLevel).Debugf("System: Process: Register: %s", kind)
	item.Meta.SetDefault()
	item.Meta.Kind = kind

	if item.Meta.Hostname, err = system.GetHostname(); err != nil {
		log.Errorf("System: Process: Register: get hostname: %s", err.Error())
		return item, err
	}

	item.Meta.PID = system.GetPid()
	item.Meta.ID = encodeID(item)

	c.process = item
	c.storage = stg

	opts := storage.GetOpts()
	opts.Ttl = systemLeadTTL
	opts.Force = true

	sl := types.NewProcessSelfLink(kind, c.process.Meta.Hostname, c.process.Meta.PID).String()

	if err := c.storage.Set(ctx, c.storage.Collection().System(), sl, c.process, opts); err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("System: Process: Register: %s", err.Error())
			return item, err
		}
	}

	return item, nil
}

// HeartBeat function - updates current process state
// and master election ttl option
func (c *Process) HeartBeat(ctx context.Context, lead chan bool) {

	log.V(logLevel).Debugf("System: Process: Start HeartBeat for: %s", c.process.Meta.Kind)
	ticker := time.NewTicker(heartBeatInterval * time.Second)

	opts := storage.GetOpts()
	opts.Ttl = systemLeadTTL
	opts.Force = true

	for range ticker.C {
		// Update process state
		log.V(logLevel).Debug("System: Process: Beat")

		l := false
		opts := storage.GetOpts()
		opts.Ttl = systemLeadTTL
		var process types.Process

		leadKey := fmt.Sprintf("%s/lead", c.process.Meta.Kind)

		err := c.storage.Get(ctx, c.storage.Collection().System(), leadKey, &process, nil)
		if err != nil {

			if errors.Storage().IsErrEntityNotFound(err) {

				err = c.storage.Put(ctx, c.storage.Collection().System(), leadKey, c.process, opts)
				if err != nil && !errors.Storage().IsErrEntityExists(err) {
					log.V(logLevel).Errorf("System: Process: create process ttl err: %s", err.Error())
					return
				}

				l = true

			} else {
				log.Errorf("System: Process: get lead process: %s", err.Error())
				return
			}

		} else {
			l = process.Meta.Hostname == c.process.Meta.Hostname && process.Meta.PID == c.process.Meta.PID
		}

		if l {
			c.process.Meta.Lead = true
			c.process.Meta.Slave = false
		} else {
			c.process.Meta.Lead = false
			c.process.Meta.Slave = true
		}

		sl := types.NewProcessSelfLink(c.process.Meta.Kind, c.process.Meta.Hostname, c.process.Meta.PID).String()
		if err := c.storage.Set(ctx, c.storage.Collection().System(), sl, c.process, opts); err != nil {
			log.Errorf("System: Process: Register: %s", err.Error())
			return
		}

		// Check election
		if c.process.Meta.Lead {
			log.V(logLevel).Debug("System: Process: Beat: Lead TTL update")

			if err := c.storage.Set(ctx, c.storage.Collection().System(), fmt.Sprintf("%s/lead", c.process.Meta.Kind), c.process, opts); err != nil {
				log.Errorf("System: Process: update process: %s", err.Error())
				return
			}

		}

		lead <- l

	}
}

// Encode unique ID from pid and process hostname
func encodeID(c *types.Process) string {
	key := fmt.Sprintf("%s|%d", c.Meta.Hostname, c.Meta.PID)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
