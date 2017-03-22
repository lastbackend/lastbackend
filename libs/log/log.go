package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
	"log/syslog"
	"os"
	"path"
	"runtime"
)

type log struct {
	log *logrus.Entry
}

func Init() *log {
	l := new(log)
	entry := logrus.NewEntry(logrus.New())
	entry.Logger.Out = os.Stdout
	entry.Logger.Formatter = getJSONFormatter()
	l.log = entry.WithFields(getFileLocate(logrus.Fields{}))
	return l
}

func (l *log) SetDebugLevel() {
	l.log.Level = logrus.DebugLevel
	l.log.Logger.Formatter = getTextFormatter()
}

func (l *log) SetSyslog(network, raddr string, priority syslog.Priority, tag string) {
	if hook, err := logrus_syslog.NewSyslogHook(network, raddr, priority, tag); err == nil {
		l.log.Logger.Hooks.Add(hook)
	}
}

func (l *log) Debug(args ...interface{}) {
	l.log.Debug(args)
}

func (l *log) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args)
}

func (l *log) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *log) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args)
}

func (l *log) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *log) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args)
}

func (l *log) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *log) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args)
}

func (l *log) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l *log) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args)
}

func (l *log) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *log) Warnf(format string, args ...interface{}) {
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
