package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	urlrouter "github.com/azer/url-router"

	"github.com/NickTaporuk/gigamock/src/fileProvider"
	"github.com/NickTaporuk/gigamock/src/fileType"
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
	"github.com/NickTaporuk/gigamock/src/handlers/inMemory"
	"github.com/NickTaporuk/gigamock/src/scenarioType"
	"github.com/NickTaporuk/gigamock/src/store"
)

// Dispatcher internally maintains all part of the app
type Dispatcher struct {
	indexedFiles map[string]fileWalkers.IndexedData
	router       *urlrouter.Router
}

// NewDispatcher is the constructor
func NewDispatcher(
	indexedFiles map[string]fileWalkers.IndexedData,
	router *urlrouter.Router,
	inMemoryStore *store.InMemoryStore,
) *Dispatcher {
	return &Dispatcher{
		indexedFiles: indexedFiles,
		router:       router,
	}
}

// inMemoryHandlers
func (di *Dispatcher) inMemoryHandlers(w http.ResponseWriter, req *http.Request) (bool, error) {

	if req.URL.Path == "/internal/v1/in-memory" {
		h := inMemory.NewHandler(&di.indexedFiles)
		switch req.Method {
		case http.MethodPost:
			h.AddRecord(w, req)
			return true, nil
		case http.MethodGet:
			h.List(w, req)
			return true, nil
		}
	}

	return false, nil
}

// RouteMatching
func (di *Dispatcher) RouteMatching(w http.ResponseWriter, req *http.Request) error {
	matched, err := di.inMemoryHandlers(w, req)
	if err != nil {
		return err
	}

	if matched {
		return nil
	}

	match := di.router.Match(req.URL.Path)

	if v, ok := di.indexedFiles[match.Pattern+"|"+req.Method]; ok && match != nil {
		fmt.Println("MATCH PATTERN VALUE=>", v, "PATTERN=>", match.Pattern)

		ext, err := fileType.FileExtensionDetection(v.FilePath)
		if err != nil {
			return err
		}

		provider, err := fileProvider.Factory(ext)
		if err != nil {
			return err
		}
		// should to get type of a scenario
		// can be http, graphql, grpc, kafka and so one
		scenario, err := provider.Unmarshal(v.FilePath)
		if err != nil {
			return err
		}

		println(scenario.Type)
		scenarioTypeProvider, err := scenarioType.Factory(scenario.Type, w)
		if err != nil {
			return err
		}

		err = scenarioTypeProvider.Unmarshal(scenario.Scenarios)
		if err != nil {
			return err
		}

		scenarioTypeProvider.Retrieve(v.ScenarioNumber)
	} else {
		//	no pattern matched; send 404 response
		http.NotFound(w, req)
	}

	return nil
}

func (di *Dispatcher) TypeMatching(w http.ResponseWriter, req *http.Request) {

}

func (di *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path == "/favicon.ico" {
		return
	}
	// TODO: should check path
	// TODO: check type can be http or graphql, grpc or kafka
	err := di.RouteMatching(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Start initialize the HTTP mock server
func (di Dispatcher) Start(addr string) {
	var wait time.Duration

	srv := &http.Server{
		Addr:         addr,
		Handler:      &di,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Println("shutting down")
		log.Fatal(err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
