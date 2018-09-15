package tokensvc

import (
	"time"

	"github.com/alextanhongpin/go-openid/models"
)

type postTokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
type postTokenResponse struct {
	AccessToken  string    `json:"access_token"`  // "SlAV32hkKG"
	TokenType    string    `json:"token_type"`    // "Bearer"
	RefreshToken string    `json:"refresh_token"` // "tGzv3JOkF0XG5Qx2TlKWIA"
	ExpiresIn    time.Time `json:"expires_in"`    // "2017-07-15T12:05:35.688266521+08:00"
	IDToken      string    `json:"id_token"`      // "eyJ0 ... NiJ9.eyJ1c ... I6IjIifX0.DeWt4Qu ... ZXso"
}
type clientRequest struct {
	ID string
}
type clientResponse struct {
	Data models.Client
}

type codeRequest struct {
	ClientID string
	Code     string
}
type codeResponse struct {
	Exist bool
}
