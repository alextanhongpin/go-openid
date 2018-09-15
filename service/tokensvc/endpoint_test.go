package tokensvc

// import "testing"

// func TestTokenRequest(t *testing.T) {
// 	expectedGrantType := "authorization_code"
// 	expectedBearerType := "Basic"
// 	expectedContentType := "application/x-www-form-urlencoded"
// 	// 	 POST /token HTTP/1.1
// 	//   Host: server.example.com
// 	//   Content-Type: application/x-www-form-urlencoded
// 	//   Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW

// 	//   grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA
// 	//     &redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb
// 	expectedCacheControl := "no-store"
// 	expectedPragma := "no-cache"

// 	expectedIDToken := "something"
// 	//   HTTP/1.1 200 OK
// 	//   Content-Type: application/json
// 	//   Cache-Control: no-store
// 	//   Pragma: no-cache

// 	//   {
// 	//    "access_token": "SlAV32hkKG",
// 	//    "token_type": "Bearer",
// 	//    "refresh_token": "8xLOxBtZp8",
// 	//    "expires_in": 3600,
// 	//    "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ.ewogImlzc
// 	//      yI6ICJodHRwOi8vc2VydmVyLmV4YW1wbGUuY29tIiwKICJzdWIiOiAiMjQ4Mjg5
// 	//      NzYxMDAxIiwKICJhdWQiOiAiczZCaGRSa3F0MyIsCiAibm9uY2UiOiAibi0wUzZ
// 	//      fV3pBMk1qIiwKICJleHAiOiAxMzExMjgxOTcwLAogImlhdCI6IDEzMTEyODA5Nz
// 	//      AKfQ.ggW8hZ1EuVLuxNuuIJKX_V8a_OMXzR0EHR9R6jgdqrOOF4daGU96Sr_P6q
// 	//      Jp6IcmD3HP99Obi1PRs-cwh3LO-p146waJ8IhehcwL7F09JdijmBqkvPeB2T9CJ
// 	//      NqeGpe-gccMg4vfKjkM8FcGvnzZUN4_KSP0aAp1tOJ1zZwgjxqGByKHiOtX7Tpd
// 	//      QyHE5lcMiKPXfEIQILVq0pc_E2DzL7emopWoaoZTF_m0_N0YzFC6g6EJbOEoRoS
// 	//      K5hoDalrcvRYLSrQAZZKflyuVCyixEoV9GfNQC3_osjzw2PAithfubEEBLuVVk4
// 	//      XUVrWOLrLl0nx7RkKU8NXNHq-rvKMzqg"
// 	//   }
// }

// func TestTokenError(t *testing.T) {
// 	// 	  HTTP/1.1 400 Bad Request
// 	//   Content-Type: application/json
// 	//   Cache-Control: no-store
// 	//   Pragma: no-cache

// 	//   {
// 	//    "error": "invalid_request"
// 	//   }
// }

// func TestIDToken(t *testing.T) {

// }

// func TestUserInfoEndpoint(t *testing.T) {
// 	// Request
// 	//   GET /userinfo HTTP/1.1
// 	//   Host: server.example.com
// 	//   Authorization: Bearer SlAV32hkKG

// 	// Success Response
// 	// 	  HTTP/1.1 200 OK
// 	//   Content-Type: application/json

// 	//   {
// 	//    "sub": "248289761001",
// 	//    "name": "Jane Doe",
// 	//    "given_name": "Jane",
// 	//    "family_name": "Doe",
// 	//    "preferred_username": "j.doe",
// 	//    "email": "janedoe@example.com",
// 	//    "picture": "http://example.com/janedoe/me.jpg"
// 	//   }

// 	// Error Response
// 	// 	 HTTP/1.1 401 Unauthorized
// 	//   WWW-Authenticate: error="invalid_token",
// 	//     error_description="The Access Token expired"
// }
