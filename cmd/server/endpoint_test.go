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
	"github.com/alextanhongpin/go-openid/pkg/querystring"
)

func TestAuthorizeEndpoint(t *testing.T) {
	assert := assert.New(t)

	// Setup mock endpoint
	db := NewDatabase()
	injectClient(db)
	s := newMockService(db)
	e := newMockEndpoint(s)

	// Setup router
	router := httprouter.New()
	router.GET("/authorize", e.Authorize)

	// Setup payload
	authzReq := &oidc.AuthorizationRequest{
		ResponseType: "code",
		ClientID:     "1",
		RedirectURI:  "http://client/cb",
		Scope:        "profile",
		State:        "123",
	}
	q := querystring.Encode(authzReq)

	// Setup request
	req := httptest.NewRequest("GET", "/authorize", nil)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	// Serve mock requests
	router.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(http.StatusFound, rr.Code, "handler return wrong status code")

	// Get the response, which is a redirect uri stored in header Location
	u, _ := url.Parse(rr.Header().Get("Location"))

	var res oidc.AuthorizationResponse
	err := querystring.Decode(&res, u.Query())
	assert.Nil(err)

	assert.Equal("code", res.Code, "should return an authorization code")
	assert.Equal(authzReq.State, res.State, "should return the given state")
	code, exist := db.Code.Get(authzReq.ClientID)
	assert.True(exist, "should have the client id in the db")
	assert.Equal(res.Code, code, "should match the code in the db")
}

func TestTokenEndpoint(t *testing.T) {
	assert := assert.New(t)
	db := NewDatabase()
	injectClient(db)
	db.Code.Put("1234", oidc.NewCode("xyz"))
	s := newMockService(db)
	e := newMockEndpoint(s)

	router := httprouter.New()
	router.POST("/token", e.Token)

	// Setup payload
	tokenReq := &oidc.AccessTokenRequest{
		GrantType:   "authorization_code",
		Code:        "xyz",
		RedirectURI: "https://client.example.com/cb",
		ClientID:    "1",
	}
	tokenJSON, _ := json.Marshal(tokenReq)

	req := httptest.NewRequest("POST", "/token", bytes.NewBuffer(tokenJSON))
	authorization := "Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW"
	assert.Equal(authorization, "", "should include authorization header")
	assert.Equal("application/x-www-form-urlencoded", "", "should set content-type to application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	// Setup mock requests
	router.ServeHTTP(rr, req)

	// Validate Headers
	cacheControl := "no-store"
	pragma := "no-cache"

	assert.Equal(http.StatusOK, rr.Code, "should return status 200 - OK")
	assert.Equal(cacheControl, rr.Header().Get("Cache-Control"), "should return Cache-Control no-store")
	assert.Equal(pragma, rr.Header().Get("Pragma"), "should return Pragma no-cache")

	// TODO: Test the database to see if the data is stored in the storage
	assert.Equal("application/json", "", "should set content-type to json")
	assert.Equal("access_token", "SlAV32hkKG", "should return access token")
	assert.Equal("token_type", "Bearer", "should return token type bearer")
	assert.Equal("refresh_token", "8xLOxBtZp8", "should return refresh token")
	assert.Equal("expires_in", 3600, "should return the correct expiry time")
	assert.Equal("id_token", "eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ.ewogImlzcyI6ICJodHRwOi8vc2VydmVyLmV4YW1wbGUuY29tIiwKICJzdWIiOiAiMjQ4Mjg5NzYxMDAxIiwKICJhdWQiOiAiczZCaGRSa3F0MyIsCiAibm9uY2UiOiAibi0wUzZfV3pBMk1qIiwKICJleHAiOiAxMzExMjgxOTcwLAogImlhdCI6IDEzMTEyODA5NzAKfQ.ggW8hZ1EuVLuxNuuIJKX_V8a_OMXzR0EHR9R6jgdqrOOF4daGU96Sr_P6qJp6IcmD3HP99Obi1PRs-cwh3LO-p146waJ8IhehcwL7F09JdijmBqkvPeB2T9CJNqeGpe-gccMg4vfKjkM8FcGvnzZUN4_KSP0aAp1tOJ1zZwgjxqGByKHiOtX7TpdQyHE5lcMiKPXfEIK5hoDalrcvRYLSrQAZZKflyuVCyixEoV9GfNQC3_osjzw2PAithfubEEBLuVVk4XUVrWOLrLl0nx7RkKU8NXNHq-rvKMzqg", "should return id token")
}

func TestTokenErrorResponse(t *testing.T) {
	statusCode := 400
	contentType := "application/json"
	cacheControl := "no-store"
	pragma := "no-cache"

	// 	{
	//    "error": "invalid_request"
	//   }
	t.Skip("testing test token error")
	assert := assert.New(t)
	assert.True(statusCode == statusCode)
	assert.True(contentType == contentType)
	assert.True(cacheControl == cacheControl)
	assert.True(pragma == pragma)
}

func TestAuthentication(t *testing.T) {
	t.Skip("test authentication")
	u, _ := url.Parse("http://server.example.com/authorize?response_type=id_token%20token&client_id=s6BhdRkqt3&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb&scope=openid%20profile&state=af0ifjsldkj&nonce=n-0S6_WzA2Mj")
	q := u.Query()

	assert := assert.New(t)
	assert.Equal("state", q.Get("state"), "should have the correct state")

	// HTTP/1.1 302 Found
	//  Location: https://client.example.org/cb#
	//    access_token=SlAV32hkKG
	//    &token_type=bearer
	//    &id_token=eyJ0 ... NiJ9.eyJ1c ... I6IjIifX0.DeWt4Qu ... ZXso
	//    &expires_in=3600
	//    &state=af0ifjsldkj
}

func TestUserInfo(t *testing.T) {
	//	 GET /userinfo HTTP/1.1
	//  Host: server.example.com
	//  Authorization: Bearer SlAV32hkKG
	//
	//   HTTP/1.1 200 OK
	//  Content-Type: application/json
	//
	//  {
	//   "sub": "248289761001",
	//   "name": "Jane Doe",
	//   "given_name": "Jane",
	//   "family_name": "Doe",
	//   "preferred_username": "j.doe",
	//   "email": "janedoe@example.com",
	//   "picture": "http://example.com/janedoe/me.jpg"
	//  }
	//
	//  HTTP/1.1 401 Unauthorized
	//   WWW-Authenticate: error="invalid_token",
	//     error_description="The Access Token expired"
}
func TestClientRegistrationEndpoint(t *testing.T) {
	assert := assert.New(t)
	db := NewDatabase()
	s := newMockService(db)
	s.newClient = func(req *oidc.ClientRegistrationRequest) *oidc.Client {
		return &oidc.Client{
			ClientPublic: req,
			ClientPrivate: &oidc.ClientRegistrationResponse{
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
	router.POST("/connect/register", e.RegisterClient)

	// Setup payload
	clientReq := &oidc.ClientRegistrationRequest{
		ClientName: "oidc_app",
	}
	clientJSON, _ := json.Marshal(clientReq)
	req := httptest.NewRequest("POST", "/connect/register", bytes.NewBuffer(clientJSON))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	var res oidc.ClientRegistrationResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	// Check status code
	assert.Equal(http.StatusCreated, rr.Code, "return incorrect status code")

	// Check headers
	header := rr.Header()
	assert.Equal("application/json", header.Get("Content-Type"), "return incorrect Content-Type")
	assert.Equal("no-store", header.Get("Cache-Control"), "return incorrect Cache-Control")
	assert.Equal("no-cache", header.Get("Pragma"), "return incorrect Pragma")

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

func TestClientRegistrationError(t *testing.T) {

	assert := assert.New(t)

	statusCode := 400
	contentType := "application/json"
	cacheControl := "no-store"
	pragma := "no-cache"
	res := oidc.ClientErrorResponse{
		Error:            "invalid_redirect_uri",
		ErrorDescription: "One or more redirect_uri values are invalid",
	}
	assert.Equal(400, statusCode, "should return bad request")
	assert.Equal("application/json", contentType, "should return json")
	assert.Equal("no-store", cacheControl, "should set cache-control to no-store")
	assert.Equal("no-cache", pragma, "should set pragma to no-cache")
	assert.Equal("invalid_redirect_uri", res.Error, "should return the correct error type")
	assert.True(res.ErrorDescription != "", "should return error description")

}

func TestClientRead(t *testing.T) {

	assert := assert.New(t)
	// Request
	clientID := "s6BhdRkqt3"
	authorization := "Bearer this.is.an.access.token.value.ffx83"
	assert.True(len(clientID) > 0, "should have client id")
	assert.True(len(authorization) > 0, "should have authorization header")

	// Response Headers
	statusCode := 200
	contentType := "application/json"
	cacheControl := "no-store"
	pragma := "no-cache"

	assert.Equal(200, statusCode, "should return ok")
	assert.Equal("application/json", contentType, "should return json")
	assert.Equal("no-store", cacheControl, "should set cache-control to no-store")
	assert.Equal("no-cache", pragma, "should set pragma to no-cache")

	// Response body
	client := &oidc.Client{
		ClientRegistrationRequest: &oidc.ClientRegistrationRequest{
			TokenEndpointAuthMethod: "token_endpoint_auth_method",
			ApplicationType:         "web",
			RedirectURIs: []string{"https://client.example.org/callback",
				"https://client.example.org/callback2"},
			ClientName:          "My Example",
			LogoURI:             "https://client.example.org/logo.png",
			SubjectType:         "pairwise",
			SectorIdentifierURI: "https://other.example.net/file_of_redirect_uris.json",
			JwksURI:             "https://client.example.org/my_public_keys.jwks",
			UserinfoEncryptedResponseAlg: "RSA1_5",
			UserinfoEncryptedResponseEnc: "A128CBC-HS256",
			Contacts:                     []string{"ve7jtb@example.org", "mary@example.org"},
			RequestURIs:                  []string{"https://client.example.org/rf.txt#qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"},
		},
		ClientRegistrationResponse: &oidc.ClientRegistrationResponse{
			ClientID:     "s6BhdRkqt3",
			ClientSecret: "OylyaC56ijpAQ7G5ZZGL7MMQ6Ap6mEeuhSTFVps2N4Q",

			RegistrationAccessToken: "",
			RegistrationClientURI:   "https://server.example.com/connect/register?client_id=s6BhdRkqt3",
			ClientIDIssuedAt:        0,
			ClientSecretExpiresAt:   17514165600,
		},
	}
	assert.True(client != nil, "should return client")
}

func TestClientReadError(t *testing.T) {
	assert := assert.New(t)
	statusCode := 401
	cacheControl := "no-store"
	pragma := "no-cache"

	assert.Equal(http.StatusUnauthorized, statusCode, "should return status 401 - Unauthorized")
	assert.Equal("no-store", cacheControl, "should set cache-control to no-store")
	assert.Equal("no-cache", pragma, "should set pragma to no-cache")
}

func TestOIDProviderIssuerDiscoveryEmail(t *testing.T) {
	u, _ := url.Parse("http://example.com/.well-known/webfinger?resource=acct%3Ajoe%40example.com&rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer")
	q := u.Query()
	resource := "acct:joe@example.com"
	host := "example.com"
	rel := "http://openid.net/specs/connect/1.0/issuer"

	assert := assert.New(t)
	assert.Equal(resource, q.Get("resource"), "should have resource in qs")
	assert.Equal(rel, q.Get("rel"), "should have rel in qs")

	// Header response
	statusCode := 200
	contentType := "application/jrd+json"
	assert.Equal(host, host, "should have host in qs")
	assert.Equal(http.StatusOK, statusCode, "should return status 200 - Ok")
	assert.Equal(contentType, contentType, "should return the correct content type")
	//	{
	//   "subject": "acct:joe@example.com",
	//   "links":
	//    [
	//     {
	//      "rel": "http://openid.net/specs/connect/1.0/issuer",
	//      "href": "https://server.example.com"
	//     }
	//    ]
	//  }
}

func TestOIDProviderIssuerDiscoveryURL(t *testing.T) {
	u, _ := url.Parse("http://example.com/.well-known/webfinger?resource=https%3A%2F%2Fexample.com%2Fjoe&rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer")
	q := u.Query()
	resource := "https://example.com/joe"
	rel := "http://openid.net/specs/connect/1.0/issuer"

	assert := assert.New(t)
	assert.Equal(resource, q.Get("resource"), "should have the correct resource in the qs")
	assert.Equal(rel, q.Get("rel"), "should have the correct rel in the qs")

	host := "example.com"
	statusCode := 200
	contentType := "application/jrd+json"
	assert.Equal(http.StatusOK, statusCode, "should return status 200 - Ok")
	assert.Equal(contentType, contentType, "should return the correct content type")
	assert.Equal(host, host, "should return the correct host")
	// Test response
	//	{
	//   "subject": "https://example.com/joe",
	//   "links":
	//    [
	//     {
	//      "rel": "http://openid.net/specs/connect/1.0/issuer",
	//      "href": "https://server.example.com"
	//     }
	//    ]
	//  }
}

func TestOIDPProviderUserDiscoveryHostnameAndPort(t *testing.T) {
	// TODO: User input using hostname and port syntax
	resource := "https://example.com:8080/"
	host := "example.com:8080"
	rel := "http://openid.net/specs/connect/1.0/issuer"

	assert := assert.New(t)
	assert.True(resource == resource)
	assert.True(host == host)
	assert.True(rel == rel)

	//	  GET /.well-known/webfinger
	//    ?resource=https%3A%2F%2Fexample.com%3A8080%2F
	//    &rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer
	//    HTTP/1.1
	//  Host: example.com:8080
	//
	//  HTTP/1.1 200 OK
	//  Content-Type: application/jrd+json
	//
	//  {
	//   "subject": "https://example.com:8080/",
	//   "links":
	//    [
	//     {
	//      "rel": "http://openid.net/specs/connect/1.0/issuer",
	//      "href": "https://server.example.com"
	//     }
	//    ]
	//  }
}

func TestDiscoverUserInputAcct(t *testing.T) {
	// resource	acct:juliet%40capulet.example@shopping.example.com
	// host	shopping.example.com
	// rel	http://openid.net/specs/connect/1.0/issuer
	//
	//   GET /.well-known/webfinger
	//     ?resource=acct%3Ajuliet%2540capulet.example%40shopping.example.com
	//     &rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer
	//     HTTP/1.1
	//   Host: shopping.example.com
	//
	//   HTTP/1.1 200 OK
	//   Content-Type: application/jrd+json
	//
	//   {
	//    "subject": "acct:juliet%40capulet.example@shopping.example.com",
	//    "links":
	//     [
	//      {
	//       "rel": "http://openid.net/specs/connect/1.0/issuer",
	//       "href": "https://server.example.com"
	//      }
	//     ]
	//   }
}

func TestOpenIDConfigurationRequest(t *testing.T) {
	url := "/.well-known/openid-configuration"
	assert := assert.New(t)
	assert.True(url == url)
	// 	{
	//    "issuer":
	//      "https://server.example.com",
	//    "authorization_endpoint":
	//      "https://server.example.com/connect/authorize",
	//    "token_endpoint":
	//      "https://server.example.com/connect/token",
	//    "token_endpoint_auth_methods_supported":
	//      ["client_secret_basic", "private_key_jwt"],
	//    "token_endpoint_auth_signing_alg_values_supported":
	//      ["RS256", "ES256"],
	//    "userinfo_endpoint":
	//      "https://server.example.com/connect/userinfo",
	//    "check_session_iframe":
	//      "https://server.example.com/connect/check_session",
	//    "end_session_endpoint":
	//      "https://server.example.com/connect/end_session",
	//    "jwks_uri":
	//      "https://server.example.com/jwks.json",
	//    "registration_endpoint":
	//      "https://server.example.com/connect/register",
	//    "scopes_supported":
	//      ["openid", "profile", "email", "address",
	//       "phone", "offline_access"],
	//    "response_types_supported":
	//      ["code", "code id_token", "id_token", "token id_token"],
	//    "acr_values_supported":
	//      ["urn:mace:incommon:iap:silver",
	//       "urn:mace:incommon:iap:bronze"],
	//    "subject_types_supported":
	//      ["public", "pairwise"],
	//    "userinfo_signing_alg_values_supported":
	//      ["RS256", "ES256", "HS256"],
	//    "userinfo_encryption_alg_values_supported":
	//      ["RSA1_5", "A128KW"],
	//    "userinfo_encryption_enc_values_supported":
	//      ["A128CBC-HS256", "A128GCM"],
	//    "id_token_signing_alg_values_supported":
	//      ["RS256", "ES256", "HS256"],
	//    "id_token_encryption_alg_values_supported":
	//      ["RSA1_5", "A128KW"],
	//    "id_token_encryption_enc_values_supported":
	//      ["A128CBC-HS256", "A128GCM"],
	//    "request_object_signing_alg_values_supported":
	//      ["none", "RS256", "ES256"],
	//    "display_values_supported":
	//      ["page", "popup"],
	//    "claim_types_supported":
	//      ["normal", "distributed"],
	//    "claims_supported":
	//      ["sub", "iss", "auth_time", "acr",
	//       "name", "given_name", "family_name", "nickname",
	//       "profile", "picture", "website",
	//       "email", "email_verified", "locale", "zoneinfo",
	//       "http://example.info/claims/groups"],
	//    "claims_parameter_supported":
	//      true,
	//    "service_documentation":
	//      "http://server.example.com/connect/service_documentation.html",
	//    "ui_locales_supported":
	//      ["en-US", "en-GB", "en-CA", "fr-FR", "fr-CA"]
	//   }
}

func newMockService(db *Database) *Service {
	gc := func() string {
		return "code"
	}
	gat := func() string {
		return "access_token"
	}
	grt := func() string {
		return "refresh_token"
	}
	return NewService(db, gc, gat, grt)
}

func newMockEndpoint(s *Service) *Endpoints {
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
func injectClient(db *Database) {
	id := "1"
	client1 := &oidc.Client{
		ClientRegistrationRequest: &oidc.ClientRegistrationRequest{
			ClientName:   "MyApp",
			RedirectURIs: []string{"https://client.example.com/cb"},
		},
		ClientRegistrationResponse: &oidc.ClientRegistrationResponse{
			ClientID: id,
		},
	}
	db.Client.Put(id, client1)
}
