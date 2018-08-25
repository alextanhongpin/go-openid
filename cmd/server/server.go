package main

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	db := NewDatabase()

	e := Endpoints{
		service: NewService(db),
	}
	// Create multiple endpoints
	r := httprouter.New()
	r.GET("/connect/register", e.Register)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
