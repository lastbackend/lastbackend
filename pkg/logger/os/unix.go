// +build !windows

package os

import "github.com/Sirupsen/logrus"

func SetSyslog(l *logrus.Logger) {
    if hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, ""); err == nil {
        l.Hooks.Add(hook)
    }
}