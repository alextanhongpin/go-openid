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

	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

var (
	templates map[string]*template.Template
	once      sync.Once
	files     = []string{"login"}
)

func init() {
	once.Do(func() {
		load := func(f string) string {
			return fmt.Sprintf("templates/%s.tmpl", f)
		}

		templates = make(map[string]*template.Template)
		layout := template.Must(template.New("base").ParseFiles(load("base")))
		for _, f := range files {
			clone := template.Must(layout.Clone())
			templates[f] = template.Must(clone.ParseFiles(load(f)))
		}
	})
}

func main() {
	port := flag.Int("port", 8080, "the port of the application")
	flag.Parse()

	idle := make(chan struct{})

	var e *Endpoints
	{
		// Factory setup
		db := database.NewInMem()
		c := crypto.New(defaultJWTSigningKey)
		svc := NewService(db, c)
		e = NewEndpoints(svc)
	}

	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		renderTemplate(w, "login", nil)
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

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	if t, ok := templates[name]; !ok {
		err := fmt.Sprintf("template with the name %s does not exist", name)
		http.Error(w, err, http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t.Execute(w, data)
	}
}
