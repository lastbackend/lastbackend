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

package system

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/util/system"
)

// HeartBeat Interval
const heartBeatInterval = 5 // in seconds
const systemLeadTTL = 15

type Process struct {
	// Process storage
	storage storage.IStorage
	// Managed process
	process *models.Process
}

// Process register function
// The main purpose is to register process in the system
// If we need to distribution and need master/replicas, use WaitElected function
func (c *Process) Register(ctx context.Context, kind string, stg storage.IStorage) (*models.Process, error) {

	var (
		err  error
		item = new(models.Process)
	)

	item.Meta.SetDefault()
	item.Meta.Kind = kind

	if item.Meta.Hostname, err = system.GetHostname(); err != nil {
		return item, err
	}

	item.Meta.PID = system.GetPid()
	item.Meta.ID = encodeID(item)

	c.process = item
	c.storage = stg

	opts := storage.GetOpts()
	opts.Ttl = systemLeadTTL
	opts.Force = true

	sl := models.NewProcessSelfLink(kind, c.process.Meta.Hostname, c.process.Meta.PID).String()

	if err := c.storage.Set(ctx, c.storage.Collection().System(), sl, c.process, opts); err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return item, err
		}
	}

	return item, nil
}

// HeartBeat function - updates current process state
// and master election ttl option
func (c *Process) HeartBeat(ctx context.Context, lead chan bool) {

	ticker := time.NewTicker(heartBeatInterval * time.Second)

	opts := storage.GetOpts()
	opts.Ttl = systemLeadTTL
	opts.Force = true

	for range ticker.C {
		// Update process state

		l := false
		opts := storage.GetOpts()
		opts.Ttl = systemLeadTTL
		var process models.Process

		leadKey := fmt.Sprintf("%s/lead", c.process.Meta.Kind)

		err := c.storage.Get(ctx, c.storage.Collection().System(), leadKey, &process, nil)
		if err != nil {

			if errors.Storage().IsErrEntityNotFound(err) {

				err = c.storage.Put(ctx, c.storage.Collection().System(), leadKey, c.process, opts)
				if err != nil && !errors.Storage().IsErrEntityExists(err) {
					return
				}

				l = true

			} else {
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

		sl := models.NewProcessSelfLink(c.process.Meta.Kind, c.process.Meta.Hostname, c.process.Meta.PID).String()
		if err := c.storage.Set(ctx, c.storage.Collection().System(), sl, c.process, opts); err != nil {
			return
		}

		// Check election
		if c.process.Meta.Lead {

			if err := c.storage.Set(ctx, c.storage.Collection().System(), fmt.Sprintf("%s/lead", c.process.Meta.Kind), c.process, opts); err != nil {
				return
			}

		}

		lead <- l

	}
}

// Encode unique ID from pid and process hostname
func encodeID(c *models.Process) string {
	key := fmt.Sprintf("%s|%d", c.Meta.Hostname, c.Meta.PID)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
