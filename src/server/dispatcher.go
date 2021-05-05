package server

import (
	"context"
	"fmt"
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
	"github.com/NickTaporuk/gigamock/src/handlers/inMemory"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	urlrouter "github.com/azer/url-router"

	"github.com/NickTaporuk/gigamock/src/fileProvider"
	"github.com/NickTaporuk/gigamock/src/fileType"
	"github.com/NickTaporuk/gigamock/src/store"
)

// Dispatcher
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
		scenario, err := provider.Parse(v.FilePath)
		if err != nil {
			return err
		}

		var body string
		var statusCode int
		if len(scenario.Scenarios) > 0 {
			body = scenario.Scenarios[v.ScenarioNumber].Response.Body
			if scenario.Scenarios[v.ScenarioNumber].Response.StatusCode > 0 {
				statusCode = scenario.Scenarios[v.ScenarioNumber].Response.StatusCode
				if len(scenario.Scenarios[v.ScenarioNumber].Response.Headers) > 0 {
					for headerName, headerValue := range scenario.Scenarios[v.ScenarioNumber].Response.Headers {
						w.Header().Add(headerName, headerValue)
						fmt.Printf("HEADERS==>", v)
					}
				}
			} else {
				statusCode = http.StatusOK
			}
		}

		w.WriteHeader(statusCode)
		w.Write([]byte(body))
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
