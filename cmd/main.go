package main

import (
	"github.com/NickTaporuk/gigamock/src/app"
	"log"
)

func main() {
	inst := app.NewApp()

	err := inst.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = inst.Stop()
		log.Fatal(err)
	}()
}
