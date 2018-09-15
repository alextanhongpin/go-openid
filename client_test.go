package oidc_test

import (
	"testing"

	"github.com/alextanhongpin/go-openid"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshallClientJSON(t *testing.T) {
	assert := assert.New(t)

	body := []byte(`{
		"application_type": "web",
		"redirect_uris": [
			"https://client.example.org/callback",
			"https://client.example.org/callback2"
		],
		"client_name": "My Example",
		"client_name#ja-Jpan-JP": "クライアント名",
		"logo_uri": "https://client.example.org/logo.png",
		"subject_type": "pairwise",
		"sector_identifier_uri": "https://other.example.net/file_of_redirect_uris.json",
		"token_endpoint_auth_method": "client_secret_basic",
		"jwks_uri": "https://client.example.org/my_public_keys.jwks",
		"userinfo_encrypted_response_alg": "RSA1_5",
		"userinfo_encrypted_response_enc": "A128CBC-HS256",
		"contacts": [
			"ve7jtb@example.org",
			"mary@example.org"
		],
		"request_uri": [
			"https://client.example.org/rf.txt#qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"
		]
	}`)

	c := oidc.ClientRegistrationRequest{}

	err := c.UnmarshalJSON(body)
	assert.Nil(err)

	err = c.Validate()
	assert.Nil(err)

	assert.Equal("web", c.ApplicationType)
}

// func makeAccessTokenRequest () {
// 	t := &http.Transport{
// 		Dial: (&net.Dialer{
// 			Timeout: 5 *time.Second,
// 			KeepAlive: 5 *time.Second,
// 		}).Dial,
// 		TLSHandshakeTimeout: 5 * time.Second,
// 		ResponseHeaderTimeout: 5 * time.Second,
// 		ExpectContinueTimeout: 1  * time.Second,
// 	}
// 	client := &http.Client {
// 		Timeout: 10 * time.Second,
// 		Transport: t,
// 	}
// 	// res, err := client.Get(url)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	req, err := http.NewRequest("GET", "url", nil)
// 	if err != nil {
// 	log.Fatal(err)
// 	}
// 	req = req.WithContext(ctx)
// 	res, err := client.Do(req)
// 	if err != nil {
//
// 	}
// 	defer res.Body.Close()
//
// }
//
// func HandleAccessTokenRequest (w http.ResponseWriter, r *http.Request) {
//
// 	r.Header().Get("Authorization")
// 	r.Header().Get("Content-Type") == "application/x-www-form-urlencoded"
//
// 	var req AccessTokenRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 	}
// 	if err := req.Validate(); err != nil {
// 	}
// 	if err := FindCodeInCache(req.Code, req.ClientID) {
//
// 	}
// 	cdb, err := FindClient(req.ClientID)
// 	if err != nil {
//
// 	}
// 	if cdb.RedirectURI != req.RedirectURI {}
//
// 	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
// 	w.Header().Set("Cache-Control", "no-store")
// 	w.Header().Set("Pragma", "no-cache")
// }
