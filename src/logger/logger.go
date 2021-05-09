package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

// ILocalLogger is interface for inheritance of structure
type ILocalLogger interface {
	Writers() []io.Writer
	Logger() *logrus.Entry
	SetLogger(logger *logrus.Entry)
	Init(level string, prettyPrint bool) error
}
