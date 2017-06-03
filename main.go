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

	// Create a logger context

	env := app.Env{
		Db:     app.Database(),
		Router: httprouter.New(),
		Tmpl:   app.MakeTemplate(),
		// Log: app.Logger()
	}

	auth.FeatureToggle(true)(env)

	env.Router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("hello"))
	})

	fmt.Printf("listening to port*:%d. press ctrl + c to cancel", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), env.Router))
}
