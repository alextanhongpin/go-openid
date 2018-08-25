package openid

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var SigningKey = []byte("JWT_SECRET")
var ErrForbidden = errors.New("forbidden request")

type AccessTokenRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
	ClientID    string `json:"client_id"`
}

func (r *AccessTokenRequest) Validate() error {
	if r.GrantType != "authorization_code" {
		return ErrInvalidRequest
	}
	// Check required field
	if r.Code == "" {
		return ErrForbidden
	}

	if r.RedirectURI == "" {
		return ErrForbidden
	}
	if r.ClientID == "" {
		return ErrForbidden
	}
	return nil
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
type AccessTokenService interface {
	Do(AccessTokenRequest) (*AccessTokenResponse, error)
}

func HandleAccessTokenRequest(s OAuthService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		authHdr := r.Header.Get("Authorization")
		log.Println("authHdr", authHdr)
		var req AccessTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res, err := s.RequestAccessToken(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		json.NewEncoder(w).Encode(res)
	}
}
