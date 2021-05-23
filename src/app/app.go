package app

import (
	"flag"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"

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

func NewApp() *App {
	return &App{}
}

func (a App) Stop() error {
	return nil
}

// Run
func (a App) Run() error {
	path, err := filepath.Abs("./config")
	if err != nil {
		return err
	}

	serverIP := flag.String("server-ip", "0.0.0.0", "Definition server IP")
	serverPort := flag.String("server-port", ":7777", "Definition server Port")
	dirPath := flag.String("dir-path", path, "Mocks config folder")
	loggerLevel := flag.String("logger-level", "DEBUG", "logger level")
	loggerPrettyPrint := flag.Bool("logger-pretty-print", false, "logger level")
	flag.Parse()

	// router is an instance of urlrouter to match urls with parameters
	router := urlrouter.New()
	// init logger
	writers := []io.Writer{os.Stdout}
	localLogger := logger.NewLocalLogger(writers)
	err = localLogger.Init(*loggerLevel, *loggerPrettyPrint)
	if err != nil {
		return err
	}
	lgr := localLogger.Logger()
	//
	filesWalker := fileWalkers.NewDirWalk(*dirPath, lgr)

	files, err := filesWalker.Walk(router)
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

	di.Start(*serverIP + *serverPort)

	return nil
}
