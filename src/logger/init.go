package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

const (
	// LoggerKey
	LoggerKey = "service"
	// LoggerValue
	LoggerValue = "gigamock"
)

// LocalLogger is wrapper of logrus logger
type LocalLogger struct {
	logger  *logrus.Entry
	writers []io.Writer
}

func NewLocalLogger(writers []io.Writer) *LocalLogger {
	return &LocalLogger{writers: writers}
}

// Writers is getter of io writers slice
func (l *LocalLogger) Writers() []io.Writer {
	return l.writers
}

// Logger is getter of logrus logger
func (l *LocalLogger) Logger() *logrus.Entry {
	return l.logger
}

// SetLogger is setter of logrus logger
func (l *LocalLogger) SetLogger(logger *logrus.Entry) {
	l.logger = logger
}

// Init is initializer of local logger
func (l *LocalLogger) Init(level string, prettyPrint bool) error {
	// formatter use for see like a list some log data
	formatter := &logrus.JSONFormatter{PrettyPrint: prettyPrint}

	writers := l.Writers()
	mw := io.MultiWriter(writers...)

	lgr := logrus.New()
	lgr.SetFormatter(formatter)
	lgr.SetOutput(mw)
	// Only log the warning severity or above.
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	lgr.SetLevel(lvl)

	builder := lgr.WithField(LoggerKey, LoggerValue)
	l.SetLogger(builder)

	return nil
}
