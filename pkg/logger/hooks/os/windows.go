// +build windows

package os

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
)

func SyslogHook(network, raddr string, tag string) (*logrus_syslog.SyslogHook, error) {
	return nil
}
