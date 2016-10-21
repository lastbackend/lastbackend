package log

import (
	"fmt"
	"github.com/deployithq/deployit/libs/log/color"
	"runtime"
	"time"
)

const _NEWLINE = "\n"

type Log struct {
	skip  int
	debug bool
}

func (l *Log) Init() {
	l.skip = 2
}

func (l *Log) SetDebugLevel() {
	l.debug = true
}

func (l *Log) Debug(args ...interface{}) {
	if l.debug {
		l.print(color.White(l.sprintlnn(args...)))
	}
}

func (l *Log) Debugf(format string, args ...interface{}) {
	if l.debug {
		l.printf(color.White(format), args...)
	}
}

func (l *Log) Info(args ...interface{}) {
	l.print(color.Yellow(l.sprintlnn(args...)))
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.printf(color.Yellow(format), args...)
}

func (l *Log) Error(args ...interface{}) {
	l.print(color.Red(l.sprintlnn(args...)))
}

func (l *Log) Errorf(format string, args ...interface{}) {
	l.printf(color.Red(format), args...)
}

func (l *Log) Fatal(args ...interface{}) {
	l.print(color.Red(l.sprintlnn(args...)))
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	l.printf(color.Red(format), args...)
}

func (l *Log) Panic(args ...interface{}) {
	l.print(color.Red(l.sprintlnn(args...)))
}

func (l *Log) Panicf(format string, args ...interface{}) {
	l.printf(color.Red(format), args...)
}

func (l *Log) Warn(args ...interface{}) {
	l.print(color.Magenta(l.sprintlnn(args...)))
}

func (l *Log) Warnf(format string, args ...interface{}) {
	l.printf(color.Magenta(format), args...)
}

func (l *Log) print(message string) {
	fmt.Printf("%v %s", time.Now().Format("2006-01-02 15:04:05"), message+_NEWLINE)
}

func (l *Log) printf(format string, a ...interface{}) {
	fmt.Printf("%v %s", time.Now().Format("2006-01-02 15:04:05"), color.Cyan(fileLine(l.skip)))
	fmt.Printf(format+_NEWLINE, a...)
}

func fileLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (l *Log) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
