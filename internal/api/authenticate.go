package main

import (
	"errors"
	"strings"
)

type (
	// AuthenticateRequest represents the authentication request.
	AuthenticateRequest struct {
		AcrValues    string `json:"acr_values,omitempty"`
		ClientID     string `json:"client_id,omitempty"`
		Display      string `json:"display,omitempty"`
		IDTokenHint  string `json:"id_token_hint,omitempty"`
		LoginHint    string `json:"login_hint,omitempty"`
		MaxAge       int64  `json:"max_age,omitempty"`
		Nonce        string `json:"nonce,omitempty"`
		Prompt       string `json:"prompt,omitempty"`
		RedirectURI  string `json:"redirect_uri,omitempty"`
		ResponseMode string `json:"response_mode,omitempty"`
		ResponseType string `json:"response_type,omitempty"`
		Scope        string `json:"scope,omitempty"`
		State        string `json:"state,omitempty"`
		UILocales    string `json:"ui_locales,omitempty"`
	}

	// AuthenticateResponse represents the successfull authentication
	// response.
	AuthenticateResponse struct {
		Code  string `json:"code,omitempty"`
		State string `json:"state,omitempty"`
	}
)

type Scope string

func (s Scope) Has(scope string) bool {
	return strings.Contains(string(s), scope)
}

// NewAuthenticateRequest returns a new AuthenticateRequest with the response
// type code.
func NewAuthenticateRequest(scope, clientID, redirectURI string) *AuthenticateRequest {
	return &AuthenticateRequest{
		Scope:        scope,
		ResponseType: "code",
		ClientID:     clientID,
		RedirectURI:  redirectURI,
	}
}

// NewAuthenticateResponse returns a new AuthenticateResponse with a unique,
// code that is valid for a fixed duration.
func NewAuthenticateResponse(code, state string) *AuthenticateResponse {
	return &AuthenticateResponse{
		Code:  code,
		State: state,
	}
}

// NOTE: Full/Partial? Partial will not return any response. By default, it
// will always be treated as full, and hence does not need the suffix -Full.
// AuthenticatePartial if no response is required. If there are multiple steps,
// can add the suffix -Flow.

// AuthenticateFlow represents the authentication flow.
func Authenticate(
	repo ClientRepository,
	// Bad, there are two Code parameters here, we are dealing with
	// implmentation logic.
	// codeRepo CodeRepository,
	// codeFactory CodeFactory,
	codeInteractor *CodeInteractor,
	req *AuthenticateRequest,
) (*AuthenticateResponse, error) {
	// Validate the required fields. The logic is tied to the struct.
	if err := req.HasRequiredFields(); err != nil {
		return nil, err
	}
	// Can it be simplified to get client by clientID and redirectURI?
	// client, err := ValidateClient(repo, req.ClientID)
	client, err := repo.GetClientByClientID(req.ClientID)
	if err != nil {
		return nil, err
	}
	if !client.HasRedirectURI(req.RedirectURI) {
		return nil, errors.New("redirect_uri is invalid")
	}
	code, err := codeInteractor.NewCode()
	if err != nil {
		return nil, err
	}
	// Redundant, since we only have two fields.
	// res := NewAuthenticateResponse(code.ID, req.State)
	// return res, nil
	return &AuthenticateResponse{
		Code:  code.ID,
		State: req.State,
	}, nil
}

// How to avoid building an anemic model? Bind more logic to the struct. For
// entity, we can't really know for sure what data to validate, but for
// request/response struct they are pretty much consistent - we know what
// fields are required  or not.
// ValidateAuthenticateRequest checks for the required fields and the values
// set.
func (req *AuthenticateRequest) HasRequiredFields() error {
	// Validate required fields.
	fields := []struct {
		label, value string
	}{
		{"scope", req.Scope},
		{"response_type", req.ResponseType},
		{"client_id", req.ClientID},
		{"redirect_uri", req.RedirectURI},
	}
	for _, field := range fields {
		if stringIsEmpty(field.value) {
			return MakeErrRequired(field.label)
		}
	}
	// Validate specific fields.
	// Hide implementation details with structs.
	if !Scope(req.Scope).Has("openid") {
		return MakeErrRequired("openid")
	}
	if req.ResponseType != "code" {
		return errors.New(`response_type "code" is invalid`)
	}
	return nil
}
