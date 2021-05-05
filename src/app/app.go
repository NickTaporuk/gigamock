package app

import (
	"flag"
	"fmt"
	"github.com/NickTaporuk/gigamock/src/server"
	"github.com/NickTaporuk/gigamock/src/store"
	urlrouter "github.com/azer/url-router"
	"path/filepath"

	"github.com/NickTaporuk/gigamock/src/fileWalkers"
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
	flag.Parse()

	fmt.Println(serverIP, serverPort, dirPath, loggerLevel)

	// router is an instance of urlrouter to match urls with parameters
	router := urlrouter.New()
	//
	inMemoryData := make(map[string]int, 0)
	inMemoryStore := store.NewInMemoryStore(inMemoryData)
	filesWalker := fileWalkers.NewDirWalk(*dirPath)

	files, err := filesWalker.Walk(router)
	if err != nil {
		return err
	}

	fmt.Printf("FILES ==> %#v\n", files)

	di := server.NewDispatcher(files, router, inMemoryStore)

	di.Start(*serverIP+*serverPort)

	return nil
}
