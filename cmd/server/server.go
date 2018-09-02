package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	var (
		port = 8080
	)
	db := NewDatabase()
	svc := NewService(db, nil)
	e := NewEndpoints(svc)
	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("hello world"))
	})
	r.GET("/connect/register", e.RegisterClient)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("listening to port *:%d. press ctrl + c to cancel.", port)
	log.Fatal(srv.ListenAndServe())
}
