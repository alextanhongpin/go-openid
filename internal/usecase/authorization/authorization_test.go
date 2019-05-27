package authorization_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/alextanhongpin/go-openid/pkg/querystring"
	"github.com/stretchr/testify/assert"
)

func TestEncodeAuthorizationRequest(t *testing.T) {
	assert := assert.New(t)
	req := &AuthorizationRequest{
		ClientID:     "abc",
		RedirectURI:  "http://client.example.com/cb",
		ResponseType: "code",
		Scope:        "email",
		State:        "xyz",
	}
	u := querystring.Encode(url.Values{}, req)
	assert.Equal("client_id=abc&redirect_uri=http%3A%2F%2Fclient.example.com%2Fcb&response_type=code&scope=email&state=xyz", u.Encode(), "should encode authorization request")
}

func TestDecodeAuthorizationRequest(t *testing.T) {
	assert := assert.New(t)
	u, err := url.Parse("http://server.example.com?client_id=abc&redirect_uri=http%3A%2F%2Fclient.example.com%2Fcb&response_type=code&scope=email&state=xyz")
	assert.Nil(err)

	var req AuthorizationRequest
	err = querystring.Decode(u.Query(), &req)
	assert.Nil(err)

	var (
		clientID     = "abc"
		redirectURI  = "http://client.example.com/cb"
		responseType = "code"
		scope        = "email"
		state        = "xyz"
	)

	assert.Equal(responseType, req.ResponseType, "should have field code")
	assert.Equal(clientID, req.ClientID, "should have field client id")
	assert.Equal(redirectURI, req.RedirectURI, "should have field redirect uri")
	assert.Equal(scope, req.Scope, "should have field scope")
	assert.Equal(state, req.State, "should have field state")
}

func TestAuthorizationFlow(t *testing.T) {
	assert := assert.New(t)

	var (
		redirectURI = "https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz"
		statusCode  = http.StatusFound
	)

	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectURI, statusCode)
	}

	req := httptest.NewRequest("GET", "http://server.example.com/authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb", nil)

	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()

	assert.Equal(statusCode, res.StatusCode, "should return 302 - Status Found")
	assert.Equal(redirectURI, res.Header.Get("Location"), "should return the redirect uri in Location header")
}

func TestEncodeAuthorizationError(t *testing.T) {
	assert := assert.New(t)

	res := &ErrorJSON{
		Code:  "access_denied",
		State: "xyz",
	}

	u, err := url.Parse("https://client.example.com/cb")
	assert.Nil(err)

	q := querystring.Encode(url.Values{}, res)
	u.RawQuery = q.Encode()

	assert.Equal("https://client.example.com/cb?error=access_denied&state=xyz", u.String(), "should encode the correct authorization error")
}

func TestAuthorizationErrorFlow(t *testing.T) {
	assert := assert.New(t)

	var (
		redirectURI = "https://client.example.com/cb?error=access_denied&state=xyz"
		statusCode  = http.StatusFound
	)

	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectURI, statusCode)
	}

	req := httptest.NewRequest("GET", "http://server.example.com/authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb", nil)

	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()

	// Test header
	assert.Equal(statusCode, res.StatusCode, "should return 302 - Status Found")
	assert.Equal(redirectURI, res.Header.Get("Location"), "should return the redirect uri in Location header")
}
