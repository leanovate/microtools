package logging

import (
	"fmt"
	"github.com/go-errors/errors"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

const (
	Fatal = iota
	Error
	Warn
	Info
	Debug
)

type loggerSimple struct {
	logger *log.Logger
	level  int
	out    io.Writer
}

func NewSimpleLogger(out io.Writer, level int) Logger {
	return &loggerSimple{
		logger: log.New(out, "", log.LstdFlags),
		level:  level,
		out:    out,
	}
}

func NewSimpleLoggerNull() Logger {
	return NewSimpleLogger(ioutil.Discard, Fatal)
}

func (l *loggerSimple) ErrorErr(err error) {
	if l.level >= Error {
		richErr := errors.Wrap(err, 1)

		l.logger.Print(richErr.ErrorStack())
	}
}

func (l *loggerSimple) Errorf(format string, args ...interface{}) {
	if l.level >= Error {
		l.logger.Printf(format, args...)
	}
}

func (l *loggerSimple) Error(args ...interface{}) {
	if l.level >= Error {
		l.logger.Print(args...)
	}
}

func (l *loggerSimple) Warnf(format string, args ...interface{}) {
	if l.level >= Warn {
		l.logger.Printf(format, args...)
	}
}

func (l *loggerSimple) Warn(args ...interface{}) {
	if l.level >= Warn {
		l.logger.Print(args...)
	}
}

func (l *loggerSimple) Infof(format string, args ...interface{}) {
	if l.level >= Info {
		l.logger.Printf(format, args...)
	}
}

func (l *loggerSimple) Info(args ...interface{}) {
	if l.level >= Info {
		l.logger.Print(args...)
	}
}

func (l *loggerSimple) Debugf(format string, args ...interface{}) {
	if l.level >= Debug {
		l.logger.Printf(format, args...)
	}
}

func (l *loggerSimple) Debug(args ...interface{}) {
	if l.level >= Debug {
		l.logger.Print(args...)
	}
}

func (l *loggerSimple) WithContext(fields map[string]interface{}) Logger {
	elements := make([]string, 0)
	elements = append(elements, l.logger.Prefix())
	for k, v := range fields {
		elements = append(elements, fmt.Sprintf("%s=%v", k, v))
	}

	prefix := strings.Join(elements, " ")
	return &loggerSimple{
		logger: log.New(l.out, prefix, log.LstdFlags),
		level:  l.level,
		out:    l.out,
	}
}
