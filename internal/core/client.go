package core

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/authheader"
)

type Client struct {
	TokenRegistrationURI string
	ClientID             string
	ClientSecret         string
}

// Exchange trades an authorization code with the access token.
func (c *Client) Exchange(ctx context.Context, code, redirectURI string) (*oidc.AuthenticationResponse, error) {
	tokenReq := oidc.AccessTokenRequest{
		GrantType:   "authorization_code",
		Code:        code,
		RedirectURI: redirectURI,
	}

	jsonBody, err := json.Marshal(tokenReq)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequest("POST", c.TokenRegistrationURI, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Add("Authorization", "Basic "+authheader.EncodeBase64(c.ClientID, c.ClientSecret))

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var res oidc.AuthenticationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
