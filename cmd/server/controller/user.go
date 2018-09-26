package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/alextanhongpin/go-openid/internal/user"
	"github.com/alextanhongpin/go-openid/pkg/appsensor"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/alextanhongpin/go-openid/service"

	"github.com/julienschmidt/httprouter"
)

type (
	// Credentials represent the user credentials for the application.
	Credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// User represents the configuration for the user controller.
	User struct {
		service   service.User
		template  *html5.Template
		session   *session.Manager
		appsensor appsensor.LoginDetector
	}
)

// NewUser returns a new user controller with a predefined service.
func NewUser() User {
	return User{
		service: user.NewService(),
	}
}

// -- setters

// SetService sets the current service.
func (u *User) SetService(s service.User) {
	u.service = s
}

// SetTemplate sets the current template.
func (u *User) SetTemplate(t *html5.Template) {
	u.template = t
}

// SetSession sets the current session.
func (u *User) SetSession(s *session.Manager) {
	u.session = s
}

// SetAppSensor sets the current appsensor.
func (u *User) SetAppSensor(a appsensor.LoginDetector) {
	u.appsensor = a
}

// GetLogin renders the login html.
func (u *User) GetLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// TODO: Add CSRF.
	// Check if the querystring contains the authentication request.
	// If yes, send it into the body as the request body.
	type templateData struct {
		ReturnURL string
	}

	// TODO: The user might have a session, but the session has
	// expired. Need to invalidate the user first by deleting the
	// old session, and creating a new one.
	if ok := u.session.HasSession(r); ok {
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

	d := templateData{uri}
	u.template.Render(w, "login", d)
}

// PostLogin represents the post login endpoint.
func (u *User) PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")

	if ok := u.session.HasSession(r); ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var req Credentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if locked := u.appsensor.IsLocked(req.Email); locked {
		writeError(w, http.StatusTooManyRequests, errors.New("too many attempts"))
		return
	}

	user, err := u.service.Login(req.Email, req.Password)
	if err != nil {
		// Log attempts here.
		u.appsensor.Increment(req.Email)
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

	u.session.SetSession(w, user.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(M{
		"access_token": accessToken,
	})
}

// GetRegister renders the user registration page.
func (u *User) GetRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// If the user already has a session (is logged in), redirect
	// them back to the home page.
	if ok := u.session.HasSession(r); ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	u.template.Render(w, "register", nil)
}

// PostRegister represents the post register endpoint.
func (u *User) PostRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")

	if ok := u.session.HasSession(r); ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var req Credentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := u.service.Register(req.Email, req.Password)
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

	u.session.SetSession(w, user.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(M{
		"access_token": accessToken,
	})
}

// PostLogout logs the user out of the session.
func (u *User) PostLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie, err := r.Cookie(session.Key)
	if err != nil {
		// ErrNoCookie should be handled as success.
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := u.session.Delete(cookie.Value); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	// TODO: Look into the PRG pattern.
	http.Redirect(w, r, "/", http.StatusFound)
}
