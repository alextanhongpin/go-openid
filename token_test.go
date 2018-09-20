package oidc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alextanhongpin/go-openid/pkg/querystring"

	"github.com/stretchr/testify/assert"
)

func TestTokenFlow(t *testing.T) {
	assert := assert.New(t)
	authorization := "Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW"
	tokenReq := AccessTokenRequest{
		GrantType:   "authorization_code",
		Code:        "SplxlOBeZQQYbYS6WxSbIA",
		RedirectURI: "https://client.example.com/cb",
	}
	tokenRes := AccessTokenResponse{
		AccessToken:  "2YotnFZFEjr1zCsicMWpAA",
		TokenType:    "example",
		ExpiresIn:    3600,
		RefreshToken: "tGzv3JOkF0XG5Qx2TlKWIA",
	}
	form := querystring.Encode(url.Values{}, tokenReq)
	req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		contentType := r.Header.Get("Content-Type")
		assert.Equal(authorization, auth)
		assert.Equal("application/x-www-form-urlencoded", contentType)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokenRes)
	}
	w := httptest.NewRecorder()
	handler(w, req)

	// res := w.Result()
	assert.Equal(http.StatusOK, w.Code, "should return status 200 - OK")
}

//
//     POST /token HTTP/1.1
//     Host: server.example.com
//     Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
//     Content-Type: application/x-www-form-urlencoded
//
//     grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA
//     &redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb
//
//
//
//     HTTP/1.1 200 OK
//     Content-Type: application/json;charset=UTF-8
//     Cache-Control: no-store
//     Pragma: no-cache
//
//     {
//       "access_token":"2YotnFZFEjr1zCsicMWpAA",
//       "token_type":"example",
//       "expires_in":3600,
//       "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
//       "example_parameter":"example_value"
//     }
