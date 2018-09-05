package oidc_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/querystring"

	"github.com/stretchr/testify/assert"
)

func TestDiscoveryEmail(t *testing.T) {
	assert := assert.New(t)

	u, _ := url.Parse("http://example.com/.well-known/webfinger?resource=acct%3Ajoe%40example.com&rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer")
	q := u.Query()

	var (
		resource = "acct:joe@example.com"
		rel      = "http://openid.net/specs/connect/1.0/issuer"
		host     = "example.com"
	)

	var res oidc.Discovery
	err := querystring.Decode(&res, q)
	assert.Nil(err)

	assert.Equal(host, u.Host, "should have the matching host")
	assert.Equal(resource, res.Resource, "should have resource in qs")
	assert.Equal(rel, res.Rel, "should have rel in qs")

	// Header response
	var (
		statusCode  = 200
		contentType = "application/jrd+json"
	)

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

func TestDiscoveryURL(t *testing.T) {
	assert := assert.New(t)

	u, _ := url.Parse("http://example.com/.well-known/webfinger?resource=https%3A%2F%2Fexample.com%2Fjoe&rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer")
	q := u.Query()

	var (
		resource = "https://example.com/joe"
		rel      = "http://openid.net/specs/connect/1.0/issuer"
	)

	var res oidc.Discovery
	err := querystring.Decode(&res, q)
	assert.Nil(err)

	assert.Equal(resource, res.Resource, "should have the correct resource in the qs")
	assert.Equal(rel, res.Rel, "should have the correct rel in the qs")

	var (
		host        = "example.com"
		statusCode  = 200
		contentType = "application/jrd+json"
	)

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

func TestHostnameAndPort(t *testing.T) {
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

// func TestDiscoverUserInputAcct(t *testing.T) {
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
// }
