// +build windows

package os

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
)

func SyslogHook(entry *logrus.Entry, network, raddr, tag string) error {
	return nil
}
