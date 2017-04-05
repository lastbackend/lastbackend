// +build windows

package os

import (
	"github.com/Sirupsen/logrus"
)

func SyslogHook(_ *logrus.Entry) error {
	return nil
}
