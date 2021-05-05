package fileType

import (
	"errors"
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
func FileExtensionDetection(filePath string) (string, error) {

	ext := filepath.Ext(filePath)

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
