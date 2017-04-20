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

package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/lastbackend/lastbackend/pkg/logger/hooks"
	"os"
)

type Logger struct {
	log *logrus.Logger
}

func New(debug bool, skip int) *Logger {
	l := new(Logger)
	l.log = logrus.New()
	l.log.Out = os.Stdout
	l.log.Formatter = getJSONFormatter()

	l.log.Hooks.Add(hooks.SyslogHook{})
	//l.log.Hooks.Add(hooks.ContextHook{skip})

	if debug {
		l.SetDebugLevel()
		l.Debug("Logger debug mode enabled")
	}

	return l
}

func (l *Logger) SetDebugLevel() {
	l.log.Level = logrus.DebugLevel
	l.log.Formatter = getTextFormatter()
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func getTextFormatter() *logrus.TextFormatter {
	var formatter = new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	return formatter
}

func getJSONFormatter() *logrus.JSONFormatter {
	var formatter = new(logrus.JSONFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	return formatter
}
