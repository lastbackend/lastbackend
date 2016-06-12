package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/deployithq/deployit/utils"
)

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
		entry.Data["file"] = utils.FileLine()
		entry.Debug(args)
	}
}

func (l *Log) Debugf(format string, args ...interface{}) {
	if l.Logger.Level >= logrus.DebugLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = utils.FileLine()
		entry.Debugf(format, args)
	}
}

func (l *Log) Info(args ...interface{}) {
	if l.Logger.Level >= logrus.InfoLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = utils.FileLine()
		entry.Info(args...)
	}
}

func (l *Log) Infof(format string, args ...interface{}) {
	if l.Logger.Level >= logrus.InfoLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = utils.FileLine()
		entry.Infof(format, args...)
	}
}

func (l *Log) Error(args ...interface{}) {
	if l.Logger.Level >= logrus.ErrorLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = utils.FileLine()
		entry.Error(args...)
	}
}

func (l *Log) Errorf(format string, args ...interface{}) {
	if l.Logger.Level >= logrus.DebugLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = utils.FileLine()
		entry.Errorf(format, args)
	}
}

func (l *Log) Fatal(args ...interface{}) {
	if l.Logger.Level >= logrus.FatalLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = utils.FileLine()
		entry.Fatal(args...)
	}
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	if l.Logger.Level >= logrus.DebugLevel {
		entry := l.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = utils.FileLine()
		entry.Fatalf(format, args)
	}
}
