package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	port := flag.Int("port", 8080, "the port of the application")
	flag.Parse()

	e := initEndpoints(defaultJWTSigningKey)
	r := httprouter.New()

	htmls := NewHTMLs("templates")
	htmls.Load("login")

	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		htmls.Render(w, "login", nil)
	})
	r.GET("/authorize", e.Authorize)
	r.POST("/token", e.Token)
	r.GET("/connect/register", e.RegisterClient)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	idle := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Receive interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listener, or context timeout.
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idle)
	}()

	log.Printf("listening to port *:%d. press ctrl + c to cancel.", *port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener.
		log.Printf("HTTP server ListenAndServe: %v", err)
	}
	<-idle
	log.Println("Gracefully shutdown HTTP server.")
}
