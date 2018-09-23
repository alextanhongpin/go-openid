package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/authheader"
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

	getAuthorizeCallback := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()
		var authzReq oidc.AuthorizationResponse
		if err := querystring.Decode(q, &authzReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenReq := oidc.AccessTokenRequest{
			GrantType:   "authorization_code",
			Code:        authzReq.Code,
			RedirectURI: "http://localhost:4000/authorize/callback",
		}
		jsonBody, err := json.Marshal(tokenReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequest("POST", "http://localhost:8080/token", bytes.NewBuffer(jsonBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req = req.WithContext(ctx)
		req.Header.Add("Authorization", "Basic "+authheader.EncodeBase64(cfg.ClientID, cfg.ClientSecret))

		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var res oidc.AuthenticationResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			panic(err)
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
