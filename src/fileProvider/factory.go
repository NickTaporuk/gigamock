package fileProvider

import (
	"errors"
	"github.com/NickTaporuk/gigamock/src/fileType"
)

func Factory(ext string) (FileProvider, error) {
	if ext == fileType.FileExtensionYAML {
		return NewYAMLProvider(), nil
	}

	return nil, errors.New("extension type is not found")
}
