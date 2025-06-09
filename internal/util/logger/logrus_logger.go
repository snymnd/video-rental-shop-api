package logger

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	log *logrus.Logger
}

func SetLogrusLogger() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		ForceColors:     true,
	})
	log.SetOutput(os.Stdout)

	logrusLogger := &logrusLogger{log: log}

	SetLogger(logrusLogger)
}

func (l *logrusLogger) GetWriter() io.Writer {
	return l.log.Out
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logrusLogger) WithFields(args map[string]interface{}) Logger {
	return &logrusEntry{
		entry: l.log.WithFields(args),
	}
}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}
