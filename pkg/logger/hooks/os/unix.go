// +build !windows
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

package os

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

var hook *logrus_syslog.SyslogHook

func SyslogHook(entry *logrus.Entry, network, raddr, tag string) error {
	var err error
	if hook == nil {
		hook, err = logrus_syslog.NewSyslogHook(network, raddr, syslog.LOG_DEBUG, tag)
		entry.Logger.Hooks.Add(hook)
	}
	return err
}
