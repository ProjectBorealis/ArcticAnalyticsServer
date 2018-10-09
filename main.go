package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/handlers"
	"projectborealisgitlab.site/project-borealis/programming/dev-ops/aa-server/server"
)

const (
	sharedSecretEnv  = "SHARED_SECRET"
	adminPasswordEnv = "ADMIN_PASSWORD"
)

var (
	behindProxy = flag.Bool("behind-proxy", false, "whether server is running behind a reverse proxy.")
	httpAddr    = flag.String("addr", ":9095", "http server port")
	dataDir     = flag.String("data-dir", "./", "data directory")
)

func main() {
	flag.Parse()

	rw, err := server.NewResultWriter(path.Join(*dataDir, "results.json"))
	if err != nil {
		panic(err)
	}
	defer rw.Close()

	handler := server.New(rw, os.Getenv(sharedSecretEnv), os.Getenv(adminPasswordEnv)).Handler()

	// enable proxy middleware if behind a proxy
	if *behindProxy {
		handler = handlers.ProxyHeaders(handler)
	}

	// enable logging
	handler = handlers.LoggingHandler(os.Stdout, handler)

	// enable compression
	handler = handlers.CompressHandler(handler)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handler,
		Addr:         *httpAddr,
	}

	log.Printf("Data directory: %v\n", *dataDir)
	log.Printf("Listening %v...\n", *httpAddr)
	log.Println(srv.ListenAndServe())
}
