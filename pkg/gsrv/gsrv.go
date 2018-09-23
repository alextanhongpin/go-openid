package gsrv

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// New returns a server with graceful shutdown.
func New(port int, r http.Handler) <-chan struct{} {
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	idle := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, os.Kill)

		<-sigint

		// Receive interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listener, or context timeout.
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idle)
	}()

	log.Printf("listening to port *:%d. press ctrl + c to cancel.", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener.
		log.Printf("HTTP server ListenAndServe: %v", err)
	}
	return idle
}
