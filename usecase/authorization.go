package usecase

import (
	"context"
	"fmt"
	"strings"

	openid "github.com/alextanhongpin/go-openid"
)

type Authorizer interface {
	Authorize(ctx context.Context, req AuthorizationRequest) (*AuthorizationResponse, error)
}

// AuthorizationRequest represents the request payload for authorization.
type AuthorizationRequest struct {
	ClientID     string             `json:"client_id,omitempty"`
	RedirectURI  openid.RedirectURI `json:"redirect_uri,omitempty"`
	ResponseType string             `json:"response_type,omitempty"`
	Scope        string             `json:"scope,omitempty"`
	State        string             `json:"state,omitempty"`
}

// GetResponseType returns the response type.
func (a *AuthorizationRequest) GetResponseType() ResponseType {
	return parseResponseType(a.ResponseType)
}

// Validate performs validation on required fields.
func (r *AuthorizationRequest) Validate() error {
	// Required fields
	if err := validateAuthzResponseType(r.ResponseType); err != nil {
		return err
	}
	if err := validateAuthzClientID(r.ClientID); err != nil {
		return err
	}
	if err := r.RedirectURI.Validate(); err != nil {
		return err
	}
	return nil
}

func validateAuthzResponseType(in string) error {
	parsed := parseResponseType(in)
	if !parsed.Is(ResponseTypeCode) {
		desc := fmt.Sprintf(`"%s" is not a valid response_type`, in)
		return ErrUnsupportedResponseType.WithDescription(desc)
	}
	return nil
}

func validateAuthzClientID(id string) error {
	// Is it necessary to check for empty string?
	// Probably, it will save us from making a call to the db.
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrAccessDenied.WithDescription("client_id cannot be empty")
	}
	return nil
}

// AuthorizationResponse represents the authorization response body.
type AuthorizationResponse struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}
