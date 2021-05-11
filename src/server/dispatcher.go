package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	urlrouter "github.com/azer/url-router"
	"github.com/sirupsen/logrus"

	"github.com/NickTaporuk/gigamock/src/fileProvider"
	"github.com/NickTaporuk/gigamock/src/fileType"
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
	"github.com/NickTaporuk/gigamock/src/handlers/inMemory"
	"github.com/NickTaporuk/gigamock/src/scenarioType"
)

// Dispatcher internally maintains all part of the app
type Dispatcher struct {
	indexedFiles map[string]fileWalkers.IndexedData
	router       *urlrouter.Router
	logger       *logrus.Entry
}

// NewDispatcher is the constructor
func NewDispatcher(
	indexedFiles map[string]fileWalkers.IndexedData,
	router *urlrouter.Router,
	lgr *logrus.Entry,
) *Dispatcher {
	return &Dispatcher{
		indexedFiles: indexedFiles,
		router:       router,
		logger:       lgr,
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

	if v, ok := di.indexedFiles[fileWalkers.PrepareImMemoryStoreKey(match.Pattern, req.Method)]; ok && match != nil {
		di.logger.Debug(
			fmt.Sprintf(
				"route %s for method %s is matched to file path %s, use scenario number %d",
				match.Pattern,
				req.Method,
				v.FilePath,
				v.ScenarioNumber,
			))

		ext, err := fileType.FileExtensionDetection(v.FilePath)
		if err != nil {
			return err
		}
		di.logger.Debug(fmt.Sprintf("file %s extension is %s", v.FilePath, ext))

		provider, err := fileProvider.Factory(ext, di.logger)
		if err != nil {
			return err
		}
		di.logger.Debug(fmt.Sprintf("file provider is %v, extension is %s", provider, ext))

		// should to get type of a scenario
		// can be http, graphql, grpc, kafka and so one
		scenario, err := provider.Unmarshal(v.FilePath)
		if err != nil {
			return err
		}
		di.logger.Debug(fmt.Sprintf("scenario data parsed, scenario data : %v", scenario))

		scenarioTypeProvider, err := scenarioType.Factory(scenario.Type, w, req)
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

func (di *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path == "/favicon.ico" {
		return
	}

	err := di.RouteMatching(w, req)
	if err != nil {
		di.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"trace":   string(debug.Stack()),
				"request": req,
				"method":  "di.RouteMatching",
			}).
			Error("route matching retrieved an error")

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
			di.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"trace":   string(debug.Stack()),
					"address": addr,
				}).
				Error("server retrieved an error")
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
		di.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"trace": string(debug.Stack()),
			}).Fatal("shutting down")
	}

	di.logger.Info("shutting down")
	os.Exit(0)
}
