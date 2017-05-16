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

// Verbose is a boolean type that implements Infof (like Printf) etc.
type Verbose bool

func (v Verbose) Debug(args ...interface{}) {
	if v {
		l.log.Debug(args...)
	}
}

func (v Verbose) Debugf(format string, args ...interface{}) {
	if v {
		l.log.Debugf(format, args...)
	}
}

func (v Verbose) Info(args ...interface{}) {
	if v {
		l.log.Info(args...)
	}
}

func (v Verbose) Infof(format string, args ...interface{}) {
	if v {
		l.log.Infof(format, args...)
	}
}

func (v Verbose) Warn(args ...interface{}) {
	if v {
		l.log.Warn(args...)
	}
}

func (v Verbose) Warnf(format string, args ...interface{}) {
	if v {
		l.log.Warnf(format, args...)
	}
}

func (v Verbose) Error(args ...interface{}) {
	if v {
		l.log.Error(args...)
	}
}

func (v Verbose) Errorf(format string, args ...interface{}) {
	if v {
		l.log.Errorf(format, args...)
	}
}

func (v Verbose) Fatal(args ...interface{}) {
	if v {
		l.log.Fatal(args...)
	}
}

func (v Verbose) Fatalf(format string, args ...interface{}) {
	if v {
		l.log.Fatalf(format, args...)
	}
}

func (v Verbose) Panic(args ...interface{}) {
	if v {
		l.log.Panic(args...)
	}
}

func (v Verbose) Panicf(format string, args ...interface{}) {
	if v {
		l.log.Panicf(format, args...)
	}
}
