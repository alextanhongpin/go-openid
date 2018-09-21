package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
)

type M map[string]interface{}

func main() {
	var (
		port   = flag.Int("port", 8080, "the port of the application")
		tplDir = flag.String("tpldir", "templates", "the datadir of the html templates")
	)
	flag.Parse()

	// Create new router.
	r := httprouter.New()

	// Load templates.
	tpl := html5.New(*tplDir)
	tpl.Load("login", "register", "client-register")

	svc := NewService()

	getLogin := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tpl.Render(w, "login", nil)
	}

	getRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tpl.Render(w, "register", nil)
	}

	getClientRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tpl.Render(w, "client-register", nil)
	}

	postLogin := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.ParseForm()

		var (
			email    = r.FormValue("email")
			password = r.FormValue("password")
		)
		u, err := svc.Login(email, password)
		if err != nil {
			json.NewEncoder(w).Encode(M{
				"error": "email of password is invalid",
			})
			return
		}
		json.NewEncoder(w).Encode(M{"user": u})
	}

	postRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.ParseForm()

		var (
			email    = r.FormValue("email")
			password = r.FormValue("password")
		)
		if err := svc.Register(email, password); err != nil {
			json.NewEncoder(w).Encode(M{
				"error": err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(M{"success": true})
	}

	getAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}

	postAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}

	r.GET("/", getLogin)
	r.GET("/register", getRegister)
	r.POST("/login", postLogin)
	r.POST("/register", postRegister)
	r.GET("/client/register", getClientRegister)
	r.GET("/authorize", getAuthorize)
	r.POST("/authorize", postAuthorize)

	srv := gsrv.New(*port, r)
	<-srv
	log.Println("Gracefully shutdown HTTP server.")
}
