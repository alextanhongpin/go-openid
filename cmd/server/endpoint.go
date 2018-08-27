package main

import (
	"encoding/json"
	"net/http"

	oidc "github.com/alextanhongpin/go-openid"
	"github.com/julienschmidt/httprouter"
)

type Endpoints struct {
	service OIDService
}

func (e *Endpoints) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check if authorization header exists, and is valid
	// Can be extracted as a middleware

	// Construct request parameters
	req := oidc.DecodeAuthorizationRequest(r.URL.Query())

	// Call service
	res, err := e.service.Authorize(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	redirectURI, err := oidc.EncodeAuthorizationResponse(res, req.RedirectURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	http.Redirect(w, r, redirectURI.String(), http.StatusFound)
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
	// Set the appropriate headers
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	json.NewEncoder(w).Encode(res)
}
