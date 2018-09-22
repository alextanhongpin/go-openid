package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/core"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
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
	tpl.Load("login", "register", "client-register", "consent")

	svc := NewService()

	getLogin := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// TODO: Add CSRF.
		// Check if the querystring contains the authentication request.
		// If yes, send it into the body as the request body.
		type data struct {
			ReturnURL string
		}
		// By default, redirect to the index page.
		uri := "/"
		if base64uri := r.URL.Query().Get("return_url"); base64uri != "" {
			u, err := decodeBase64(base64uri)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			uri = u
		}

		d := data{uri}
		tpl.Render(w, "login", d)
	}

	postLogin := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		type loginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		u, err := svc.user.Login(req.Email, req.Password)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		// TODO: Move logic to the service.
		res, err := core.NewToken(u)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		cookie := http.Cookie{
			Name:     "auth",
			Value:    res.AccessToken,
			HttpOnly: true,
			Secure:   true, // Only for https
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}

	getRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tpl.Render(w, "register", nil)
	}

	postRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		type registerRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")

		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		u, err := svc.user.Register(req.Email, req.Password)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// TODO: Move logic to the service.
		res, err := core.NewToken(u)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}

	makeAuthorizeURI := func(q url.Values) (string, error) {
		u, err := url.Parse("http://localhost:8080/authorize")
		if err != nil {
			return "", err
		}
		u.RawQuery = q.Encode()
		return u.String(), nil
	}

	getAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()
		authorizeURI, err := makeAuthorizeURI(q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var req oidc.AuthenticationRequest
		if err := querystring.Decode(q, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, _ := r.Cookie("auth")
		isAuthorized := token != nil && token.Value != ""
		if isAuthorized {
			log.Println("got token:", token.Value)
			t := strings.Split(token.Value, "=")[1]
			claims, err := crypto.ParseJWT([]byte("access_token_secret"), t)
			if err != nil {
				log.Println("getAuthorize: err parsing jwt", err)
				b64uri := encodeBase64(authorizeURI)
				u := fmt.Sprintf(`http://localhost:8080/login?return_url=%s`, b64uri)
				http.Redirect(w, r, u, http.StatusFound)
				return
			}
			log.Println(claims)
		}
		log.Println("getAuthorize: got cookie", token, isAuthorized)
		// Check the prompt type here. If login is required, direct them to the login page.
		if prompt := req.GetPrompt(); prompt.Is(oidc.PromptLogin) && !isAuthorized {
			b64uri := encodeBase64(authorizeURI)
			u := fmt.Sprintf(`http://localhost:8080/login?return_url=%s`, b64uri)
			http.Redirect(w, r, u, http.StatusFound)
			return
		}
		tpl.Render(w, "consent", nil)
	}

	postAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Generate code to be exchanged as token.
		log.Println("postAuthorize", r.URL.Query())
		// Delete cookie here.
		json.NewEncoder(w).Encode(M{"ok": true})
	}

	getClientRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		clientID := r.URL.Query().Get("client_id")
		if clientID != "" {
			client, err := svc.client.Read(clientID)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			json.NewEncoder(w).Encode(client)
			return
		}
		tpl.Render(w, "client-register", nil)
	}

	postClientRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.ParseForm()

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")

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
		json.NewEncoder(w).Encode(M{"success": true})
	}
	r.GET("/", getLogin)

	r.GET("/register", getRegister)
	r.GET("/login", getLogin)
	r.POST("/login", postLogin)
	r.POST("/register", postRegister)
	r.GET("/connect/register", getClientRegister)
	r.POST("/connect/register", postClientRegister)
	r.GET("/authorize", getAuthorize)
	r.POST("/authorize", postAuthorize)

	srv := gsrv.New(*port, r, "server.crt", "server.key")
	<-srv
	log.Println("Gracefully shutdown HTTP server.")
}

func decodeBase64(in string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(in)
	return string(b), err
}

func encodeBase64(in string) string {
	return base64.URLEncoding.EncodeToString([]byte(in))
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(M{
		"error": err.Error(),
	})
}
