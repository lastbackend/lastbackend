//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/tools/log/formatter"
	"github.com/lastbackend/lastbackend/tools/log/hooks"
	"github.com/sirupsen/logrus"
	"os"
)

type Log interface {
	Debug(args ...interface{})
	Debugf(args ...interface{})
	Info(args ...interface{})
	Infof(args ...interface{})
	Warn(args ...interface{})
	Warnf(args ...interface{})
	Error(args ...interface{})
	Errorf(args ...interface{})
}

type Logger struct {
	name  string
	level Level
	log   *logrus.Logger
	lv    Verbose
	ev    Verbose
}

type Level int

func New(level int) *Logger {
	var l = new(Logger)
	l.level = Level(level)
	l.log = logrus.New()
	l.log.Out = os.Stdout
	l.log.Formatter = getJSONFormatter()

	l.lv = l.log
	l.ev = new(Empty)

	if level >= 0 {
		l.log.Level = logrus.DebugLevel
		l.log.Formatter = getTextFormatter()
	}

	return l
}

func (l *Logger) SetLevel(level int) {

	l.level = Level(level)
	if level >= 0 {
		l.log.Level = logrus.DebugLevel
		l.log.Formatter = getTextFormatter()
	}
}

func (l *Logger) EnableFileInfo(skip int) *Logger {
	l.log.Hooks.Add(hooks.ContextHook{Skip: skip})
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
	if level <= l.level {
		return l.lv
	}
	return l.ev
}

func getJSONFormatter() *logrus.JSONFormatter {
	var f = new(logrus.JSONFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05"
	return f
}

func getTextFormatter() *formatter.TextFormatter {
	var f = new(formatter.TextFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05"
	return f
}
