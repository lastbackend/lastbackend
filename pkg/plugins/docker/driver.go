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
	"fmt"
	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	protoio "github.com/gogo/protobuf/io"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/pkg/errors"
	"github.com/tonistiigi/fifo"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
)

type driver struct {
	mu     sync.Mutex
	logs   map[string]*logPair
	idx    map[string]*logPair
	logger logger.Logger
}

type logPair struct {
	active  bool
	file    string
	info    logger.Info
	logLine jsonLogLine
	stream  io.ReadCloser
}

func NewDriver() *driver {
	return &driver{
		logs: make(map[string]*logPair),
		idx:  make(map[string]*logPair),
	}
}

func (d *driver) StartLogging(file string, logCtx logger.Info) error {

	d.mu.Lock()
	if _, exists := d.logs[file]; exists {
		d.mu.Unlock()
		return fmt.Errorf("logger for %q already exists", file)
	}
	d.mu.Unlock()

	if logCtx.LogPath == "" {
		logCtx.LogPath = filepath.Join("/var/log/docker", logCtx.ContainerID)
	}
	if err := os.MkdirAll(filepath.Dir(logCtx.LogPath), 0755); err != nil {
		return errors.Wrap(err, "error setting up logger dir")
	}

	log.Debug("start logging")
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

	logLine := jsonLogLine{
		ContainerId:      logCtx.FullID(),
		ContainerName:    logCtx.Name(),
		ContainerCreated: jsonTime{logCtx.ContainerCreated},
		ImageId:          logCtx.ImageFullID(),
		ImageName:        logCtx.ImageName(),
		Command:          logCtx.Command(),
		Tag:              tag,
		Extra:            extra,
		Host:             hostname,
	}

	d.mu.Lock()
	lp := &logPair{true, file, logCtx, logLine, stream}
	d.logs[file] = lp
	d.idx[logCtx.ContainerID] = lp
	d.mu.Unlock()

	go consumeLog(lp)
	return nil
}

func (d *driver) StopLogging(file string) error {
	log.Debug("Stop logging")
	d.mu.Lock()
	lp, ok := d.logs[path.Base(file)]
	if ok {
		lp.active = false
		delete(d.logs, path.Base(file))
	} else {
		log.Errorf("Failed to stop logging. File %q is not active", file)
	}
	d.mu.Unlock()
	return nil
}

func consumeLog(lp *logPair) {

	var buf logdriver.LogEntry

	dec := protoio.NewUint32DelimitedReader(lp.stream, binary.BigEndian, 1e6)
	defer dec.Close()
	defer shutdownLogPair(lp)

	for {
		if !lp.active {
			log.Debug("shutting down logger goroutine due to stop request")
			return
		}

		err := dec.ReadMsg(&buf)
		if err != nil {
			if err == io.EOF {
				log.Debug("shutting down logger goroutine due to file EOF")
				return
			} else {
				log.Warn("error reading from FIFO, trying to continue")
				dec = protoio.NewUint32DelimitedReader(lp.stream, binary.BigEndian, 1e6)
				continue
			}
		}

		err = logMessage(lp, buf.Line)
		if err != nil {
			log.Warn("error logging message, dropping it and continuing")
		}

		buf.Reset()
	}
}

func (d *driver) ReadLogs(info logger.Info, config logger.ReadConfig) (io.ReadCloser, error) {
	d.mu.Lock()
	lf, exists := d.idx[info.ContainerID]
	d.mu.Unlock()
	if !exists {
		return nil, fmt.Errorf("logger does not exist for %s", info.ContainerID)
	}

	r, w := io.Pipe()
	lr, ok := lf.stream.(logger.LogReader)
	if !ok {
		return nil, fmt.Errorf("logger does not support reading")
	}

	go func() {
		watcher := lr.ReadLogs(config)

		enc := protoio.NewUint32DelimitedWriter(w, binary.BigEndian)
		defer enc.Close()
		defer watcher.ConsumerGone()

		var buf logdriver.LogEntry
		for {
			select {
			case msg, ok := <-watcher.Msg:
				if !ok {
					w.Close()
					return
				}

				buf.Line = msg.Line
				buf.TimeNano = msg.Timestamp.UnixNano()
				buf.Source = msg.Source

				if err := enc.WriteMsg(&buf); err != nil {
					w.CloseWithError(err)
					return
				}
			case err := <-watcher.Err:
				w.CloseWithError(err)
				return
			}

			buf.Reset()
		}
	}()

	return r, nil
}

func shutdownLogPair(lp *logPair) {
	if lp.stream != nil {
		lp.stream.Close()
	}

	lp.active = false
}
