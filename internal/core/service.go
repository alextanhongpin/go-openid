package core

import (
	"context"
	"errors"
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/asaskevich/govalidator"
)

type serviceImpl struct {
	model model.Core
}

func (s *serviceImpl) Authenticate(ctx context.Context, req *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error) {
	userID, ok := ctx.Value(oidc.UserContextKey).(string)
	if !ok {
		return nil, errors.New("user_id not present")
	}
	user, err := s.model.GetUser(userID)
	if err != nil {
		return nil, err
	}

	if time.Since(time.Unix(user.Profile.UpdatedAt, 0)) > time.Duration(req.MaxAge) {
		// TODO: Must re-authenticate.
	}

	if req.LoginHint == user.Email.Email {
		// user.Email.Email
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

	if err := validateRequiredFields(req); err != nil {
		return nil, err
	}
	client, err := s.model.GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	if !client.RedirectURIs.Contains(req.RedirectURI) {
		return nil, errors.New("redirect_uri incorrect")
	}

	// Return the code to be exchanged for a token.
	code := s.model.NewCode()

	return &oidc.AuthenticationResponse{
		Code:  code,
		State: req.State,
	}, nil

}
