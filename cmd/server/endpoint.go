package main

import (
	"encoding/json"
	"net/http"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/julienschmidt/httprouter"
)

type Endpoints struct {
	service Service
}

func (e *Endpoints) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Construct request parameters
	req := openid.AuthorizationRequest{}

	// Call service
	res, err := e.service.Authorize(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Set the appropriate headers
	w.Header().Set("Cache-Control", "no-control")
	w.Header().Set("Pragma", "no-cache")

	// Return json response/redirect
	json.NewEncoder(w).Encode(res)
}

func (e *Endpoints) Token(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
