package logging

import (
	"os"
	"path"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"github.com/go-errors/errors"
	"gopkg.in/codegangsta/cli.v2"
)

type logrusLogger struct {
	logger *logrus.Logger
	fields map[string]interface{}
}

func NewLogrusLogger(ctx *cli.Context) Logger {
	logFile := ctx.String("log-file")
	if logFile != "" {
		if err := os.MkdirAll(path.Dir(logFile), 0755); err != nil {
			logrus.Errorf("Failed to create path %s: %s", path.Dir(logFile), err.Error())
		} else {
			file, err := os.OpenFile(logFile, syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, 0644)
			if err != nil {
				logrus.Errorf("Failed to open log file %s: %s", logFile, err.Error())
			} else {
				logrus.SetOutput(file)
			}
		}
	}
	switch ctx.String("log-format") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "logstash":
		logrus.SetFormatter(&logstash.LogstashFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	if ctx.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return &logrusLogger{
		logger: logrus.StandardLogger(),
	}
}

func (l *logrusLogger) ErrorErr(err error) {
	if l.logger.Level >= logrus.ErrorLevel {
		richErr := errors.Wrap(err, 1)
		l.logger.WithFields(l.fields).Errorf(richErr.ErrorStack())
	}
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	if l.logger.Level >= logrus.ErrorLevel {
		l.logger.WithFields(l.fields).Errorf(format, args...)
	}
}

func (l *logrusLogger) Error(args ...interface{}) {
	if l.logger.Level >= logrus.ErrorLevel {
		l.logger.WithFields(l.fields).Error(args...)
	}
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	if l.logger.Level >= logrus.WarnLevel {
		l.logger.WithFields(l.fields).Warnf(format, args...)
	}
}

func (l *logrusLogger) Warn(args ...interface{}) {
	if l.logger.Level >= logrus.WarnLevel {
		l.logger.WithFields(l.fields).Warn(args...)
	}
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	if l.logger.Level >= logrus.InfoLevel {
		l.logger.WithFields(l.fields).Infof(format, args...)
	}
}

func (l *logrusLogger) Info(args ...interface{}) {
	if l.logger.Level >= logrus.InfoLevel {
		l.logger.WithFields(l.fields).Info(args...)
	}
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		l.logger.WithFields(l.fields).Debugf(format, args...)
	}

}

func (l *logrusLogger) Debug(args ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		l.logger.WithFields(l.fields).Debug(args...)
	}
}

func (l *logrusLogger) WithContext(fields map[string]interface{}) Logger {
	newFields := make(map[string]interface{}, 0)
	for key, value := range l.fields {
		newFields[key] = value
	}
	for key, value := range fields {
		newFields[key] = value
	}

	return &logrusLogger{
		logger: l.logger,
		fields: newFields,
	}
}
