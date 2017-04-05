// +build !windows

package os

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

var hook *logrus_syslog.SyslogHook

func SyslogHook(entry *logrus.Entry, network, raddr string, priority syslog.Priority, tag string) error {
	var err error
	if hook == nil {
		hook, err = logrus_syslog.NewSyslogHook(network, raddr, priority, tag)
		if err == nil {
			entry.Logger.Hooks.Add(hook)
		}
	}
	return err
}
