package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"

	oidc "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
)

// Endpoints represent the endpoints for the OpenIDConnect.
type Endpoints struct {
	service Service
}

// NewEndpoints returns a pointer to new endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		service: s,
	}
}

// Authorize performs the authorization logic.
func (e *Endpoints) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Construct request parameters
	var req oidc.AuthorizationRequest
	if err := querystring.Decode(&req, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if !govalidator.IsURL(req.RedirectURI) {
		http.Error(w, oidc.InvalidRedirectURI.String(), http.StatusForbidden)
		return
	}

	// Prepare redirect uri
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Call service
	res, authErr := e.service.Authorize(r.Context(), &req)
	if authErr != nil {
		q := querystring.Encode(authErr)
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
		return
	}

	q := querystring.Encode(res)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// Token represents the token service.
func (e *Endpoints) Token(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	auth := r.Header.Get("Authorization")
	if len(auth) < 7 || auth[0:5] != "Basic" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(oidc.InvalidRequest.JSON())
		return
	}

	clientID, clientSecret := oidc.DecodeClientAuth(auth[6:])
	if err := e.service.ValidateClient(clientID, clientSecret); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
		return
	}

	var req oidc.AccessTokenRequest
	if err := querystring.Decode(&req, r.Form); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(oidc.InvalidRequest.JSON())
		return
	}

	res, err := e.service.Token(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// RegisterClient represents the endpoint for client registration.
func (e *Endpoints) RegisterClient(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	auth := r.Header.Get("Authorization")
	if len(auth) < 8 || auth[0:6] != "Bearer" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(oidc.InvalidRequest.JSON())
		return
	}

	_, err := e.service.ParseJWT(auth[7:])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	// Check for authorization headers to see if the client can register
	var req oidc.ClientRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(oidc.InvalidRequest.JSON())
		return
	}

	res, err := e.service.RegisterClient(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(oidc.InvalidRequest.JSON())
		return
	}

	// Set the appropriate headers
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(res)
}

// Client returns the authorized client information.
func (e *Endpoints) Client(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	// TODO: Check authorization header to ensure the client has the right credentials.
	auth := r.Header.Get("Authorization")
	if len(auth) < 8 || auth[0:6] != "Bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(oidc.InvalidRequest.JSON())
		return
	}

	// TODO: Check the user status from cache.
	_, err := e.service.ParseJWT(auth[7:])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	id := r.URL.Query().Get("client_id")
	client, err := e.service.Client(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(client)
	// GET /connect/register?client_id=
	// Authorization: Bearer this.is.an.access.token.value
	// return 200, cache-control: no-store, pragma: no-cache
	// Client does not exist, invalid client, invalid token returns 401 unauthorized
	// No permission: 403 forbidden
	// Do not return 404
}

// .well-known/webfinger
func (e *Endpoints) Webfinger()     {}
func (e *Endpoints) Configuration() {}

// .well-known/openid-configuration
func (e *Endpoints) Authenticate(ctx context.Context, req *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error) {

	return nil, nil
}

// RefreshToken returns a new refresh token alongside with the id token.
func (e *Endpoints) RefreshToken() {}

func validateTokenHeader(token string) (string, error) {
	return "1", nil
}

func (e *Endpoints) UserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	auth := r.Header.Get("Authorization")
	if auth[0:6] != "Bearer" {
		// Error
		err := oidc.ErrUnauthorizedClient
		msg := fmt.Sprintf(`error="%s" error_description="%s"`, err.Error(), "The access token expired")
		w.Header().Set("WWW-Authenticate", msg)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	token := auth[7:]

	// TODO: Receive the correct token type.
	claims, err := e.service.ParseJWT(token)
	if err != nil {
		// TODO: Return the correct error.
		err := oidc.ErrUnauthorizedClient
		msg := fmt.Sprintf(`error="%s" error_description="%s"`, err.Error(), "The access token expired")
		w.Header().Set("WWW-Authenticate", msg)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	res, err := e.service.UserInfo(r.Context(), claims.UserID)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(res)
}
