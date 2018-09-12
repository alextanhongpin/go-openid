package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

// HTMLs represent the html templates stored as dictionary.
type HTMLs map[string]*template.Template

var (
	htmls HTMLs
	once  sync.Once
)

func init() {
	htmls = make(HTMLs)
	once.Do(initTemplates(htmls, "login"))
}

func main() {
	port := flag.Int("port", 8080, "the port of the application")
	flag.Parse()

	idle := make(chan struct{})

	e := initEndpoints(defaultJWTSigningKey)
	r := httprouter.New()

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
}

// Render renders the html output with the given data.
func (h HTMLs) Render(w http.ResponseWriter, name string, data interface{}) {
	if t, ok := h[name]; !ok {
		err := fmt.Sprintf("template with the name %s does not exist", name)
		http.Error(w, err, http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t.Execute(w, data)
	}
}
