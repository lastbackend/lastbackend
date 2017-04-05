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

package hooks

import (
	"github.com/Sirupsen/logrus"
	"github.com/lastbackend/lastbackend/pkg/logger/hooks/os"
	"log/syslog"
	"path"
	"runtime"
	"strings"
)

type ContextHook struct {
	Skip int
}

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	if hook.Skip == 0 {
		hook.Skip = 8
	}

	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(hook.Skip, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/Sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["func"] = path.Base(name)
			entry.Data["file"] = file
			entry.Data["line"] = line
			break
		}
	}
	return nil
}

type SyslogHook struct {
	Tag     string
	Network string
	Raddr   string
}

func (hook SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook SyslogHook) Fire(entry *logrus.Entry) error {
	return os.SyslogHook(entry, hook.Network, hook.Raddr, syslog.LOG_DEBUG, hook.Tag)
}
