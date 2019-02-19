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

package docker

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	protoio "github.com/gogo/protobuf/io"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/proxy"
	"github.com/pkg/errors"
	"github.com/tonistiigi/fifo"
	"io"
	"path"
	"sync"
	"syscall"
	"time"
)

type driver struct {
	mu     sync.Mutex
	logs   map[string]*logPair
	idx    map[string]*logPair
	logger logger.Logger
	proxy  *proxy.Client
}

type logPair struct {
	active  bool
	file    string
	info    logger.Info
	Message types.LogMessage
	stream  io.ReadCloser
}

const (
	hostAddr = "127.0.0.1:2963"
)

func NewDriver() *driver {
	return &driver{
		proxy: proxy.NewClient("driver", hostAddr, nil),
		logs:  make(map[string]*logPair),
		idx:   make(map[string]*logPair),
	}
}

func (d *driver) StartLogging(file string, logCtx logger.Info) error {

	d.mu.Lock()
	if _, exists := d.logs[path.Base(file)]; exists {
		d.mu.Unlock()
		return fmt.Errorf("logger for %q already exists", file)
	}
	d.mu.Unlock()

	fmt.Println("start logging:> ", path.Base(file), logCtx.ContainerName)
	stream, err := fifo.OpenFifo(context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		return errors.Wrapf(err, "error opening logger fifo: %q", file)
	}

	tag, err := loggerutils.ParseLogTag(logCtx, loggerutils.DefaultTemplate)
	if err != nil {
		return err
	}

	extra, err := logCtx.ExtraAttributes(nil)
	if err != nil {
		return err
	}

	hostname, err := logCtx.Hostname()
	if err != nil {
		return err
	}

	msg := types.LogMessage{
		ContainerId:      logCtx.FullID(),
		ContainerName:    logCtx.Name(),
		Selflink:         logCtx.ContainerLabels[types.ContainerTypeLBC],
		ContainerCreated: types.JsonTime{logCtx.ContainerCreated},
		Tag:              tag,
		Extra:            extra,
		Host:             hostname,
	}

	d.mu.Lock()
	lp := &logPair{true, file, logCtx, msg, stream}
	d.logs[path.Base(file)] = lp
	d.idx[logCtx.ContainerID] = lp
	d.mu.Unlock()

	go consumeLog(d, lp)
	return nil
}

func (d *driver) StopLogging(file string) error {
	fmt.Println("Stop logging >", path.Base(file))
	d.mu.Lock()
	lp, ok := d.logs[path.Base(file)]
	if ok {
		lp.active = false
		delete(d.logs, path.Base(file))
		delete(d.idx, lp.info.ContainerID)
	} else {
		log.Errorf("Failed to stop logging. File %q is not active", file)
	}
	d.mu.Unlock()
	return nil
}

func consumeLog(d *driver, lp *logPair) {

	var buf logdriver.LogEntry

	dec := protoio.NewUint32DelimitedReader(lp.stream, binary.BigEndian, 1e6)
	defer func() { _ = dec.Close() }()
	defer shutdownLogPair(lp)

	for {

		if !lp.active {
			return
		}

		err := dec.ReadMsg(&buf)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				dec = protoio.NewUint32DelimitedReader(lp.stream, binary.BigEndian, 1e6)
				continue
			}
		}

		lp.Message.Data = string(buf.Line[:])
		lp.Message.Timestamp = types.JsonTime{time.Now()}

		msg, err := json.Marshal(lp.Message)
		if err != nil {
			continue
		}

		if err := d.proxy.Send(msg); err != nil {
			continue
		}

		buf.Reset()
	}
}

func shutdownLogPair(lp *logPair) {
	if lp.stream != nil {
		_ = lp.stream.Close()
	}

	lp.active = false
}
