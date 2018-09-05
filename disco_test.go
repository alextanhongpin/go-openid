package oidc

func Test(t * testing.T) {
	t.Error("hello")
}
//
// func TestOIDProviderIssuerDiscoveryEmail(t *testing.T) {
//         u, _ := url.Parse("http://example.com/.well-known/webfinger?resource=acct%3Ajoe%40example.com&rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer")
//         q := u.Query()
//         resource := "acct:joe@example.com"
//         host := "example.com"
//         rel := "http://openid.net/specs/connect/1.0/issuer"
//
//         assert := assert.New(t)
//         assert.Equal(resource, q.Get("resource"), "should have resource in qs")
//         assert.Equal(rel, q.Get("rel"), "should have rel in qs")
//
//         // Header response
//         statusCode := 200
//         contentType := "application/jrd+json"
//         assert.Equal(host, host, "should have host in qs")
//         assert.Equal(http.StatusOK, statusCode, "should return status 200 - Ok")
//         assert.Equal(contentType, contentType, "should return the correct content type")
//         //	{
//         //   "subject": "acct:joe@example.com",
//         //   "links":
//         //    [
//         //     {
//         //      "rel": "http://openid.net/specs/connect/1.0/issuer",
//         //      "href": "https://server.example.com"
//         //     }
//         //    ]
//         //  }
// }
//
// func TestOIDProviderIssuerDiscoveryURL(t *testing.T) {
//         u, _ := url.Parse("http://example.com/.well-known/webfinger?resource=https%3A%2F%2Fexample.com%2Fjoe&rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer")
//         q := u.Query()
//         resource := "https://example.com/joe"
//         rel := "http://openid.net/specs/connect/1.0/issuer"
//
//         assert := assert.New(t)
//         assert.Equal(resource, q.Get("resource"), "should have the correct resource in the qs")
//         assert.Equal(rel, q.Get("rel"), "should have the correct rel in the qs")
//
//         host := "example.com"
//         statusCode := 200
//         contentType := "application/jrd+json"
//         assert.Equal(http.StatusOK, statusCode, "should return status 200 - Ok")
//         assert.Equal(contentType, contentType, "should return the correct content type")
//         assert.Equal(host, host, "should return the correct host")
//         // Test response
//         //	{
//         //   "subject": "https://example.com/joe",
//         //   "links":
//         //    [
//         //     {
//         //      "rel": "http://openid.net/specs/connect/1.0/issuer",
//         //      "href": "https://server.example.com"
//         //     }
//         //    ]
//         //  }
// }
//
// func TestOIDPProviderUserDiscoveryHostnameAndPort(t *testing.T) {
//         // TODO: User input using hostname and port syntax
//         resource := "https://example.com:8080/"
//         host := "example.com:8080"
//         rel := "http://openid.net/specs/connect/1.0/issuer"
//
//         assert := assert.New(t)
//         assert.True(resource == resource)
//         assert.True(host == host)
//         assert.True(rel == rel)
//
//         //	  GET /.well-known/webfinger
//         //    ?resource=https%3A%2F%2Fexample.com%3A8080%2F
//         //    &rel=http%3A%2F%2Fopenid.net%2Fspecs%2Fconnect%2F1.0%2Fissuer
//         //    HTTP/1.1
//         //  Host: example.com:8080
//         //
//         //  HTTP/1.1 200 OK
//         //  Content-Type: application/jrd+json
//         //
//         //  {
//         //   "subject": "https://example.com:8080/",
//         //   "links":
//         //    [
//         //     {
//         //      "rel": "http://openid.net/specs/connect/1.0/issuer",
//         //      "href": "https://server.example.com"
//         //     }
//         //    ]
//         //  }
// }
//

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
