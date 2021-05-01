package fileWalkers

import (
	"os"
	"path/filepath"
)

// DirWalker represents an interface to index and filter files inside particular dir
type DirWalker interface {
	Walk() (map[string]string, error)
	Validate() error
}

type DirWalk struct {
	rootDirPath string
}

// SetRootPath
func (dw *DirWalk) SetRootDirPath(rootDirPath string) {
	dw.rootDirPath = rootDirPath
}

func NewDirWalk(rootDirPath string) *DirWalk {
	return &DirWalk{rootDirPath: rootDirPath}
}

func (dw *DirWalk) Walk() (map[string]string, error) {

	err := dw.Validate()
	if err != nil {
		return nil, err
	}

	filesTree := map[string]string{}
	walkFunk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filePath := path + info.Name()
		filesTree[info.Name()] = filePath

		ext, err := FileExtensionDetection(info)
		if err != nil {
			return nil
		}

		provider := fileProvi
		return nil
	}

	err = filepath.Walk(dw.rootDirPath, walkFunk)
	if err != nil {
		return nil, err
	}

	return filesTree, nil
}

// Validate
func (dw *DirWalk) Validate() error {
	absPath, err := filepath.Abs(dw.rootDirPath)
	if err != nil {
		return err
	}

	dw.SetRootDirPath(absPath)

	return nil
}
