package webot

import (
	"io"
	"log"
	"os"
)

type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
}

// NewLogger create a Logger wraps the *log.Logger
func NewLogger(output io.Writer, prefix string, flag int) Logger {
	return &logger{l: log.New(output, prefix, flag)}
}

func createDefaultLogger() Logger {
	return NewLogger(os.Stdout, "", log.Ldate|log.Lmicroseconds)
}

type logger struct {
	l *log.Logger
}

func (l *logger) Errorf(format string, v ...interface{}) {}

func (l *logger) Warnf(format string, v ...interface{}) {}

func (l *logger) Debugf(format string, v ...interface{}) {}

func (l *logger) Infof(format string, v ...interface{}) {}
