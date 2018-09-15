package openidsvc

func getWebfinger() {
	// curl -H "Host: example.com"
	// http://localhost:8080/.well-known/webfinger
	//	?resource=acct:joe@example.com
	//  &rel=http://openid.net/specs/connect/1.0/issuer
	//   HTTP/1.1 200 OK
	//   Content-Type: application/jrd+json

	//   {
	//    "subject": "acct:joe@example.com",
	//    "links":
	//     [
	//      {
	//       "rel": "http://openid.net/specs/connect/1.0/issuer",
	//       "href": "https://server.example.com"
	//      }
	//     ]
	//   }
}
