package log

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"strings"
)

const bucket = "map"

type Log struct {
	Logger *logrus.Logger
}

func New() *logrus.Logger {
	return logrus.New()
}

func (l *Log) SetDebugLevel() {
	l.Logger.Level = logrus.DebugLevel
}

func (l *Log) Debug(args ...interface{}) {
	if l.Logger.Level >= logrus.DebugLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args)
	}
}

func (l *Log) Debugf(format string, args ...interface{}) {
	if l.Logger.Level >= logrus.DebugLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debugf(format, args)
	}
}

func (l *Log) Error(args ...interface{}) {
	if l.Logger.Level >= logrus.ErrorLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

func (l *Log) Errorf(format string, args ...interface{}) {
	if l.Logger.Level >= logrus.DebugLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Errorf(format, args)
	}
}

func (l *Log) Fatal(args ...interface{}) {
	if l.Logger.Level >= logrus.FatalLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	if l.Logger.Level >= logrus.DebugLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatalf(format, args)
	}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
