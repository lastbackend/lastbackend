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
	_os "github.com/lastbackend/lastbackend/pkg/logger/os"
	"os"
	"path"
	"runtime"
)

type Logger struct {
	log   *logrus.Logger
	entry *logrus.Entry
}

func New(debug bool) *Logger {
	l := new(Logger)
	l.log = logrus.New()
	l.log.Out = os.Stdout
	l.log.Formatter = getJSONFormatter()
	l.entry = l.log.WithFields(getFileLocate(logrus.Fields{}))

	_os.SetSyslog(l.log)

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
	l.entry.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args)
}

func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args)
}

func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args)
}

func (l *Logger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args)
}

func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args)
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

func getFileLocate(fields logrus.Fields) logrus.Fields {
	if pc, file, line, ok := runtime.Caller(2); ok {
		funcName := runtime.FuncForPC(pc).Name()
		fields["func"] = path.Base(funcName)
		fields["file"] = file
		fields["line"] = line
	}
	return fields
}
