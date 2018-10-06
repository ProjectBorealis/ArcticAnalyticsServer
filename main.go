package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/ProjectBorealis/ArcticAnalyticsServer/server"
)

const (
	sharedSecretEnv = "SHARED_SECRET"
)

func main() {
	rw, err := server.NewResultWriter("results.json")
	if err != nil {
		panic(err)
	}
	defer rw.Close()

	s := server.New(rw, os.Getenv(sharedSecretEnv))
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.Handler())

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      loggedRouter,
		Addr:         ":9095",
	}
	log.Println(srv.ListenAndServe())
}
