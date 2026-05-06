package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/sirupsen/logrus"

	urlrouter "github.com/azer/url-router"

	"github.com/NickTaporuk/gigamock/src/fileWalkers"
	"github.com/NickTaporuk/gigamock/src/logger"
	"github.com/NickTaporuk/gigamock/src/server"
)

// Application is a root structure to init that app
type Application interface {
	Run() error
	Stop() error
}

type App struct {
}

// Config contains runtime settings for the mock server.
type Config struct {
	ServerIP          string
	ServerPort        string
	DirPaths          []string
	LoggerLevel       string
	LoggerPrettyPrint bool
}

func NewApp() *App {
	return &App{}
}

func (a App) Stop() error {
	return nil
}

// DefaultConfig returns app settings used when no CLI flags are provided.
func DefaultConfig() (Config, error) {
	path, err := filepath.Abs("./config")
	if err != nil {
		return Config{}, err
	}

	return Config{
		ServerIP:          "0.0.0.0",
		ServerPort:        ":7777",
		DirPaths:          []string{path},
		LoggerLevel:       "DEBUG",
		LoggerPrettyPrint: false,
	}, nil
}

// Run starts the app with default settings.
func (a App) Run() error {
	cfg, err := DefaultConfig()
	if err != nil {
		return err
	}

	return a.RunWithConfig(cfg)
}

// RunWithConfig starts the app with explicit runtime settings.
func (a App) RunWithConfig(cfg Config) error {
	// router is an instance of urlrouter to match urls with parameters
	router := urlrouter.New()
	// init logger
	writers := []io.Writer{os.Stdout}
	localLogger := logger.NewLocalLogger(writers)
	err := localLogger.Init(cfg.LoggerLevel, cfg.LoggerPrettyPrint)
	if err != nil {
		return err
	}
	lgr := localLogger.Logger()

	files, err := a.walkConfigDirs(cfg.DirPaths, router, lgr)
	if err != nil {
		lgr.
			WithError(err).
			WithFields(logrus.Fields{
				"trace":  string(debug.Stack()),
				"router": router,
				"method": "func (a App) Run() error",
				"action": "filesWalker.Walk(router)",
			}).
			Error("file walker retrieved an error")
		return err
	}

	di := server.NewDispatcher(files, router, lgr)

	di.Start(cfg.ServerIP + cfg.ServerPort)

	return nil
}

func (a App) walkConfigDirs(
	dirPaths []string,
	router *urlrouter.Router,
	lgr *logrus.Entry,
) (map[string]fileWalkers.IndexedData, error) {
	if len(dirPaths) == 0 {
		return nil, fmt.Errorf("at least one dir-path is required")
	}

	files := map[string]fileWalkers.IndexedData{}
	for _, dirPath := range dirPaths {
		filesWalker := fileWalkers.NewDirWalk(dirPath, lgr)

		dirFiles, err := filesWalker.Walk(router)
		if err != nil {
			return nil, err
		}

		for key, indexedData := range dirFiles {
			if existingData, ok := files[key]; ok {
				return nil, fmt.Errorf(
					"duplicate mock endpoint %s found in %s and %s",
					key,
					existingData.FilePath,
					indexedData.FilePath,
				)
			}

			files[key] = indexedData
		}
	}

	return files, nil
}
