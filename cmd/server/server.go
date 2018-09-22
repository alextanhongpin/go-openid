package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
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
		// TODO: Add CSRF.
		// Check if the querystring contains the authentication request.
		// If yes, send it into the body as the request body.
		q := r.URL.Query()
		var req oidc.AuthenticationRequest
		if err := querystring.Decode(q, &req); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Sign the payload with JWT?
		tpl.Render(w, "login", req)
	}

	getRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tpl.Render(w, "register", nil)
	}

	getClientRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		clientID := r.URL.Query().Get("client_id")
		if clientID != "" {
			client, err := svc.client.Read(clientID)
			if err != nil {
				json.NewEncoder(w).Encode(M{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(client)
			return
		}
		tpl.Render(w, "client-register", nil)
	}

	postClientRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.ParseForm()

		var (
			clientName   = r.FormValue("client_name")
			redirectURIs = strings.Split(r.FormValue("redirect_uris"), " ")
		)

		client := oidc.NewClient()
		client.ClientName = clientName
		client.RedirectURIs = redirectURIs

		newClient, err := svc.client.Register(client)
		if err != nil {
			res := M{"error": err.Error()}
			v, ok := err.(*oidc.ErrorJSON)
			if ok {
				res["error"] = v.Code
				res["error_description"] = v.Description
			}
			json.NewEncoder(w).Encode(res)
			return
		}
		log.Println("registered:", newClient.ClientID)

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")
		json.NewEncoder(w).Encode(M{"success": true})
	}

	postLogin := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.ParseForm()

		var (
			email    = r.FormValue("email")
			password = r.FormValue("password")
		)
		u, err := svc.user.Login(email, password)
		if err != nil {
			json.NewEncoder(w).Encode(M{
				"error": "email of password is invalid",
			})
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")
		json.NewEncoder(w).Encode(M{"user": u})
	}

	postRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.ParseForm()

		var (
			email    = r.FormValue("email")
			password = r.FormValue("password")
		)
		if err := svc.user.Register(email, password); err != nil {
			json.NewEncoder(w).Encode(M{
				"error": err.Error(),
			})
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")
		json.NewEncoder(w).Encode(M{"success": true})
	}

	getAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()
		var req oidc.AuthenticationRequest
		if err := querystring.Decode(q, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Check the prompt type here. If login is required, direct them to the login page.
		if prompt := req.GetPrompt(); prompt.Is(oidc.PromptLogin) {
			// Redirect to login page while maintaining the data.
			u, err := url.Parse("http://localhost:8080/login")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			u.RawQuery = q.Encode()
			http.Redirect(w, r, u.String(), http.StatusFound)
			return
		}
		json.NewEncoder(w).Encode(req)
	}

	postAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Generate code to be exchanged as token.
	}

	r.GET("/", getLogin)
	r.GET("/register", getRegister)
	r.POST("/login", postLogin)
	r.POST("/register", postRegister)
	r.GET("/connect/register", getClientRegister)
	r.POST("/connect/register", postClientRegister)
	r.GET("/authorize", getAuthorize)
	r.POST("/authorize", postAuthorize)

	srv := gsrv.New(*port, r)
	<-srv
	log.Println("Gracefully shutdown HTTP server.")
}
