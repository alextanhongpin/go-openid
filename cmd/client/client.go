package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/core"
	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
)

type M map[string]interface{}

func main() {
	cfg := NewConfig()

	tpl := html5.New(cfg.TemplateDir)
	tpl.Load("authorize")

	getIndex := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tpl.Render(w, "authorize", nil)
	}

	openidClient := core.Client{
		ClientID:             cfg.ClientID,
		ClientSecret:         cfg.ClientSecret,
		TokenRegistrationURI: "http://localhost:8080/token",
	}

	getAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Construct redirect url here.
		// Based on the prompt type, if it is the popup, an iframe
		// should be displayed to let the users login without being
		// redirected. The question is, how to maintain the login
		// token/session in the popup.
		req := openid.AuthenticationRequest{
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

	getAuthorizeCallback := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()
		var authzReq openid.AuthorizationResponse
		if err := querystring.Decode(q, &authzReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var (
			code        = authzReq.Code
			redirectURI = "http://localhost:4000/authorize/callback"
		)

		res, err := openidClient.Exchange(r.Context(), code, redirectURI)
		if err != nil {
			json.NewEncoder(w).Encode(M{
				"error": err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(res)
	}

	r := httprouter.New()

	r.GET("/", getIndex)
	r.GET("/authorize", getAuthorize)
	r.GET("/authorize/callback", getAuthorizeCallback)

	srv := gsrv.New(cfg.Port, r)
	<-srv
	log.Println("shutting down server.")
}
