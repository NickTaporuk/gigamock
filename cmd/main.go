package main

import (
	"context"
	"fmt"
	"log"

	"github.com/NickTaporuk/gigamock/src/app"
)

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	inst := app.NewApp(ctx)

	err := inst.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Printf("STOP MAIN==>\n")
		cancel()
		err = inst.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
