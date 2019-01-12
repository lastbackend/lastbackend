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

package log

import (
	"github.com/lastbackend/lastbackend/pkg/log/logger"
)

var l *logger.Logger

func init() {
	l = logger.New(-1)
}

// Initialize loggers map

func New(level int) *logger.Logger {
	l.SetLevel(level)
	return l
}

func Debug(args ...interface{}) {
	l.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func Info(args ...interface{}) {
	l.Info(args...)
}

func Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func Warn(args ...interface{}) {
	l.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func Error(args ...interface{}) {
	l.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	l.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	l.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	l.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	l.Panicf(format, args...)
}

func V(level logger.Level) logger.Verbose {
	return l.V(level)
}
