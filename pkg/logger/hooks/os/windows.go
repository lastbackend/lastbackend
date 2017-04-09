// +build windows

package os

import (
	"github.com/Sirupsen/logrus"
)

func SyslogHook(entry *logrus.Entry, network, raddr, tag string) error {
	return nil
}
