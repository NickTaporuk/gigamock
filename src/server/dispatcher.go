package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
	_ "net/http/pprof"

	"github.com/NickTaporuk/gigamock/src/fileProvider"
	"github.com/NickTaporuk/gigamock/src/fileType"
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
	"github.com/NickTaporuk/gigamock/src/handlers/inMemory"
	"github.com/NickTaporuk/gigamock/src/scenarioType"
	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/NickTaporuk/gigamock/src/webhookType"
	urlrouter "github.com/azer/url-router"
)

// Dispatcher internally maintains all part of the app
type Dispatcher struct {
	ctx          context.Context
	indexedFiles map[string]fileWalkers.IndexedData
	router       *urlrouter.Router
	logger       *logrus.Entry
	usePPROF     bool
}

// NewDispatcher is the constructor
func NewDispatcher(
	ctx context.Context,
	indexedFiles map[string]fileWalkers.IndexedData,
	router *urlrouter.Router,
	lgr *logrus.Entry,
	usePPROF bool,
) *Dispatcher {
	return &Dispatcher{
		indexedFiles: indexedFiles,
		router:       router,
		logger:       lgr,
		ctx:          ctx,
		usePPROF:     usePPROF,
	}
}

// inMemoryHandlers is dispatched all request to show
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

//pprof debug
func (di Dispatcher) pprrofHandling(w http.ResponseWriter, req *http.Request) (bool, error) {
	switch req.URL.Path {
	case "/debug/pprof/":
		{
			pprof.Index(w, req)
			return true, nil
		}
	case "/debug/pprof/cmdline":
		{
			pprof.Cmdline(w, req)
			return true, nil
		}
	case "/debug/pprof/profile":
		{
			pprof.Profile(w, req)
			return true, nil
		}
	case "/debug/pprof/symbol":
		{
			pprof.Symbol(w, req)
			return true, nil
		}
	case "/debug/pprof/trace":
		{
			pprof.Trace(w, req)
			return true, nil
		}

	}
	return false, nil
}

// RouteMatching
func (di *Dispatcher) RouteMatching(w http.ResponseWriter, req *http.Request) error {
	if di.usePPROF {
		matched, err := di.pprrofHandling(w, req)
		if err != nil {
			return err
		}

		if matched {
			return nil
		}
	}

	matched, err := di.inMemoryHandlers(w, req)
	if err != nil {
		return err
	}

	if matched {
		return nil
	}

	match := di.router.Match(req.URL.Path)

	if v, ok := di.indexedFiles[fileWalkers.PrepareInMemoryStoreKey(match.Pattern, req.Method)]; ok && match != nil {
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
			di.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"indexed file": v,
					"method":       "dispatcher.RouteMatching",
					"action":       "fileType.FileExtensionDetection(v.FilePath)",
					"stack":        string(debug.Stack()),
				}).Error("fileType.FileExtensionDetection is retrieved an error")
			return err
		}
		di.logger.Debug(fmt.Sprintf("file %s extension is %s", v.FilePath, ext))

		provider, err := fileProvider.Factory(ext, di.logger)
		if err != nil {
			di.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"extension": ext,
					"provider":  provider,
					"method":    "dispatcher.RouteMatching",
					"action":    "fileProvider.Factory(ext, di.logger)",
					"stack":     string(debug.Stack()),
				}).Error("fileProvider.Factory is retrieved an error")
			return err
		}
		di.logger.Debug(fmt.Sprintf("file provider is %v, extension is %s", provider, ext))

		// should to get type of a scenario
		// can be http, graphql, grpc, kafka and so one
		scenario, err := provider.Unmarshal(v.FilePath)
		if err != nil {
			di.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"scenario": scenario,
					"method":   "dispatcher.RouteMatching",
					"action":   "provider.Unmarshal(v.FilePath)",
					"stack":    string(debug.Stack()),
				}).Error("provider.Unmarshal is retrieved an error")
			return err
		}
		di.logger.Debug(fmt.Sprintf("scenario data parsed, scenario data : %v", scenario))

		scenarioTypeProvider, err := scenarioType.Factory(scenario.Type, w, req, di.logger, di.ctx)
		if err != nil {
			di.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"scenario": scenario,
					"method":   "dispatcher.RouteMatching",
					"action":   "scenarioType.Factory(scenario.Type, w, req)",
					"stack":    string(debug.Stack()),
				}).Error("scenarioType.Factory is retrieved an error")
			return err
		}

		err = scenarioTypeProvider.Unmarshal(scenario.Scenarios)
		if err != nil {
			di.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"scenario": scenario,
					"method":   "dispatcher.RouteMatching",
					"action":   "scenarioTypeProvider.Unmarshal(scenario.Scenarios)",
					"stack":    string(debug.Stack()),
				}).Error("scenarioTypeProvider.Unmarshal is retrieved an error")

			return err
		}

		err = scenarioTypeProvider.Validate()
		if err != nil {
			di.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"scenario": scenario,
					"method":   "dispatcher.RouteMatching",
					"action":   "scenarioTypeProvider.Validate()",
					"stack":    string(debug.Stack()),
				}).Error("scenarioTypeProvider.Validate is retrieved an error")
			return err
		}

		scenarioTypeProvider.Retrieve(v.ScenarioNumber)

		// need to run webhook
		if scenario.WebHook != nil {
			go func(
				scenario *scenarios.BaseGigaMockScenario,
				scenarioNumber int,
			) {
				err = di.runWebHook(scenario, scenarioNumber)
				if err != nil {
					di.logger.
						WithError(err).
						WithFields(logrus.Fields{
							"scenario": scenario,
							"method":   "dispatcher.RouteMatching",
							"action":   "di.runWebHook(scenario)",
							"stack":    string(debug.Stack()),
						}).Error("di.runWebHook is retrieved an error")
					return
				}
			}(scenario, v.ScenarioNumber)
		}
	} else {
		//	no pattern matched; send 404 response
		http.NotFound(w, req)
	}

	return nil
}

func (di *Dispatcher) runWebHook(
	scenario *scenarios.BaseGigaMockScenario,
	scenarioNumber int,
) error {
	err := scenario.WebHook.Validate()
	if err != nil {
		di.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"scenario": scenario,
				"method":   "dispatcher.RouteMatching",
				"action":   "scenario.WebHook.Validate()",
				"stack":    string(debug.Stack()),
			}).Error("scenario.WebHook.Validate is retrieved an error")
		return err
	}

	webHookProvider, err := webhookType.Factory(scenario.WebHook, di.logger, scenarioNumber)
	if err != nil {
		di.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"scenario": scenario,
				"method":   "dispatcher.RouteMatching",
				"action":   "webhookType.Factory(scenario.WebHook.Type, di.logger)",
				"stack":    string(debug.Stack()),
			}).Error("webhookType.Factory is retrieved an error")
		return err
	}

	err = webHookProvider.Unmarshal(scenario.WebHook.Scenarios)
	if err != nil {
		di.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"scenario": scenario,
				"method":   "dispatcher.RouteMatching",
				"action":   "webHookProvider.Unmarshal()",
				"stack":    string(debug.Stack()),
			}).Error("webHookProvider.Unmarshal is retrieved an error")
		return err
	}

	err = webHookProvider.Validate()
	if err != nil {
		di.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"scenario": scenario,
				"method":   "dispatcher.RouteMatching",
				"action":   "webHookProvider.Validate()",
				"stack":    string(debug.Stack()),
			}).Error("webHookProvider.Validate() is retrieved an error")
		return err
	}

	err = webHookProvider.Send()
	if err != nil {
		di.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"scenario": scenario,
				"method":   "dispatcher.RouteMatching",
				"action":   "webHookProvider.Send()",
				"stack":    string(debug.Stack()),
			}).Error("webHookProvider.Send is retrieved an error")
		return err
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
				"action":  "di.RouteMatching",
				"method":  "func (di *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request)",
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
					"trace":  string(debug.Stack()),
					"srv":    srv,
					"action": "srv.ListenAndServe()",
					"method": "func (di Dispatcher) Start(addr string)",
					"stack":  string(debug.Stack()),
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
				"trace":  string(debug.Stack()),
				"action": "srv.Shutdown(ctx)",
				"method": "func (di Dispatcher) Start(addr string)",
			}).Fatal("shutting down")
	}

	di.logger.Info("shutting down")
	os.Exit(0)
}
