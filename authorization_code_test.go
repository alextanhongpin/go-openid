package oidc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAuthorizationRequest(t *testing.T) {
	assert := assert.New(t)
	req := &AuthorizationRequest{
		ResponseType: "code",
		ClientID:     "abc",
		RedirectURI:  "http://client.example.com/cb",
		Scope:        "email",
		State:        "xyz",
	}
	u, err := EncodeAuthorizationRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("client_id=abc&redirect_uri=http%3A%2F%2Fclient.example.com%2Fcb&response_type=code&scope=email&state=xyz", u.Encode(), "should encode authorization request")
}

func TestDecodeAuthorizationRequest(t *testing.T) {
	assert := assert.New(t)
	u, err := url.Parse("http://server.example.com?client_id=abc&redirect_uri=http%3A%2F%2Fclient.example.com%2Fcb&response_type=code&scope=email&state=xyz")
	if err != nil {
		t.Fatal(err)
	}
	req := DecodeAuthorizationRequest(u.Query())
	assert.Equal("code", req.ResponseType, "should have field code")
	assert.Equal("abc", req.ClientID, "should have field client id")
	assert.Equal("http://client.example.com/cb", req.RedirectURI, "should have field redirect uri")
	assert.Equal("email", req.Scope, "should have field scope")
	assert.Equal("xyz", req.State, "should have field state")
}

func TestAuthorizationFlow(t *testing.T) {
	assert := assert.New(t)
	redirectURI := "https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz"
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectURI, http.StatusFound)
	}
	req := httptest.NewRequest("GET", "http://server.example.com/authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	res := w.Result()
	assert.Equal(http.StatusFound, res.StatusCode, "should return 302 - Status Found")
	assert.Equal(redirectURI, res.Header.Get("Location"), "should return the redirect uri in Location header")
}
func TestEncodeAuthorizationError(t *testing.T) {
	assert := assert.New(t)
	res := &AuthorizationError{
		Error: "access_denied",
		State: "xyz",
	}
	redirectURI, err := EncodeAuthorizationError(res, "https://client.example.com/cb")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("https://client.example.com/cb?error=access_denied&state=xyz", redirectURI.String(), "should encode the correct authorization error")

}
func TestAuthorizationErrorFlow(t *testing.T) {
	assert := assert.New(t)
	redirectURI := "https://client.example.com/cb?error=access_denied&state=xyz"
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectURI, http.StatusFound)
	}

	req := httptest.NewRequest("GET", "http://server.example.com/authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	res := w.Result()

	// Test header
	assert.Equal(http.StatusFound, res.StatusCode, "should return 302 - Status Found")
	assert.Equal(redirectURI, res.Header.Get("Location"), "should return the redirect uri in Location header")
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
