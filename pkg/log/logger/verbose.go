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

import "github.com/Sirupsen/logrus"

// Verbose is a boolean type that implements Infof (like Printf) etc.
type Logrus struct {
	log *logrus.Logger
}

func (v Logrus) Debug(args ...interface{}) {
	v.log.Debug(args...)
}

func (v Logrus) Debugf(format string, args ...interface{}) {
	v.log.Debugf(format, args...)
}

func (v Logrus) Info(args ...interface{}) {
	v.log.Info(args...)
}

func (v Logrus) Infof(format string, args ...interface{}) {
	v.log.Infof(format, args...)
}

func (v Logrus) Warn(args ...interface{}) {
	v.log.Warn(args...)
}

func (v Logrus) Warnf(format string, args ...interface{}) {
	v.log.Warnf(format, args...)
}

func (v Logrus) Error(args ...interface{}) {
	v.log.Error(args...)
}

func (v Logrus) Errorf(format string, args ...interface{}) {
	v.log.Errorf(format, args...)
}

type Empty struct{}

func (v Empty) Debug(args ...interface{}) {
}

func (v Empty) Debugf(format string, args ...interface{}) {
}

func (v Empty) Info(args ...interface{}) {
}

func (v Empty) Infof(format string, args ...interface{}) {
}

func (v Empty) Warn(args ...interface{}) {
}

func (v Empty) Warnf(format string, args ...interface{}) {
}

func (v Empty) Error(args ...interface{}) {
}

func (v Empty) Errorf(format string, args ...interface{}) {
}

type Verbose interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}
