package oidc

import (
	"log"
	"net/http"
	"net/url"
)

type AuthorizationRequest struct {
	ResponseType string `json:"response_type"`
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
}

// DecodeAuthorizationRequest takes in a url with the query string parameters and converts it into a struct
func DecodeAuthorizationRequest(u url.Values) *AuthorizationRequest {
	return &AuthorizationRequest{
		ResponseType: u.Get("response_type"),
		ClientID:     u.Get("client_id"),
		RedirectURI:  u.Get("redirect_uri"),
		Scope:        u.Get("scope"),
		State:        u.Get("state"),
	}
}

// Encode converts the struct into url with query string
func EncodeAuthorizationRequest(r *AuthorizationRequest) (url.Values, error) {
	u, err := url.Parse("")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("response_type", r.ResponseType)
	q.Add("client_id", r.ClientID)
	q.Add("redirect_uri", r.RedirectURI)
	q.Add("scope", r.Scope)
	q.Add("state", r.State)

	return q, nil
}

// Validate performs validation on required fields
func (r *AuthorizationRequest) Validate() error {
	// Required fields
	if r.ResponseType != "code" {
		return ErrUnsupportedResponseType
	}
	if r.ClientID == "" {
	}
	// Optional fields
	if r.RedirectURI == "" {
	}
	if r.Scope == "" {
		return ErrInvalidScope
	}
	if r.State == "" {
	}
	return nil
}

type AuthorizationResponse struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func DecodeAuthorizationResponse(u url.Values) *AuthorizationResponse {
	return &AuthorizationResponse{
		Code:  u.Get("code"),
		State: u.Get("state"),
	}
}
func EncodeAuthorizationResponse(r *AuthorizationResponse, targetURL string) (*url.URL, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("code", r.Code)
	q.Add("state", r.State)
	u.RawQuery = q.Encode()
	return u, nil
}

type AuthorizationService interface {
	GenerateCode() string
}

func HandleAuthorizationRequest(s AuthorizationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := DecodeAuthorizationRequest(r.URL.Query())
		if err := req.Validate(); err != nil {
			log.Fatal(err)
		}

		res := &AuthorizationResponse{
			Code:  s.GenerateCode(),
			State: req.State,
		}
		redirectURL, err := EncodeAuthorizationResponse(res, req.RedirectURI)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
	}
}

type AuthorizationError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
	State            string `json:"state"`
}
