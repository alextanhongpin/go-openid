package core

import (
	"context"
	"errors"
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/model"
	"github.com/asaskevich/govalidator"
)

type serviceImpl struct {
	model model.Core
}

var ErrMustReauthenticate = errors.New("re-authentication required")

// PreAuthenticate will only check if the authentication request is valid and
// the type of flow it is using.
func (s *serviceImpl) PreAuthenticate(ctx context.Context, req *oidc.AuthenticationRequest) error {

	return nil
}

// Authenticate performs the full authentication and validation of all fields.
func (s *serviceImpl) Authenticate(ctx context.Context, req *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error) {
	validateUserFields := func(ctx context.Context, req *oidc.AuthenticationRequest) error {
		userID, ok := ctx.Value(oidc.UserContextKey).(string)
		if !ok {
			return errors.New("user_id not present")
		}
		user, err := s.model.GetUser(userID)
		if err != nil {
			return err
		}
		if time.Since(time.Unix(user.Profile.UpdatedAt, 0)) > time.Duration(req.MaxAge) {
			// TODO: Must re-authenticate.
			return ErrMustReauthenticate
		}

		if req.LoginHint == user.Email.Email {
			// user.Email.Email
		}
		return nil
	}

	validateRequiredFields := func(req *oidc.AuthenticationRequest) error {
		var (
			scope        = req.GetScope()
			responseType = req.GetResponseType()
			clientID     = req.ClientID
			redirectURI  = req.RedirectURI
			prompt       = req.GetPrompt()
		)
		if scope.Is(oidc.ScopeNone) || !scope.Has(oidc.ScopeOpenID) {
			return errors.New("scope required")
		}
		if responseType.Is(oidc.ResponseTypeNone) || !responseType.Has(oidc.ResponseTypeCode) {
			return errors.New("response_type required")
		}
		if clientID == "" {
			return errors.New("client_id required")
		}
		if redirectURI == "" {
			return errors.New("redirect_uri required")
		}
		if !govalidator.IsURL(redirectURI) {
			return errors.New("redirect_uri invalid")
		}
		if prompt.Has(oidc.PromptNone) && prompt.Has(oidc.PromptLogin|oidc.PromptConsent|oidc.PromptSelectAccount) {
			return errors.New("prompt none may not contain other values")
		}
		return nil
	}

	validateClientFields := func(req *oidc.AuthenticationRequest) error {
		// Fields to validate.
		var (
			clientID    = req.ClientID
			redirectURI = req.RedirectURI
		)
		client, err := s.model.GetClient(clientID)
		if err != nil {
			return err
		}
		if !client.GetRedirectURIs().Contains(redirectURI) {
			return errors.New("redirect_uri incorrect")
		}
		return nil
	}

	if err := validateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := validateUserFields(ctx, req); err != nil {
		return nil, err
	}
	if err := validateClientFields(req); err != nil {
		return nil, err
	}
	// Return the code to be exchanged for a token.
	code := s.model.NewCode()

	return &oidc.AuthenticationResponse{
		Code:  code,
		State: req.State,
	}, nil

}
