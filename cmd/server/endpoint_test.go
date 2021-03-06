package main_test

//
// import (
//         "bytes"
//         "encoding/json"
//         "net/http"
//         "net/http/httptest"
//         "net/url"
//         "strings"
//         "testing"
//         "time"
//
//         "github.com/julienschmidt/httprouter"
//         "github.com/stretchr/testify/assert"
//
//         openid "github.com/alextanhongpin/go-openid"
//         "github.com/alextanhongpin/go-openid/pkg/crypto"
//         "github.com/alextanhongpin/go-openid/pkg/querystring"
// )
//
// var (
//         // Crypto defaults
//         defaultXIDToken = "x2y9aS"
//         defaultJWTToken = "SlAV32hkKG"
//         defaultUUID     = "0000-0000-0000-0000"
//
//         // Request defaults
//         defaultResponseType = "code"
//         defaultClientID     = defaultXIDToken
//         defaultRedirectURI  = "https://client.example.com/cb"
//         defaultScope        = "profile email"
//         defaultState        = "xyz"
//
//         // Client defaults
//         defaultClientName   = "MyApp"
//         defaultClientSecret = defaultUUID
//
//         // Token defaults
//         defaultCode         = defaultXIDToken
//         defaultAccessToken  = defaultJWTToken
//         defaultRefreshToken = defaultJWTToken
//         defaultIDToken      = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ"
//
//         // User defaults
//         defaultUserID = "1"
// )
//
// func testAuthzEndpoint(e *Endpoints, r *openid.AuthorizationRequest) *httptest.ResponseRecorder {
//         router := httprouter.New()
//         router.GET("/authorize", e.Authorize)
//
//         q := querystring.Encode(r)
//
//         req := httptest.NewRequest("GET", "http://client.example.com/authorize", nil)
//         req.URL.RawQuery = q.Encode()
//
//         rr := httptest.NewRecorder()
//         router.ServeHTTP(rr, req)
//         return rr
// }
//
// func TestAuthorizeEndpoint(t *testing.T) {
//         assert := assert.New(t)
//
//         db := newMockDatabase()
//         s := newMockService(db)
//         e := newMockEndpoint(s)
//
//         // Setup payload
//         req := &openid.AuthorizationRequest{
//                 ResponseType: "code",
//                 ClientID:     defaultClientID,
//                 RedirectURI:  "https://client.example.com/cb",
//                 Scope:        "profile email",
//                 State:        "xyz",
//         }
//
//         rr := testAuthzEndpoint(e, req)
//
//         // Check status code
//         assert.Equal(http.StatusFound, rr.Code, "handler return wrong status code")
//
//         // Get the response, which is a redirect uri stored in header Location
//         u, _ := url.Parse(rr.Header().Get("Location"))
//
//         var res openid.AuthorizationResponse
//         err := querystring.Decode(&res, u.Query())
//         assert.Nil(err)
//
//         var (
//                 code  = defaultCode
//                 state = req.State
//         )
//
//         assert.Equal(code, res.Code, "should return an authorization code")
//         assert.Equal(state, res.State, "should return the given state")
//
//         codedb, exist := db.Code.Get(req.ClientID)
//         assert.True(exist, "should have the client id in the db")
//         assert.Equal(code, codedb.Code, "should match the code in the db")
// }
//
// func testTokenEndpoint(e *Endpoints, r *openid.AccessTokenRequest, bearer string) *httptest.ResponseRecorder {
//         router := httprouter.New()
//         router.POST("/token", e.Token)
//
//         form := querystring.Encode(r)
//
//         req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
//         req.Header.Set("Authorization", "Basic "+bearer)
//         req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//         rr := httptest.NewRecorder()
//         router.ServeHTTP(rr, req)
//         return rr
// }
//
// func TestTokenEndpoint(t *testing.T) {
//         assert := assert.New(t)
//
//         e := defaultMockEndpoint()
//
//         // Setup payload
//         req := &openid.AccessTokenRequest{
//                 GrantType:   "authorization_code",
//                 Code:        defaultCode,
//                 RedirectURI: defaultRedirectURI,
//                 ClientID:    defaultClientID,
//         }
//
//         var (
//                 bearer = openid.EncodeBasicAuth(defaultClientID, defaultClientSecret)
//         )
//         rr := testTokenEndpoint(e, req, bearer)
//
//         // Test response headers
//         var (
//                 cacheControl = "no-store"
//                 contentType  = "application/json"
//                 pragma       = "no-cache"
//                 statusCode   = 200
//         )
//
//         assert.Equal(statusCode, rr.Code, "should return status 200 - OK")
//         assert.Equal(contentType, rr.Header().Get("content-type"), "should have Content-Type application/json")
//         assert.Equal(cacheControl, rr.Header().Get("Cache-Control"), "should return Cache-Control no-store")
//         assert.Equal(pragma, rr.Header().Get("Pragma"), "should return Pragma no-cache")
//
//         // Test response body
//         var (
//                 accessToken  = defaultAccessToken
//                 tokenType    = "Bearer"
//                 refreshToken = defaultRefreshToken
//                 expiresIn    = int64(3600)
//                 idToken      = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ..."
//         )
//
//         var res openid.AccessTokenResponse
//         err := json.NewDecoder(rr.Body).Decode(&res)
//         assert.Nil(err)
//
//         assert.Equal(accessToken, res.AccessToken, "should return access token")
//         assert.Equal(tokenType, res.TokenType, "should return token type bearer")
//         assert.Equal(refreshToken, res.RefreshToken, "should return refresh token")
//         assert.Equal(expiresIn, res.ExpiresIn, "should return the correct expiry time")
//         assert.Equal(idToken, res.IDToken, "should return the id token")
// }
//
// func TestTokenErrorResponse(t *testing.T) {
//         assert := assert.New(t)
//
//         e := defaultMockEndpoint()
//         // Setup payload
//         req := &openid.AccessTokenRequest{}
//
//         rr := testTokenEndpoint(e, req, "")
//
//         // Validate headers
//         var (
//                 statusCode   = 403
//                 contentType  = "application/json"
//                 cacheControl = "no-store"
//                 pragma       = "no-cache"
//
//                 header = rr.Header()
//         )
//
//         assert.Equal(statusCode, rr.Code, "should return status 403 - Forbidden")
//         assert.Equal(contentType, header.Get("Content-Type"), "should return Content-Type application/json")
//         assert.Equal(cacheControl, header.Get("Cache-Control"), "should return Cache-Control no-store")
//         assert.Equal(pragma, header.Get("Pragma"), "should return Pragma no-cache")
//
//         var (
//                 errorCode = "invalid_request"
//         )
//
//         // Validate body
//         var res openid.ErrorJSON
//         err := json.NewDecoder(rr.Body).Decode(&res)
//         assert.Nil(err)
//
//         assert.Equal(errorCode, res.Code, "should return field error")
//         assert.True(res.Description != "", "should return field error description")
// }
//
// func TestAuthentication(t *testing.T) {
//         t.Skip("test authentication")
//         u, _ := url.Parse("http://server.example.com/authorize?response_type=id_token%20token&client_id=s6BhdRkqt3&redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb&scope=openid%20profile&state=af0ifjsldkj&nonce=n-0S6_WzA2Mj")
//         q := u.Query()
//
//         assert := assert.New(t)
//         assert.Equal("state", q.Get("state"), "should have the correct state")
//
//         // HTTP/1.1 302 Found
//         //  Location: https://client.example.org/cb#
//         //    access_token=SlAV32hkKG
//         //    &token_type=bearer
//         //    &id_token=eyJ0 ... NiJ9.eyJ1c ... I6IjIifX0.DeWt4Qu ... ZXso
//         //    &expires_in=3600
//         //    &state=af0ifjsldkj
// }
//
// func testUserInfo(e *Endpoints, id string) *httptest.ResponseRecorder {
//         router := httprouter.New()
//         router.GET("/userinfo", e.UserInfo)
//
//         req := httptest.NewRequest("GET", "/userinfo", nil)
//         req.Header.Set("Authorization", "Bearer slav32hkkg")
//
//         rr := httptest.NewRecorder()
//         router.ServeHTTP(rr, req)
//
//         return rr
// }
//
// func TestUserInfo(t *testing.T) {
//         assert := assert.New(t)
//
//         e := defaultMockEndpoint()
//
//         rr := testUserInfo(e, "1")
//
//         var (
//                 statusCode = 200
//         )
//
//         assert.Equal(statusCode, rr.Code, "should return status 200 - OK")
//
//         var (
//                 sub               = "248289761001"
//                 name              = "Jane Doe"
//                 givenName         = "Jane"
//                 familyName        = "Doe"
//                 preferredUsername = "j.doe"
//                 email             = "janedoe@example.com"
//                 picture           = "http://example.com/janedoe/me.jpg"
//         )
//
//         var u User
//         err := json.NewDecoder(rr.Body).Decode(&u)
//         assert.Nil(err)
//
//         assert.Equal(sub, u.Profile.Sub, "should match the subject")
//         assert.Equal(name, u.Profile.Name, "should match the name")
//         assert.Equal(givenName, u.Profile.GivenName, "should match the given name")
//         assert.Equal(familyName, u.Profile.FamilyName, "should match the family name")
//         assert.Equal(preferredUsername, u.Profile.PreferredUsername, "should match the preferred username")
//         assert.Equal(email, u.Email.Email, "should match the email")
//         assert.Equal(picture, u.Profile.Picture, "should match the picture")
// }
//
// func TestUserInfoError(t *testing.T) {
//         assert := assert.New(t)
//
//         e := defaultMockEndpoint()
//
//         rr := testUserInfo(e, "0")
//
//         var (
//                 statusCode = 401
//                 header     = `WWW-Authenticate: error="invalid_token" error_description="The access token expired"`
//         )
//
//         assert.Equal(statusCode, rr.Code, "should return status 401 - Unauthorized")
//         assert.Equal(header, rr.Header().Get("WWW-Authenticate"), "should return WWW-Authenticate response header with error message")
// }
//
// func testClientRegistration(e *Endpoints, r *openid.ClientRegistrationRequest, bearer string) *httptest.ResponseRecorder {
//         router := httprouter.New()
//         router.POST("/connect/register", e.RegisterClient)
//
//         reqJSON, _ := json.Marshal(r)
//         req := httptest.NewRequest("POST", "/connect/register", bytes.NewBuffer(reqJSON))
//         req.Header.Add("Authorization", "Bearer "+bearer)
//
//         rr := httptest.NewRecorder()
//
//         router.ServeHTTP(rr, req)
//
//         return rr
// }
//
// func TestClientRegistrationEndpoint(t *testing.T) {
//         assert := assert.New(t)
//
//         db := newMockDatabase()
//         s := newMockService(db)
//         e := newMockEndpoint(s)
//
//         // Setup payload
//         req := &openid.ClientPublic{
//                 ClientName: defaultClientName,
//         }
//
//         var (
//                 aud = "audience"
//                 sub = "subject"
//                 iss = "go-openid"
//                 dur = time.Minute
//                 rc  = crypto.New("secret")
//         )
//
//         bearer, err := rc.NewJWT(aud, sub, iss, dur)
//         assert.Nil(err)
//
//         rr := testClientRegistration(e, req, bearer)
//
//         var (
//                 statusCode = http.StatusCreated
//
//                 contentType  = "application/json"
//                 cacheControl = "no-store"
//                 pragma       = "no-cache"
//
//                 header = rr.Header()
//         )
//
//         // Check status code
//         assert.Equal(statusCode, rr.Code, "should return status code 201 - Created")
//
//         // Check headers
//         assert.Equal(contentType, header.Get("Content-Type"), "return incorrect Content-Type")
//         assert.Equal(cacheControl, header.Get("Cache-Control"), "return incorrect Cache-Control")
//         assert.Equal(pragma, header.Get("Pragma"), "return incorrect Pragma")
//
//         var (
//                 clientID                = defaultClientID
//                 clientSecret            = defaultClientSecret
//                 registrationAccessToken = defaultAccessToken
//                 registrationClientURI   = ""
//         )
//
//         var res openid.ClientPrivate
//         if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
//                 t.Fatal(err)
//         }
//
//         // Check response
//         assert.Equal(clientID, res.ClientID, "return incorrect client id")
//         assert.Equal(clientSecret, res.ClientSecret, "return incorrect client secret")
//         assert.Equal(registrationAccessToken, res.RegistrationAccessToken, "return incorrect registration access token")
//         assert.Equal(registrationClientURI, res.RegistrationClientURI, "return wrong registration client uri")
//
//         var (
//                 expireAt = time.Unix(res.ClientSecretExpiresAt, 0)
//                 issuedAt = time.Unix(res.ClientIDIssuedAt, 0)
//                 duration = time.Hour * 24 * 30
//         )
//         assert.True(expireAt.Sub(issuedAt) == duration, "return wrong issued date")
//
//         // Check the database to see if the client has been stored successfully
//         clientdb, exist := db.Client.Get(req.ClientName)
//         assert.Equal(true, exist, "client does not exist in the storage")
//         assert.Equal(res, *clientdb.ClientPrivate, "should point to the same object")
// }
//
// func TestClientRegistrationError(t *testing.T) {
//
//         assert := assert.New(t)
//
//         e := defaultMockEndpoint()
//
//         // Setup payload
//         req := &openid.ClientPublic{
//                 ClientName:   "openid_app",
//                 RedirectURIs: []string{"not_valid_url"},
//         }
//
//         rr := testClientRegistration(e, req, "")
//
//         var (
//                 statusCode   = 400
//                 contentType  = "application/json"
//                 cacheControl = "no-store"
//                 pragma       = "no-cache"
//
//                 header = rr.Header()
//         )
//
//         assert.Equal(statusCode, rr.Code, "should return status 400 - Bad Request")
//         assert.Equal(contentType, header.Get("Content-Type"), "should return json")
//         assert.Equal(cacheControl, header.Get("Cache-Control"), "should set cache-control to no-store")
//         assert.Equal(pragma, header.Get("Pragma"), "should set pragma to no-cache")
//
//         var (
//                 msg  = "invalid_redirect_uri"
//                 desc = "One or more redirect_uri values are incorrect"
//         )
//
//         var res openid.ClientErrorResponse
//         err := json.NewDecoder(rr.Body).Decode(&res)
//         assert.Nil(err)
//
//         assert.Equal(msg, res.Error, "should return the matching error")
//         assert.Equal(desc, res.ErrorDescription, "should return the matching error description")
// }
//
// func testClientRead(e *Endpoints, id string, bearer string) *httptest.ResponseRecorder {
//         router := httprouter.New()
//         router.GET("/connect/register", e.Client)
//
//         req := httptest.NewRequest("GET", "/connect/register?client_id="+id, nil)
//         req.Header.Set("Authorization", "Bearer "+bearer)
//         rr := httptest.NewRecorder()
//
//         router.ServeHTTP(rr, req)
//
//         return rr
// }
//
// func TestClientRead(t *testing.T) {
//
//         assert := assert.New(t)
//
//         e := defaultMockEndpoint()
//
//         var (
//                 cid = "s6BhdRkqt3"
//
//                 aud = "audience"
//                 sub = "subject"
//                 iss = "go-openid"
//                 dur = time.Minute
//                 rc  = crypto.New("secret")
//         )
//
//         bearer, err := rc.NewJWT(aud, sub, iss, dur)
//         assert.Nil(err)
//         rr := testClientRead(e, cid, bearer)
//
//         // Response Headers
//         var (
//                 statusCode   = 200
//                 contentType  = "application/json"
//                 cacheControl = "no-store"
//                 pragma       = "no-cache"
//
//                 header = rr.Header()
//         )
//
//         assert.Equal(statusCode, rr.Code, "should return status 200 - OK")
//         assert.Equal(contentType, header.Get("Content-Type"), "should return Content-Type application/json")
//         assert.Equal(cacheControl, header.Get("Cache-Control"), "should return Cache-Control no-store")
//         assert.Equal(pragma, header.Get("Pragma"), "should return Pragma no-cache")
//
//         var (
//                 clientID     = "s6BhdRkqt3"
//                 clientSecret = "OylyaC56ijpAQ7G5ZZGL7MMQ6Ap6mEeuhSTFVps2N4Q"
//         )
//
//         var client openid.Client
//         err = json.NewDecoder(rr.Body).Decode(&client)
//         assert.Nil(err)
//
//         assert.Equal(clientID, client.ClientPrivate.ClientID)
//         assert.Equal(clientSecret, client.ClientPrivate.ClientSecret)
//         // Response body
//         //	client := &openid.Client{
//         //		ClientPublic: &openid.ClientPublic{
//         //			TokenEndpointAuthMethod: "token_endpoint_auth_method",
//         //			ApplicationType:         "web",
//         //			RedirectURIs: []string{"https://client.example.org/callback",
//         //				"https://client.example.org/callback2"},
//         //			ClientName:          "My Example",
//         //			LogoURI:             "https://client.example.org/logo.png",
//         //			SubjectType:         "pairwise",
//         //			SectorIdentifierURI: "https://other.example.net/file_of_redirect_uris.json",
//         //			JwksURI:             "https://client.example.org/my_public_keys.jwks",
//         //			UserinfoEncryptedResponseAlg: "RSA1_5",
//         //			UserinfoEncryptedResponseEnc: "A128CBC-HS256",
//         //			Contacts:                     []string{"ve7jtb@example.org", "mary@example.org"},
//         //			RequestURIs:                  []string{"https://client.example.org/rf.txt#qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"},
//         //		},
//         //		ClientPrivate: &openid.ClientPrivate{
//         //			ClientID:     "s6BhdRkqt3",
//         //			ClientSecret: "OylyaC56ijpAQ7G5ZZGL7MMQ6Ap6mEeuhSTFVps2N4Q",
//         //
//         //			RegistrationAccessToken: "",
//         //			RegistrationClientURI:   "https://server.example.com/connect/register?client_id=s6BhdRkqt3",
//         //			ClientIDIssuedAt:        0,
//         //			ClientSecretExpiresAt:   17514165600,
//         //		},
//         //	}
// }
//
// func TestClientReadError(t *testing.T) {
//         assert := assert.New(t)
//
//         e := defaultMockEndpoint()
//
//         cid := "unknown_client_id"
//
//         rr := testClientRead(e, cid, "")
//
//         var (
//                 statusCode   = 401
//                 cacheControl = "no-store"
//                 pragma       = "no-cache"
//
//                 header = rr.Header()
//         )
//
//         assert.Equal(statusCode, rr.Code, "should return status 401 - Unauthorized")
//         assert.Equal(cacheControl, header.Get("Cache-Control"), "should set cache-control to no-store")
//         assert.Equal(pragma, header.Get("Pragma"), "should set pragma to no-cache")
// }
//
// func TestOpenIDConfigurationRequest(t *testing.T) {
//         url := "/.well-known/openid-configuration"
//         assert := assert.New(t)
//         assert.True(url == url)
//         // 	{
//         //    "issuer":
//         //      "https://server.example.com",
//         //    "authorization_endpoint":
//         //      "https://server.example.com/connect/authorize",
//         //    "token_endpoint":
//         //      "https://server.example.com/connect/token",
//         //    "token_endpoint_auth_methods_supported":
//         //      ["client_secret_basic", "private_key_jwt"],
//         //    "token_endpoint_auth_signing_alg_values_supported":
//         //      ["RS256", "ES256"],
//         //    "userinfo_endpoint":
//         //      "https://server.example.com/connect/userinfo",
//         //    "check_session_iframe":
//         //      "https://server.example.com/connect/check_session",
//         //    "end_session_endpoint":
//         //      "https://server.example.com/connect/end_session",
//         //    "jwks_uri":
//         //      "https://server.example.com/jwks.json",
//         //    "registration_endpoint":
//         //      "https://server.example.com/connect/register",
//         //    "scopes_supported":
//         //      ["openid", "profile", "email", "address",
//         //       "phone", "offline_access"],
//         //    "response_types_supported":
//         //      ["code", "code id_token", "id_token", "token id_token"],
//         //    "acr_values_supported":
//         //      ["urn:mace:incommon:iap:silver",
//         //       "urn:mace:incommon:iap:bronze"],
//         //    "subject_types_supported":
//         //      ["public", "pairwise"],
//         //    "userinfo_signing_alg_values_supported":
//         //      ["RS256", "ES256", "HS256"],
//         //    "userinfo_encryption_alg_values_supported":
//         //      ["RSA1_5", "A128KW"],
//         //    "userinfo_encryption_enc_values_supported":
//         //      ["A128CBC-HS256", "A128GCM"],
//         //    "id_token_signing_alg_values_supported":
//         //      ["RS256", "ES256", "HS256"],
//         //    "id_token_encryption_alg_values_supported":
//         //      ["RSA1_5", "A128KW"],
//         //    "id_token_encryption_enc_values_supported":
//         //      ["A128CBC-HS256", "A128GCM"],
//         //    "request_object_signing_alg_values_supported":
//         //      ["none", "RS256", "ES256"],
//         //    "display_values_supported":
//         //      ["page", "popup"],
//         //    "claim_types_supported":
//         //      ["normal", "distributed"],
//         //    "claims_supported":
//         //      ["sub", "iss", "auth_time", "acr",
//         //       "name", "given_name", "family_name", "nickname",
//         //       "profile", "picture", "website",
//         //       "email", "email_verified", "locale", "zoneinfo",
//         //       "http://example.info/claims/groups"],
//         //    "claims_parameter_supported":
//         //      true,
//         //    "service_documentation":
//         //      "http://server.example.com/connect/service_documentation.html",
//         //    "ui_locales_supported":
//         //      ["en-US", "en-GB", "en-CA", "fr-FR", "fr-CA"]
//         //   }
// }
//
// type cry struct {
//         *crypto.Impl
// }
//
// func (c *cry) Code() string {
//         return defaultXIDToken
// }
//
// func (c *cry) NewJWT(aud, sub, iss string, dur time.Duration) (string, error) {
//         return defaultJWTToken, nil
// }
//
// func (c *cry) UUID() string {
//         return defaultUUID
// }
//
// func newMockCrypto() crypto.Crypto {
//         c := crypto.New("secret")
//         return &cry{c}
// }
//
// func newMockService(db *Database) *ServiceImpl {
//         c := newMockCrypto()
//         return NewService(db, c)
// }
//
// func newMockEndpoint(s Service) *Endpoints {
//         return &Endpoints{
//                 service: s,
//         }
// }
//
// func newMockDatabase() *Database {
//         claims := &openid.StandardClaims{
//                 Profile: &openid.Profile{
//                         Sub:               "248289761001",
//                         Name:              "Jane Doe",
//                         GivenName:         "Jane",
//                         FamilyName:        "Doe",
//                         PreferredUsername: "j.doe",
//                         Picture:           "http://example.com/janedoe/me.jpg",
//                 },
//                 Email: &openid.Email{
//                         Email: "janedoe@example.com",
//                 },
//         }
//
//         user := &User{
//                 ID:             "1",
//                 StandardClaims: claims,
//         }
//
//         client := &openid.Client{
//                 ClientPublic: &openid.ClientPublic{
//                         ClientName:   defaultClientName,
//                         RedirectURIs: []string{defaultRedirectURI},
//                 },
//                 ClientPrivate: &openid.ClientPrivate{
//                         ClientID:     defaultClientID,
//                         ClientSecret: defaultClientSecret,
//                 },
//         }
//
//         client2 := &openid.Client{
//                 ClientPublic: &openid.ClientPublic{
//                         TokenEndpointAuthMethod: "token_endpoint_auth_method",
//                         ApplicationType:         "web",
//                         RedirectURIs: []string{"https://client.example.org/callback",
//                                 "https://client.example.org/callback2"},
//                         ClientName:                   "My Example",
//                         LogoURI:                      "https://client.example.org/logo.png",
//                         SubjectType:                  "pairwise",
//                         SectorIdentifierURI:          "https://other.example.net/file_of_redirect_uris.json",
//                         JwksURI:                      "https://client.example.org/my_public_keys.jwks",
//                         UserinfoEncryptedResponseAlg: "RSA1_5",
//                         UserinfoEncryptedResponseEnc: "A128CBC-HS256",
//                         Contacts:                     []string{"ve7jtb@example.org", "mary@example.org"},
//                         RequestURIs:                  []string{"https://client.example.org/rf.txt#qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"},
//                 },
//                 ClientPrivate: &openid.ClientPrivate{
//                         ClientID:     "s6BhdRkqt3",
//                         ClientSecret: "OylyaC56ijpAQ7G5ZZGL7MMQ6Ap6mEeuhSTFVps2N4Q",
//
//                         RegistrationAccessToken: "",
//                         RegistrationClientURI:   "https://server.example.com/connect/register?client_id=s6BhdRkqt3",
//                         ClientIDIssuedAt:        0,
//                         ClientSecretExpiresAt:   17514165600,
//                 },
//         }
//
//         db := NewDatabase()
//
//         db.Client.Put(client.ClientID, client)
//         db.Code.Put(client.ClientID, openid.NewCode(defaultCode))
//         db.Client.Put(client2.ClientPrivate.ClientID, client2)
//
//         db.User.Put("1", user)
//         return db
// }
//
// func defaultMockEndpoint() *Endpoints {
//         db := newMockDatabase()
//         s := newMockService(db)
//         e := newMockEndpoint(s)
//         return e
// }
