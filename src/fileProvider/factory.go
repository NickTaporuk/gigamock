package fileProvider

import (
	"errors"
	"github.com/NickTaporuk/gigamock/src/fileType"
	"github.com/sirupsen/logrus"
)

func Factory(ext string, lgr *logrus.Entry) (FileProvider, error) {
	switch ext {
	case fileType.FileExtensionYAML:
		return NewYAMLProvider(lgr), nil
	case fileType.FileExtensionJSON:
		return NewJSONProvider(lgr), nil
	default:
		return nil, errors.New("extension type is not found")
	}
}
