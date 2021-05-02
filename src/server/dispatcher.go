package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ryanuber/go-glob"
)

type Dispatcher struct {
	indexedFiles map[string]string
}

func NewDispatcher(indexedFiles map[string]string) *Dispatcher {
	return &Dispatcher{indexedFiles: indexedFiles}
}

func (di *Dispatcher) SetIndexedFiles(indexedFiles map[string]string) {
	di.indexedFiles = indexedFiles
}

// RouteMatchings
func (di *Dispatcher) RouteMatching(w http.ResponseWriter, req *http.Request) bool {
	found := false
	for i, v := range di.indexedFiles {
		splitedKey := strings.Split(i, "|")
		fmt.Println("KEY=>", splitedKey, "VALUE=>", v)
		if glob.Glob(splitedKey[0], req.URL.Path) && req.Method == splitedKey[1] {
			w.Write([]byte("method:" + req.Method + ", route:" + req.URL.Path + ",KEY" + v))
			found = true
		}
	}

	if !found {
		//	no pattern matched; send 404 response
		http.NotFound(w, req)
	}

	return found
}

func (di *Dispatcher) TypeMatching(w http.ResponseWriter, req *http.Request) {

}

func (di *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(w, req.Method)

	if req.URL.Path == "/favicon.ico" {
		return
	}
	// TODO: should check path
	// TODO: check type can be http or graphql, grpc
	di.RouteMatching(w, req)
}

//Start initialize the HTTP mock server
func (di Dispatcher) Start(addr string, indexedFiles map[string]string) {
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
