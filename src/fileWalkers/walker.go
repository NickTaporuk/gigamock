package fileWalkers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	urlrouter "github.com/azer/url-router"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"

	"github.com/NickTaporuk/gigamock/src/fileProvider"
	"github.com/NickTaporuk/gigamock/src/fileType"
)

// DirWalker represents an interface to index and filter files inside particular dir
type DirWalker interface {
	Walk() (map[string]string, error)
	validation.Validatable
}

type DirWalk struct {
	rootDirPath string
	logger      *logrus.Entry
}

// SetRootPath
func (dw *DirWalk) SetRootDirPath(rootDirPath string) {
	dw.rootDirPath = rootDirPath
}

func NewDirWalk(rootDirPath string, lgr *logrus.Entry) *DirWalk {
	return &DirWalk{rootDirPath: rootDirPath, logger: lgr}
}

type IndexedData struct {
	FilePath       string
	ScenarioNumber int
}

type ListIndexedData []IndexedData

func (dw *DirWalk) Walk(router *urlrouter.Router) (map[string]IndexedData, error) {

	err := dw.prepareAbsolutePath()
	if err != nil {
		return nil, err
	}

	err = dw.Validate()
	if err != nil {
		return nil, err
	}

	filesTree := map[string]IndexedData{}

	walkFunk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filePath := path

		ext, err := fileType.FileExtensionDetection(info.Name())
		if err != nil {
			dw.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"action":        "fileType.FileExtensionDetection(info.Name())",
					"method":        "func (dw *DirWalk) Walk(router *urlrouter.Router) (map[string]IndexedData, error)",
					"stack":         string(debug.Stack()),
					"ext":           ext,
					"filePath":      filePath,
					"fileInfo.Name": info.Name(),
				}).Error("action fileType.FileExtensionDetection is retrieved an error")

			return err
		}

		provider, err := fileProvider.Factory(ext, dw.logger)
		if err != nil {
			dw.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"action":        "fileProvider.Factory(ext, dw.logger)",
					"method":        "func (dw *DirWalk) Walk(router *urlrouter.Router) (map[string]IndexedData, error)",
					"stack":         string(debug.Stack()),
					"ext":           ext,
					"filePath":      filePath,
					"fileInfo.Name": info.Name(),
				}).Error("action fileProvider.Factory is retrieved an error")
			return err
		}

		scenario, err := provider.Unmarshal(filePath)
		if err != nil {
			dw.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"action":        "provider.Unmarshal(filePath)",
					"method":        "func (dw *DirWalk) Walk(router *urlrouter.Router) (map[string]IndexedData, error)",
					"stack":         string(debug.Stack()),
					"ext":           ext,
					"filePath":      filePath,
					"fileInfo.Name": info.Name(),
				}).Error("action provider.Unmarshal is retrieved an error")
			return err
		}

		router.Add(scenario.Path)

		filesTree[PrepareInMemoryStoreKey(scenario.Path, scenario.Method)] = IndexedData{FilePath: filePath}

		dw.logger.Info(fmt.Sprintf("file %s for path %s for method %s was indexed", info.Name(), scenario.Path, scenario.Method))

		return nil
	}

	err = filepath.Walk(dw.rootDirPath, walkFunk)
	if err != nil {
		dw.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"action":         "filepath.Walk(dw.rootDirPath, walkFunk)",
				"method":         "func (dw *DirWalk) Walk(router *urlrouter.Router) (map[string]IndexedData, error)",
				"stack":          string(debug.Stack()),
				"root directory": dw.rootDirPath,
			}).Error("action filepath.Walk is retrieved an error")

		return nil, err
	}

	return filesTree, nil
}

// prepareAbsolutePath
func (dw *DirWalk) prepareAbsolutePath() error {
	absPath, err := filepath.Abs(dw.rootDirPath)
	if err != nil {
		return err
	}

	dw.SetRootDirPath(absPath)

	//
	return nil
}

// Validate is fields checker to validate values
func (dw *DirWalk) Validate() error {
	//
	return nil
}
