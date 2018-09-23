package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
)

func main() {
	cfg := NewConfig()

	tpl := html5.New(cfg.TemplateDir)
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
		req := oidc.AuthenticationRequest{
			ResponseType: "code",
			Scope:        "openid profile email",
			ClientID:     cfg.ClientID,
			State:        "abc",
			RedirectURI:  "http://localhost:4000/authorize/callback",

			Prompt: "login",
		}
		u, err := url.Parse("http://localhost:8080/authorize")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		qs := querystring.Encode(url.Values{}, req)
		u.RawQuery = qs.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
	}

	postAuthorizeCallback := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}

	r := httprouter.New()
	r.GET("/", getIndex)
	r.GET("/authorize", getAuthorize)
	r.POST("/authorize/callback", postAuthorizeCallback)

	srv := gsrv.New(cfg.Port, r, "", "")
	<-srv
	log.Println("shutting down server.")
}
