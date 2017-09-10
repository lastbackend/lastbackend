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
	"github.com/lastbackend/lastbackend/pkg/logger/formatter"
	"github.com/lastbackend/lastbackend/pkg/logger/hooks"
	"github.com/sirupsen/logrus"
	"os"
)

var l *Logger

type Logger struct {
	ILogger

	name  string
	level Level
	log   *logrus.Logger
}

type Level int

func New(name string, level int) *Logger {
	l = new(Logger)
	l.name = name
	l.level = Level(level)
	l.log = logrus.New()
	l.log.Out = os.Stdout
	l.log.Formatter = getJSONFormatter()

	if level > 0 {
		l.log.Level = logrus.DebugLevel
		l.log.Formatter = getTextFormatter(l.name)
	}

	return l
}

func (l *Logger) EnableFileInfo(skip int) *Logger {
	l.log.Hooks.Add(hooks.ContextHook{skip})
	return l
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

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
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

func (l *Logger) V(level Level) Verbose {
	return Verbose(l.level >= level)
}

func getJSONFormatter() *logrus.JSONFormatter {
	var formatter = new(logrus.JSONFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	return formatter
}

func getTextFormatter(name string) *formatter.TextFormatter {
	var f = new(formatter.TextFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05"
	f.Name = name
	return f
}
