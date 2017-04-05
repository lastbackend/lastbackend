// +build !windows

package os

import (
	"github.com/Sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

var hook *logrus_syslog.SyslogHook

func SyslogHook(network, raddr string, tag string) (*logrus_syslog.SyslogHook, error) {
	var err error
	if hook == nil {
		hook, err = logrus_syslog.NewSyslogHook(network, raddr, syslog.LOG_DEBUG, tag)
	}
	return hook, err
}
