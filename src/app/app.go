package app

import (
	"flag"
	"fmt"
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

// server Ip
// server port
// scripts folder
// config file path
// root path for directory
func (a App) Run() error {
	path, err := filepath.Abs("./config")
	if err != nil {
		return err
	}

	serverIP := flag.String("server-ip", "0.0.0.0", "Definition server IP")
	serverPort := flag.Int("server-port", 7777, "Definition server Port")
	dirPath := flag.String("dir-path", path, "Mocks config folder")
	flag.Parse()

	fmt.Println(serverIP, serverPort, dirPath)

	filesWalker := fileWalkers.NewDirWalk(*dirPath)

	files, err := filesWalker.Walk()
	if err != nil {
		return err
	}

	fmt.Println("FILES ==>", files)

	return nil
}
