package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/auth"
)

func main() {
	var (
		port = flag.Int("port", 8080, "The port for the server")
		//env  = flag.String("env", "dev", "The working environment dev|stage|test|prod")
	)

	r := httprouter.New()
	db := app.Database()
	tmpl := app.Template()

	auth.SetupAuth(db, r, tmpl)
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("hello"))
	})

	fmt.Printf("listening to port*:%d. press ctrl + c to cancel", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}
