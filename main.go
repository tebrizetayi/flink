package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tebrizetayi/flink/httpapi"
)

func main() {

	port := os.Getenv("HISTORY_SERVER_LISTEN_ADDR")
	if port == "" {
		port = ":8080"
	}
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the HTTP service listening for requests.
	api := http.Server{
		Addr:           port,
		Handler:        httpapi.API(shutdown),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("main : API Listening %s", port)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("main : Error starting server: %+v", err)

	case sig := <-shutdown:
		log.Printf("main : %v : Start shutdown..", sig)
	}
}
