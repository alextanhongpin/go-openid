package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
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

func TestArgon2(t *testing.T) {
	key := argon2.IDKey([]byte("hello world"), []byte("hello"), 1, 64*1024, 4, 32)
	log.Println(crypto.Argon2id(key))

}
func TestAuthorizeEndpoint(t *testing.T) {
	assert := assert.New(t)

	// Setup mock endpoint
	e := newMockEndpoint()

	// Setup router
	router := httprouter.New()
	router.GET("/authorize", e.Authorize)

	// Setup payload
	authReq := &openid.AuthorizationRequest{
		ResponseType: "code",
		ClientID:     "1",
		RedirectURI:  "http://client/cb",
		Scope:        "profile",
		State:        "123",
	}
	rawQuery, err := openid.EncodeAuthorizationRequest(authReq)
	if err != nil {
		t.Fatal(err)
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/authorize", nil)
	req.URL.RawQuery = rawQuery.Encode()
	rr := httptest.NewRecorder()

	// Serve mock requests
	router.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(rr.Code, http.StatusFound, "handler return wrong status code")

	log.Println(rr.Body.String())

	u, err := url.Parse(rr.Header().Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	res := openid.DecodeAuthorizationResponse(u.Query())
	assert.Equal(res.Code, "code", "handler return wrong authorization code")
	assert.Equal(res.State, authReq.State, "handler return wrong state")
}

func TestTokenEndpoint(t *testing.T) {
	assert := assert.New(t)

	e := newMockEndpoint()

	router := httprouter.New()
	router.POST("/token", e.Token)

	// Setup payload
	atReq := openid.AccessTokenRequest{
		GrantType:   "authorization_code",
		Code:        "XYZ",
		RedirectURI: "http://client/cb",
		ClientID:    "123",
	}
	atReqJSON, err := json.Marshal(atReq)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/token", bytes.NewBuffer(atReqJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	log.Println(rr.Body.String())
	// Setup mock requests
	router.ServeHTTP(rr, req)

	log.Println(rr.Body.String())
	assert.Equal(rr.Code, http.StatusOK, "handler return wrong status code")
	assert.Equal(rr.Header().Get("Cache-Control"), "no-store", "cache-control header set incorrectly")
	assert.Equal(rr.Header().Get("Pragma"), "no-cache", "pragma header set incorrectly")
}
