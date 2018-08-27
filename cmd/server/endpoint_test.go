package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	oidc "github.com/alextanhongpin/go-openid"
)

func newMockService(db *Database) *Service {
	codeGen := func() string {
		return "code"
	}
	atGen := func() string {
		return "access_token"
	}
	rtGen := func() string {
		return "refresh_token"
	}
	return NewService(db, codeGen, atGen, rtGen)
}

func newMockEndpoint(s *Service) *Endpoints {
	if s == nil {
		s = newMockService(nil)
	}
	return &Endpoints{
		service: s,
	}
}

func setupAuthorizationRequest() *oidc.AuthorizationRequest {
	return &oidc.AuthorizationRequest{
		ResponseType: "code",
		ClientID:     "1",
		RedirectURI:  "http://client/cb",
		Scope:        "profile",
		State:        "123",
	}
}

func TestAuthorizeEndpoint(t *testing.T) {
	assert := assert.New(t)

	// Setup mock endpoint
	db := NewDatabase()
	db.Client.Put("1", &oidc.Client{
		ClientRegistrationRequest: &oidc.ClientRegistrationRequest{
			ClientName:   "oidc_app",
			RedirectURIs: []string{"http://client/cb"},
		},
		ClientRegistrationResponse: &oidc.ClientRegistrationResponse{
			ClientID: "1",
		},
	})
	s := newMockService(db)
	e := newMockEndpoint(s)

	// Setup router
	router := httprouter.New()
	router.GET("/authorize", e.Authorize)

	// Setup payload
	authzReq := setupAuthorizationRequest()
	q, _ := oidc.EncodeAuthorizationRequest(authzReq)

	// Setup request
	req, _ := http.NewRequest("GET", "/authorize", nil)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	// Serve mock requests
	router.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(http.StatusFound, rr.Code, "handler return wrong status code")

	u, _ := url.Parse(rr.Header().Get("Location"))
	res := oidc.DecodeAuthorizationResponse(u.Query())

	assert.Equal("code", res.Code, "handler return wrong authorization code")
	assert.Equal(authzReq.State, res.State, "handler return wrong state")

	// TODO: Test the db to see if the data is stored
}

func setupTokenRequest() *oidc.AccessTokenRequest {
	return &oidc.AccessTokenRequest{
		GrantType:   "authorization_code",
		Code:        "xyz",
		RedirectURI: "http://client/cb",
		ClientID:    "1234",
	}
}

func TestTokenEndpoint(t *testing.T) {
	assert := assert.New(t)
	db := NewDatabase()
	db.Client.Put("oidc_app", &oidc.Client{
		ClientRegistrationRequest: &oidc.ClientRegistrationRequest{
			ClientName:   "oidc_app",
			RedirectURIs: []string{"http://client/cb"},
		},
		ClientRegistrationResponse: &oidc.ClientRegistrationResponse{
			ClientID: "1234",
		},
	})
	db.Code.Put("1234", oidc.NewCode("xyz"))
	s := newMockService(db)
	e := newMockEndpoint(s)

	router := httprouter.New()
	router.POST("/token", e.Token)

	// Setup payload
	tokenReq := setupTokenRequest()
	tokenJSON, _ := json.Marshal(tokenReq)
	req, _ := http.NewRequest("POST", "/token", bytes.NewBuffer(tokenJSON))

	rr := httptest.NewRecorder()

	// Setup mock requests
	router.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code, "handler return wrong status code")
	assert.Equal("no-store", rr.Header().Get("Cache-Control"), "cache-control header set incorrectly")
	assert.Equal("no-cache", rr.Header().Get("Pragma"), "pragma header set incorrectly")

	// TODO: Test the database to see if the data is stored in the storage
}

func TestClientRegistrationEndpoint(t *testing.T) {
	assert := assert.New(t)
	db := NewDatabase()
	s := newMockService(db)
	s.newClient = func(req *oidc.ClientRegistrationRequest) *oidc.Client {
		return &oidc.Client{
			ClientRegistrationRequest: req,
			ClientRegistrationResponse: &oidc.ClientRegistrationResponse{
				ClientID:                "test client id",
				ClientSecret:            "test client secret",
				RegistrationAccessToken: "test registration access token",
				RegistrationClientURI:   "test registration client uri",
				ClientIDIssuedAt:        1000,
				ClientSecretExpiresAt:   1000,
			},
		}
	}
	e := newMockEndpoint(s)

	router := httprouter.New()
	router.POST("/client/register", e.RegisterClient)

	// Setup payload
	clientReq := &oidc.ClientRegistrationRequest{
		ClientName: "oidc_app",
	}
	clientJSON, _ := json.Marshal(clientReq)
	req, _ := http.NewRequest("POST", "/client/register", bytes.NewBuffer(clientJSON))

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	var res oidc.ClientRegistrationResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	// Check status code
	assert.Equal(rr.Code, http.StatusCreated, "return incorrect status code")

	// Check headers
	assert.Equal("application/json", rr.Header().Get("Content-Type"), "return incorrect Content-Type")
	assert.Equal("no-store", rr.Header().Get("Cache-Control"), "return incorrect Cache-Control")
	assert.Equal("no-cache", rr.Header().Get("Pragma"), "return incorrect Pragma")

	// Check response
	assert.Equal("test client id", res.ClientID, "return incorrect client id")
	assert.Equal("test client secret", res.ClientSecret, "return incorrect client secret")
	assert.Equal("test registration access token", res.RegistrationAccessToken, "return incorrect registration access token")
	assert.Equal("test registration client uri", res.RegistrationClientURI, "return wrong registration client uri")
	assert.Equal(int64(1000), res.ClientIDIssuedAt, "return wrong issued date")
	assert.Equal(int64(1000), res.ClientSecretExpiresAt, "return wrong client secret expired at date")

	// Check the database to see if the client has been stored successfully
	clientdb, exist := db.Client.Get(clientReq.ClientName)
	assert.Equal(true, exist, "client does not exist in the storage")
	assert.Equal(res, *clientdb.ClientRegistrationResponse, "should point to the same object")
}
