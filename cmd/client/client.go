package main

import (
	"log"
	"net/http"

	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/julienschmidt/httprouter"
)

func main() {
	port := 4000
	tpldir := "templates"

	tpl := html5.New(tpldir)
	tpl.Load("authorize")

	getIndex := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tpl.Render(w, "authorize", nil)
	}

	getAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Construct redirect url here.
		// Based on the prompt type, if it is the popup, an iframe
		// should be displayed to let the users login without being
		// redirected. The question is, how to maintain the login
		// token/session in the popup.
	}

	postAuthorizeCallback := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}

	r := httprouter.New()
	r.GET("/", getIndex)
	r.GET("/authorize", getAuthorize)
	r.POST("/authorize/callback", postAuthorizeCallback)

	srv := gsrv.New(port, r)
	<-srv
	log.Println("shutting down server.")
}
