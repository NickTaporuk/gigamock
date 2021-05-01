package fileType

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	// FileExtensionYAML
	FileExtensionYAML = ".yaml"
	// FileExtensionYML
	FileExtensionYML = ".yml"
	// FileExtensionJSON
	FileExtensionJSON = ".json"
)

// FileExtensionDetection
func FileExtensionDetection(fileInfo os.FileInfo) (string, error) {

	ext := filepath.Ext(fileInfo.Name())

	switch ext {
	case FileExtensionYML:
		return ext, nil
	case FileExtensionJSON:
		return ext, nil
	case FileExtensionYAML:
		return ext, nil
	}

	return "", errors.New("extension type is not found")
}
