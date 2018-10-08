package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"projectborealisgitlab.site/project-borealis/programming/dev-ops/aa-server/server"
)

const (
	sharedSecretEnv = "SHARED_SECRET"
)

var (
	behindProxy = flag.Bool("behind-proxy", false, "whether server is running behind a reverse proxy.")
	httpPort    = flag.Int("port", 9095, "http server port")
)

func main() {
	flag.Parse()

	rw, err := server.NewResultWriter("results.json")
	if err != nil {
		panic(err)
	}
	defer rw.Close()

	handler := server.New(rw, os.Getenv(sharedSecretEnv)).Handler()

	// enable proxy middleware if behind a proxy
	if *behindProxy {
		handler = handlers.ProxyHeaders(handler)
	}

	// enable logging
	handler = handlers.LoggingHandler(os.Stdout, handler)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handler,
		Addr:         ":9095",
	}
	log.Println(srv.ListenAndServe())
}
