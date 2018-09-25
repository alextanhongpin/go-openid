package main

import (
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
	"github.com/alextanhongpin/go-openid/pkg/appsensor"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
	"github.com/alextanhongpin/go-openid/pkg/session"
)

// M represents simple map interface.
type M map[string]interface{}

// Credentials represent the user credentials for the application.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

	// TODO: Run a cron job that handles deletion of unused data for a
	// certain period of time.
	aps := appsensor.NewLoginDetector()

	svc := NewService()

	// -- endpoints

	getLogin := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// TODO: Add CSRF.
		// Check if the querystring contains the authentication request.
		// If yes, send it into the body as the request body.
		type data struct {
			ReturnURL string
		}

		// TODO: The user might have a session, but the session has
		// expired. Need to invalidate the user first by deleting the
		// old session, and creating a new one.
		if ok := sessMgr.HasSession(r); ok {
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

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")

		if ok := sessMgr.HasSession(r); ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		var req Credentials
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		if locked := aps.IsLocked(req.Email); locked {
			writeError(w, http.StatusTooManyRequests, errors.New("too many attempts"))
			return
		}

		user, err := svc.user.Login(req.Email, req.Password)
		if err != nil {
			// Log attempts here.
			aps.Increment(req.Email)
			writeError(w, http.StatusBadRequest, err)
			return
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

		accessToken, err := provideToken(user.ID)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		sessMgr.SetSession(w, user.ID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(M{
			"access_token": accessToken,
		})
	}

	getRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// If the user already has a session (is logged in), redirect
		// them back to the home page.
		if ok := sessMgr.HasSession(r); ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		tpl.Render(w, "register", nil)
	}

	postRegister := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pragma", "no-cache")

		if ok := sessMgr.HasSession(r); ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		var req Credentials
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		user, err := svc.user.Register(req.Email, req.Password)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
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

		accessToken, err := provideToken(user.ID)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		sessMgr.SetSession(w, user.ID)

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

		if err := svc.core.PreAuthenticate(&req); err != nil {
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

		isAuthorized := sessMgr.HasSession(r)
		prompt := req.GetPrompt()

		// If the prompt is set to none, but the user is unauthorized,
		// an error should be returned indicating that login is
		// required.
		if prompt.Is(oidc.PromptNone) && !isAuthorized {
			http.Error(w, oidc.ErrLoginRequired.Error(), http.StatusBadRequest)
			return
		}

		// If the user is not authorized, login them first.
		if !isAuthorized {
			redirectToLogin()
			return
		}

		// Get the current session in order to check if the user has a
		// valid session.
		// sess, err := sessMgr.GetSession(r)
		// if err != nil {
		//         http.Error(w, oidc.ErrLoginRequired.Error(), http.StatusBadRequest)
		//         return
		// }

		type response struct {
			QueryString string
		}
		res := response{q.Encode()}
		tpl.Render(w, "consent", res)
	}

	postAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := r.Context()

		// User needs to have a session in order to call the post
		// authorize endpoint.
		sess, err := sessMgr.GetSession(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// Attach the user_id to the context.
		ctx = oidc.SetUserIDContextKey(ctx, sess.UserID)

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

		u, err := buildURL(req.RedirectURI, res.ToQueryString())
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

		buildClient := func(r *http.Request) *oidc.Client {
			var (
				clientName   = r.FormValue("client_name")
				redirectURIs = strings.Split(r.FormValue("redirect_uris"), " ")
			)
			client := oidc.NewClient()
			client.ClientName = clientName
			client.RedirectURIs = redirectURIs
			return client
		}

		client := buildClient(r)

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

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newClient)
	}

	postLogout := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cookie, err := r.Cookie(session.Key)
		if err != nil {
			// ErrNoCookie should be handled as success.
			writeError(w, http.StatusBadRequest, err)
			return
		}

		if err := sessMgr.Delete(cookie.Value); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// TODO: Look into the PRG pattern.
		http.Redirect(w, r, "/", http.StatusFound)
	}

	getIndex := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		type data struct {
			IsLoggedIn bool
		}
		var res data
		sess, err := sessMgr.GetSession(r)
		if err != nil {
			res.IsLoggedIn = false
		}
		if sess != nil {
			res.IsLoggedIn = true
		}
		tpl.Render(w, "index", res)
	}

	postToken := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := r.Context()

		// Checks if the user has an active session. If the session is
		// not active, user does not have the right credentials to
		// access the token endpoint.
		sess, err := sessMgr.GetSession(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// Put the extra data in the context to be validated by the
		// service.
		authorization := r.Header.Get("Authorization")
		ctx = oidc.SetAuthContextKey(ctx, authorization)
		ctx = oidc.SetUserIDContextKey(ctx, sess.UserID)

		// Parse request body.
		var req oidc.AccessTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		res, err := svc.core.Token(ctx, &req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		json.NewEncoder(w).Encode(res)
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
	r.POST("/token", postToken)

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

func buildURL(uri string, q url.Values) (string, error) {
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
