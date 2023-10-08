package config

import "github.com/sirupsen/logrus"

type Logger struct {
	*logrus.Entry
}

func NewLogger(name string) *Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logger.SetLevel(logrus.DebugLevel)

	entry := logger.WithFields(logrus.Fields{
		"app": name,
	})

	return &Logger{entry}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Entry.Infof(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Entry.Debugf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Entry.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Entry.Errorf(format, args...)
}
