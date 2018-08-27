package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	oidc "github.com/alextanhongpin/go-openid"
)

func newMockEndpoint() *Endpoints {
	codeGen := func() string {
		return "code"
	}
	atGen := func() string {
		return "access_token"
	}
	rtGen := func() string {
		return "refresh_token"
	}
	return &Endpoints{
		service: NewService(nil, codeGen, atGen, rtGen),
	}
}

func setupAuthorizationRequest() url.Values {
	req := &oidc.AuthorizationRequest{
		ResponseType: "code",
		ClientID:     "1",
		RedirectURI:  "http://client/cb",
		Scope:        "profile",
		State:        "123",
	}
	q, _ := oidc.EncodeAuthorizationRequest(req)
	return q
}

func TestAuthorizeEndpoint(t *testing.T) {
	assert := assert.New(t)

	// Setup mock endpoint
	e := newMockEndpoint()

	// Setup router
	router := httprouter.New()
	router.GET("/authorize", e.Authorize)

	// Setup payload
	q := setupAuthorizationRequest()

	// Setup request
	req, _ := http.NewRequest("GET", "/authorize", nil)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	// Serve mock requests
	router.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(rr.Code, http.StatusFound, "handler return wrong status code")

	log.Println(rr.Body.String())

	u, _ := url.Parse(rr.Header().Get("Location"))
	res := oidc.DecodeAuthorizationResponse(u.Query())
	assert.Equal(res.Code, "code", "handler return wrong authorization code")
	assert.Equal(res.State, authReq.State, "handler return wrong state")
}

func setupTokenRequest() []byte {
	req := oidc.AccessTokenRequest{
		GrantType:   "authorization_code",
		Code:        "xyz",
		RedirectURI: "http://client/cb",
		ClientID:    "1234",
	}
	js, _ := json.Marshal(req)
	return js
}
func TestTokenEndpoint(t *testing.T) {
	assert := assert.New(t)

	e := newMockEndpoint()

	router := httprouter.New()
	router.POST("/token", e.Token)

	// Setup payload
	payload := setupTokenRequest()
	req, _ := http.NewRequest("POST", "/token", bytes.NewBuffer(payload))

	rr := httptest.NewRecorder()
	log.Println(rr.Body.String())

	// Setup mock requests
	router.ServeHTTP(rr, req)

	log.Println(rr.Body.String())
	assert.Equal(rr.Code, http.StatusOK, "handler return wrong status code")
	assert.Equal(rr.Header().Get("Cache-Control"), "no-store", "cache-control header set incorrectly")
	assert.Equal(rr.Header().Get("Pragma"), "no-cache", "pragma header set incorrectly")
}
