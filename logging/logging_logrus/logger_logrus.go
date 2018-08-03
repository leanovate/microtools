package logging_logrus

import (
	"fmt"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logger *logrus.Logger
	fields map[string]interface{}
}

func init() {
	logging.RegisterBackend("logrus", NewLogrusLogger)
}

// NewLogrusLogger creates a Looger implementation based on the famous
// logrus package
func NewLogrusLogger(options logging.Options) logging.Logger {
	backend := logrus.New()
	backend.Out = options.GetOutput()

	switch options.LogFormat {
	case "json":
		backend.Formatter = &logrus.JSONFormatter{}
	default:
		backend.Formatter = &logrus.TextFormatter{}
	}
	switch options.Level {
	case logging.Fatal:
		backend.Level = logrus.FatalLevel
	case logging.Error:
		backend.Level = logrus.ErrorLevel
	case logging.Warn:
		backend.Level = logrus.WarnLevel
	case logging.Info:
		backend.Level = logrus.InfoLevel
	case logging.Debug:
		backend.Level = logrus.DebugLevel
	}

	return &logrusLogger{
		logger: backend,
	}
}

func (l *logrusLogger) ErrorErr(err error) {
	if l.logger.Level >= logrus.ErrorLevel {
		switch richErr := err.(type) {
		case fmt.Formatter:
			l.logger.WithFields(l.fields).Errorf("%+v", richErr)
		case logging.SimpleStackTracer:
			l.logger.WithFields(l.fields).Errorf(richErr.ErrorStack())
		default:
			wrapped := errors.Wrap(err, err.Error())

			l.logger.WithFields(l.fields).Errorf("%+v", wrapped)
		}
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

func (l *logrusLogger) WithContext(fields map[string]interface{}) logging.Logger {
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

func (l *logrusLogger) WithField(name, value string) logging.Logger {
	newFields := make(map[string]interface{}, 0)
	for key, value := range l.fields {
		newFields[key] = value
	}
	newFields[name] = value

	return &logrusLogger{
		logger: l.logger,
		fields: newFields,
	}
}
