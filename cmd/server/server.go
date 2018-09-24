package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/session"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
)

// M represents simple map interface.
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
	tpl.Load("login", "register", "client-register", "consent", "index")

	sessMgr := session.NewManager()
	sessMgr.Start()
	defer sessMgr.Stop()

	svc := NewService()

	// -- helpers

	// TODO: Allow session manager to access the *http.Request to get the cookie.
	hasSession := func(r *http.Request) bool {
		c, err := r.Cookie("id")
		if err != nil {
			return false
		}
		// If the session does not exist, an error will be thrown.
		sess, err := sessMgr.Get(c.Value)
		if err != nil {
			return false
		}
		return sess != nil
	}

	getSession := func(r *http.Request) (*session.Session, error) {
		c, err := r.Cookie("id")
		if err != nil {
			return nil, err
		}

		return sessMgr.Get(c.Value)
	}

	setSession := func(w http.ResponseWriter, userID string) {
		s := session.NewSession(userID)
		sessMgr.Put(s.SID, s)
		cookie := session.NewCookie(s.SID)
		http.SetCookie(w, cookie)
	}
	// -- endpoints

	getLogin := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// TODO: Add CSRF.
		// Check if the querystring contains the authentication request.
		// If yes, send it into the body as the request body.
		type data struct {
			ReturnURL string
		}
		if ok := hasSession(r); ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		parseURI := func(u url.Values) (string, error) {
			base64uri := u.Get("return_url")
			if base64uri == "" {
				return "/", nil
			}
			return decodeBase64(base64uri)
		}

		uri, err := parseURI(r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		d := data{uri}
		tpl.Render(w, "login", d)
	}

	postLogin := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")

		if ok := hasSession(r); ok {
			http.Error(w, "already logged in", http.StatusUnauthorized)
			return
		}
		// This inline-function is just meant to break the steps in
		// this function to determine the pipeline. It should be
		// encapsulated into the service for testability.
		performLogin := func(r *http.Request) (*oidc.User, error) {
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			var (
				email    = request.Email
				password = request.Password
			)
			return svc.user.Login(email, password)
		}

		provideToken := func(userid string) (string, error) {
			var (
				aud = "https://server.example.com/login"
				sub = userid
				iss = userid
				iat = time.Now().UTC()
				exp = iat.Add(2 * time.Hour)

				key = []byte("access_token_secret")
			)
			claims := crypto.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
			return crypto.NewJWT(key, claims)
		}

		// TODO: Check if the user has an existing session or not
		// before performing login.
		user, err := performLogin(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		accessToken, err := provideToken(user.ID)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		setSession(w, user.ID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(M{
			"access_token": accessToken,
		})
	}

	getRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// If the user already has a session (is logged in), redirect
		// them back to the home page.
		if ok := hasSession(r); ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		tpl.Render(w, "register", nil)
	}

	postRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")

		if ok := hasSession(r); ok {
			http.Error(w, "already logged in", http.StatusUnauthorized)
			return
		}

		performRegister := func(r *http.Request) (*oidc.User, error) {
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			var (
				email    = request.Email
				password = request.Password
			)
			return svc.user.Register(email, password)
		}

		provideToken := func(userid string) (string, error) {
			var (
				aud = "https://server.example.com/register"
				sub = userid
				iss = userid
				iat = time.Now().UTC()
				exp = iat.Add(2 * time.Hour)

				key = []byte("access_token_secret")
			)
			claims := crypto.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
			return crypto.NewJWT(key, claims)
		}

		// TODO: Check if the user has an existing session or not
		// before performing login.
		user, err := performRegister(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		accessToken, err := provideToken(user.ID)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		setSession(w, user.ID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(M{
			"access_token": accessToken,
		})
	}

	getAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()

		var req oidc.AuthenticationRequest
		if err := querystring.Decode(q, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		redirectToLogin := func() {
			redirectURI := getHost(r)
			redirectURI.RawQuery = q.Encode()
			base64uri := encodeBase64(redirectURI.String())
			u := fmt.Sprintf(`http://localhost:8080/login?return_url=%s`, base64uri)
			http.Redirect(w, r, u, http.StatusFound)
		}

		isAuthorized := hasSession(r)

		type response struct {
			QueryString string
		}
		res := response{q.Encode()}
		prompt := req.GetPrompt()
		switch {
		case prompt.Is(oidc.PromptNone) && isAuthorized:
			// Success
			tpl.Render(w, "consent", res)
		case prompt.Is(oidc.PromptNone) && !isAuthorized:
			http.Error(w, oidc.ErrLoginRequired.Error(), http.StatusBadRequest)
			// ErrorLoginRequired
			// http.Redirect(w, r, req.RedirectURI, http.StatusNotFound)
		case prompt.Is(oidc.PromptLogin) && !isAuthorized:
			// Force login:
			redirectToLogin()
		case prompt.Is(oidc.PromptLogin) && isAuthorized:
			tpl.Render(w, "consent", res)
		case prompt.Is(oidc.PromptConsent) && isAuthorized:
			tpl.Render(w, "consent", res)
		// case prompt.Is(oidc.PromptSelectAccount):
		default:
			http.Error(w, "invalid prompt", http.StatusBadRequest)
		}
	}

	postAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := r.Context()
		if authorized := hasSession(r); !authorized {
			// User is not logged in, throw error.
			writeError(w, http.StatusUnauthorized, errors.New("not logged in"))
			return
		}

		sess, err := getSession(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// Attach the user_id to the context.
		ctx = context.WithValue(ctx, oidc.UserContextKey, sess.UserID)

		// Construct the request payload from the querystring.
		var req oidc.AuthenticationRequest
		if err := req.FromQueryString(r.URL.Query()); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// Attempt to authenticate the user.
		res, err := svc.core.Authenticate(ctx, &req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		u, err := urlWithQuery(req.RedirectURI, res.ToQueryString())
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		http.Redirect(w, r, u, http.StatusFound)
	}

	getClientRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// TODO: Check if the user is authorized to read the client
		// details.
		id := r.URL.Query().Get("client_id")
		if id == "" {
			tpl.Render(w, "client-register", nil)
			return
		}
		client, err := svc.client.Read(id)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		json.NewEncoder(w).Encode(client)
	}

	postClientRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.ParseForm()
		// TODO: Check if the user is authorized to perform client
		// registration.

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
			v, ok := err.(*oidc.ErrorJSON)
			if ok {
				json.NewEncoder(w).Encode(v)
			} else {
				json.NewEncoder(w).Encode(M{"error": err.Error()})
			}
			return
		}
		log.Println("registered:", newClient.ClientID)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(M{"success": true})
	}

	postLogout := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cookie, err := r.Cookie("id")
		if err != nil {
			// ErrNoCookie should be handled as success.
			writeError(w, http.StatusBadRequest, err)
			return
		}
		sessMgr.Delete(cookie.Value)
		json.NewEncoder(w).Encode(M{
			"ok": true,
		})
	}

	getIndex := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		type data struct {
			IsLoggedIn bool
		}
		var res data
		sess, err := getSession(r)
		if err != nil {
			res.IsLoggedIn = false
		}
		if sess != nil {
			res.IsLoggedIn = true
		}
		tpl.Render(w, "index", res)
	}

	r.GET("/", getIndex)

	r.POST("/logout", postLogout)
	r.GET("/register", getRegister)
	r.GET("/login", getLogin)
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

func urlWithQuery(uri string, q url.Values) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// getHost tries its best to return the request host.
func getHost(r *http.Request) *url.URL {
	u := r.URL

	// The scheme is http because that's the only protocol your server handles.
	u.Scheme = "http"

	// If client specified a host header, then use it for the full URL.
	u.Host = r.Host

	// Otherwise, use your server's host name.
	if u.Host == "" {
		u.Host = "your-host-name.com"
	}
	// if r.URL.IsAbs() {
	//         host := r.Host
	//         // Slice off any port information.
	//         if i := strings.Index(host, ":"); i != -1 {
	//                 host = host[:i]
	//         }
	//         u.Host = host
	// }

	return u
}
