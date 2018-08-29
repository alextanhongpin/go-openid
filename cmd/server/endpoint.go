package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	oidc "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
	"github.com/julienschmidt/httprouter"
)

type Endpoints struct {
	service OIDService
}

func (e *Endpoints) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check if authorization header exists, and is valid
	// Can be extracted as a middleware

	// Construct request parameters
	var req oidc.AuthorizationRequest
	if err := querystring.Decode(&req, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Prepare redirect uri
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	// Call service
	res, err := e.service.Authorize(r.Context(), req)
	if err != nil {
		q := querystring.Encode(err)
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
		return
	}

	q := querystring.Encode(res)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (e *Endpoints) Token(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req oidc.AccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := e.service.Token(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// TODO: What status type to return here?
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	json.NewEncoder(w).Encode(res)
}

func (e *Endpoints) RegisterClient(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check for authorization headers to see if the client can register
	var req oidc.ClientRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := e.service.RegisterClient(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Set the appropriate headers
	w.WriteHeader(http.StatusCreated)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	json.NewEncoder(w).Encode(res)
}

func (e *Endpoints) Client() {
	// GET /connect/register?client_id=
	// Authorization: Bearer this.is.an.access.token.value
	// return 200, cache-control: no-store, pragma: no-cache
	// Client does not exist, invalid client, invalid token returns 401 unauthorized
	// No permission: 403 forbidden
	// Do not return 404
}

// .well-known/webfinger
// .well-known/openid-configuration
func (e *Endpoints) Authenticate(ctx context.Context, req *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error) {

	return nil, nil
}
func (e *Endpoints) RefreshToken() {}

func (e *Endpoints) UserInfo() {

}
