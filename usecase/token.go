package usecase

import (
	"context"
	"errors"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/gostrings"
)

type Tokenizer interface {
	Token(ctx context.Context, req TokenRequest) (*TokenResponse, error)
}

type (
	TokenRequest struct {
		GrantType   string             `json:"grant_type"`
		Code        string             `json:"code"`
		RedirectURI openid.RedirectURI `json:"redirect_uri"`
	}
	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		IDToken      string `json:"id_token"`
	}
)

func (req *TokenRequest) Validate() error {
	// Validate required fields.
	if gostrings.IsEmpty(req.Code) {
		return errors.New("code is required")
	}
	if gostrings.IsEmpty(req.GrantType) {
		return errors.New("grant_type is required")
	}
	if gostrings.IsEmpty(req.RedirectURI) {
		return errors.New("redirect_uri is required")
	}
	if err := req.RedirectURI.Validate(); err != nil {
		return err
	}
	return nil
}
