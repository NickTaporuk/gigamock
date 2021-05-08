package fileProvider

import (
	"errors"
	"github.com/NickTaporuk/gigamock/src/fileType"
)

func Factory(ext string) (FileProvider, error) {
	switch ext {
	case fileType.FileExtensionYAML:
		return NewYAMLProvider(), nil
	case fileType.FileExtensionJSON:
		return NewJSONProvider(), nil
	default:
		return nil, errors.New("extension type is not found")
	}
}
