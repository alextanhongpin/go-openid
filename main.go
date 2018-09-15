package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/middlewares"
	"github.com/alextanhongpin/go-openid/service/authsvc"
	"github.com/alextanhongpin/go-openid/service/clientsvc"
	"github.com/alextanhongpin/go-openid/service/oauthsvc"
	"github.com/alextanhongpin/go-openid/service/tokensvc"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
)

// main is where our application lives
func main() {
	var (
		port      = flag.Int("port", 8080, "The port for the server")
		dsn       = flag.String("dsn", "127.0.0.1:27017", "The data source name for the database")
		dbn       = flag.String("dbn", "go-openid", "The database name")
		redisHost = flag.String("redis_host", "localhost:6379", "The redis host name")
	)
	flag.Parse()

	govalidator.SetFieldsRequiredByDefault(true)

	r := httprouter.New()

	tracer := app.NewTracer()
	defer tracer.Close()

	db := app.NewDatabase(*dbn, *dsn)
	defer db.Close()

	cache := app.NewCache(*redisHost)

	env := app.Env{
		DB:     db,
		Router: r,
		Tmpl:   app.NewTemplate(),
		Cache:  cache,
		Tracer: tracer,
		// TODO: Create a logger context
		// Log: app.Logger()
	}

	// Serve public files
	r.ServeFiles("/public/*filepath", http.Dir("public"))

	// Toggle services here
	authsvc.FeatureToggle(true)(env)
	clientsvc.FeatureToggle(true)(env)
	oauthsvc.FeatureToggle(true)(env)
	tokensvc.FeatureToggle(true)(env)

	r.GET("/", middlewares.ProtectRoute(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		env.Tmpl.Render(w, "index", nil)
	}))

	log.Printf("listening to port *:%d. press ctrl + c to cancel\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), env.Router))
}
