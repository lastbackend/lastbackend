package log

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"runtime"
)

type Log struct {
	logger *logrus.Logger
	skip   int
}

func (l *Log) Init() {
	l.logger = logrus.New()
	l.logger.Out = os.Stdout
	l.skip = 2
}

func (l *Log) SetDebugLevel() {
	l.logger.Level = logrus.DebugLevel
}

func (l *Log) Debug(args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Debug(args)
	}
}

func (l *Log) Debugf(format string, args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Debugf(format, args)
	}
}

func (l *Log) Info(args ...interface{}) {
	if l.logger.Level >= logrus.InfoLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Info(args...)
	}
}

func (l *Log) Infof(format string, args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Infof(format, args)
	}
}

func (l *Log) Error(args ...interface{}) {
	if l.logger.Level >= logrus.ErrorLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Error(args...)
	}
}

func (l *Log) Errorf(format string, args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Errorf(format, args)
	}
}

func (l *Log) Fatal(args ...interface{}) {
	if l.logger.Level >= logrus.FatalLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Fatal(args...)
	}
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Fatalf(format, args)
	}
}

func (l *Log) Panic(args ...interface{}) {
	if l.logger.Level >= logrus.PanicLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Panic(args...)
	}
}

func (l *Log) Panicf(format string, args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Panicf(format, args)
	}
}

func (l *Log) Warn(args ...interface{}) {
	if l.logger.Level >= logrus.WarnLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Warn(args...)
	}
}

func (l *Log) Warnf(format string, args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		entry := l.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileLine(l.skip)
		entry.Warnf(format, args)
	}
}

// Logger for GORM
func (l *Log) Print(args ...interface{}) {
	if l.logger.Level >= logrus.ErrorLevel {
		l.logger.Error(args)
	}
}

func fileLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	}
	return fmt.Sprintf("%s:%d", file, line)
}
