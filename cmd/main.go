package main

import (
	"log"

	"github.com/NickTaporuk/gigamock/src/app"
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
