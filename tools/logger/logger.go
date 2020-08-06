package logger

import (
	"context"
)

// A global variable so that log functions can be directly accessed
var log Logger

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	//Debug has verbose message
	Debug = "debug"
	//Info is default log level
	Info = "info"
	//Warn is for logging messages about possible issues
	Warn = "warn"
	//Error is for logging errors
	Error = "error"
	//Fatal is for logging fatal messages. The system shutdown after logging the message.
	Fatal = "fatal"

	InstanceZapLogger int = iota
)

//Logger is our contract for the logger
type Logger interface {
	Debug(args ...interface{})

	Debugf(format string, args ...interface{})

	Info(args ...interface{})

	Infof(format string, args ...interface{})

	Warn(args ...interface{})

	Warnf(format string, args ...interface{})

	Error(args ...interface{})

	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})

	Fatalf(format string, args ...interface{})

	Panic(args ...interface{})

	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

//NewLogger returns an instance of logger
func NewLogger(config Configuration) error {
	logger, err := newZapLogger(config)
	if err != nil {
		return err
	}
	log = logger
	return nil
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

func NewContext(ctx context.Context, fields Fields) context.Context {
	if fields == nil {
		return context.WithValue(ctx, InstanceZapLogger, WithContext(ctx))
	}
	return context.WithValue(ctx, InstanceZapLogger, WithContext(ctx).WithFields(fields))
}

func WithContext(ctx context.Context) Logger {

	if ctx == nil {
		return log
	}

	if ctxLogger, ok := ctx.Value(InstanceZapLogger).(Logger); ok {
		return ctxLogger
	}

	return log
}
