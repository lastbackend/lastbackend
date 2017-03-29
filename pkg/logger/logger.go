package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
	"log/syslog"
	"os"
	"path"
	"runtime"
)

type Logger struct {
	log *logrus.Entry
}

func Init() *Logger {
	l := new(Logger)
	entry := logrus.NewEntry(logrus.New())
	entry.Logger.Out = os.Stdout
	entry.Logger.Formatter = getJSONFormatter()
	l.log = entry.WithFields(getFileLocate(logrus.Fields{}))
	return l
}

func (l *Logger) SetDebugLevel() {
	l.log.Level = logrus.DebugLevel
	l.log.Logger.Formatter = getTextFormatter()
}

func (l *Logger) SetSyslog(network, raddr string, priority syslog.Priority, tag string) {
	if hook, err := logrus_syslog.NewSyslogHook(network, raddr, priority, tag); err == nil {
		l.log.Logger.Hooks.Add(hook)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args)
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
