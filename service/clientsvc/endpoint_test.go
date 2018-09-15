package clientsvc

// import (
// 	"io/ioutil"
// 	"net/http/httptest"
// 	"testing"

// 	"strings"

// 	"bytes"

// 	"github.com/alextanhongpin/go-openid/models"
// 	"github.com/julienschmidt/httprouter"
// )

// type mocksvc struct{}

// func (s mocksvc) GetClientMetadata(request getClientMetadataRequest) (*getClientMetadataResponse, error) {
// 	return &getClientMetadataResponse{
// 		Data: models.ClientMetadata{
// 			ApplicationType: "web",
// 		},
// 	}, nil
// }
// func (s mocksvc) GetClientsMetadata(getClientsMetadataRequest) (*getClientsMetadataResponse, error) {
// 	return &getClientsMetadataResponse{
// 		Data: []models.ClientMetadata{
// 			models.ClientMetadata{
// 				ApplicationType: "web",
// 			},
// 		},
// 	}, nil
// }
// func (s mocksvc) PostClientMetadata(postClientMetadataRequest) (*postClientMetadataResponse, error) {
// 	return &postClientMetadataResponse{
// 		Data: models.ClientMetadata{

// 		// "client_id": "s6BhdRkqt3",
// 		// "client_secret": "ZJYCqe3GGRvdrudKyZS0XhGv_Z45DuKhCUk0gBR1vZk",
// 		// "client_secret_expires_at": 1577858400,
// 		// "registration_access_token": "this.is.an.access.token.value.ffx83",
// 		// "registration_client_uri": "https://server.example.com/connect/register?client_id=s6BhdRkqt3",
// 		// "token_endpoint_auth_method": "client_secret_basic",
// 		// "application_type": "web",
// 		// "redirect_uris": ["https://client.example.org/callback","https://client.example.org/callback2"],
// 		// "client_name": "My Example",
// 		// "client_name#ja-Jpan-JP": "クライアント名",
// 		// "logo_uri": "https://client.example.org/logo.png",
// 		// "subject_type": "pairwise",
// 		// "sector_identifier_uri": "https://other.example.net/file_of_redirect_uris.json",
// 		// "jwks_uri": "https://client.example.org/my_public_keys.jwks",
// 		// "userinfo_encrypted_response_alg": "RSA1_5",
// 		// "userinfo_encrypted_response_enc": "A128CBC-HS256",
// 		// "contacts": ["ve7jtb@example.org", "mary@example.org"],
// 		// "request_uris": ["https://client.example.org/rf.txt #qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"]

// 		},
// 	}, nil
// }

// var s = mocksvc{}
// var e = MakeServerEndpoints(s)

// func TestGetClientsMetadata(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/api/clients", nil)

// 	router := httprouter.New()
// 	router.Handle("GET", "/api/clients", e.GetClients())
// 	router.ServeHTTP(w, r)
// 	res := w.Result()
// 	body, _ := ioutil.ReadAll(res.Body)

// 	expectedBody := `{"data":[{"id":"","application_type":"web"}],"count":0}`
// 	actualBody := strings.TrimSpace(string(body))
// 	if actualBody != expectedBody {
// 		t.Errorf("got %v, want %v", actualBody, expectedBody)
// 	}

// 	expectedStatusCode := 200
// 	actualStatusCode := res.StatusCode
// 	if actualStatusCode != expectedStatusCode {
// 		t.Errorf("got %v, want %v", actualStatusCode, expectedStatusCode)
// 	}

// 	expectedContentType := "application/json"
// 	actualContentType := res.Header.Get("Content-Type")
// 	if actualContentType != expectedContentType {
// 		t.Errorf("got %v, want %v", actualContentType, expectedContentType)
// 	}

// 	expectedCacheControl := "no-store"
// 	actualCacheControl := res.Header.Get("Cache-Control")
// 	if actualCacheControl != expectedCacheControl {
// 		t.Errorf("got %v, want %v", actualCacheControl, expectedCacheControl)
// 	}

// 	expectedPragma := "no-cache"
// 	actualPragma := res.Header.Get("Pragma")
// 	if actualPragma != expectedPragma {
// 		t.Errorf("got %v, want %v", actualPragma, expectedPragma)
// 	}
// }

// func TestGetClientMetadata(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/api/clients", nil)

// 	router := httprouter.New()
// 	router.Handle("GET", "/api/clients", e.GetClient())
// 	router.ServeHTTP(w, r)
// 	res := w.Result()
// 	body, _ := ioutil.ReadAll(res.Body)

// 	expectedBody := `{"data":{"id":"","application_type":"web"}}`
// 	actualBody := strings.TrimSpace(string(body))
// 	if actualBody != expectedBody {
// 		t.Errorf("got %v, want %v", actualBody, expectedBody)
// 	}

// 	expectedStatusCode := 200
// 	actualStatusCode := res.StatusCode
// 	if actualStatusCode != expectedStatusCode {
// 		t.Errorf("got %v, want %v", actualStatusCode, expectedStatusCode)
// 	}
// }

// func TestPostClientMetadata(t *testing.T) {

// 	//   POST /connect/register HTTP/1.1
// 	//   Content-Type: application/json
// 	//   Accept: application/json
// 	//   Host: server.example.com
// 	//   Authorization: Bearer eyJhbGciOiJSUzI1NiJ9.eyJ ...
// 	payload := []byte(`{
// 		"application_type": "web",
// 		"redirect_uris": ["https://client.example.org/callback", "https://client.example.org/callback2"],
// 		"client_name": "My Example",
// 		"client_name#ja-Jpan-JP": "クライアント名",
// 		"logo_uri": "https://client.example.org/logo.png",
// 		"subject_type": "pairwise",
// 		"sector_identifier_uri": "https://other.example.net/file_of_redirect_uris.json",
// 		"token_endpoint_auth_method": "client_secret_basic",
// 		"jwks_uri": "https://client.example.org/my_public_keys.jwks",
// 		"userinfo_encrypted_response_alg": "RSA1_5",
// 		"userinfo_encrypted_response_enc": "A128CBC-HS256",
// 		"contacts": ["ve7jtb@example.org", "mary@example.org"],
// 		"request_uris": ["https://client.example.org/rf.txt#qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"]
// 	}`)

// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("POST", "/connect/register", bytes.NewBuffer(payload))
// 	r.Header.Set("Content-Type", "application/json")
// 	r.Header.Set("Accept", "application/json")

// 	router := httprouter.New()
// 	router.Handle("POST", "/connect/register", e.PostClient())
// 	router.ServeHTTP(w, r)
// 	res := w.Result()
// 	body, _ := ioutil.ReadAll(res.Body)

// 	expectedBody := strings.TrimSpace(`{
// 		"client_id": "s6BhdRkqt3",
// 		"client_secret": "ZJYCqe3GGRvdrudKyZS0XhGv_Z45DuKhCUk0gBR1vZk",
// 		"client_secret_expires_at": 1577858400,
// 		"registration_access_token": "this.is.an.access.token.value.ffx83",
// 		"registration_client_uri": "https://server.example.com/connect/register?client_id=s6BhdRkqt3",
// 		"token_endpoint_auth_method": "client_secret_basic",
// 		"application_type": "web",
// 		"redirect_uris": ["https://client.example.org/callback","https://client.example.org/callback2"],
// 		"client_name": "My Example",
// 		"client_name#ja-Jpan-JP": "クライアント名",
// 		"logo_uri": "https://client.example.org/logo.png",
// 		"subject_type": "pairwise",
// 		"sector_identifier_uri": "https://other.example.net/file_of_redirect_uris.json",
// 		"jwks_uri": "https://client.example.org/my_public_keys.jwks",
// 		"userinfo_encrypted_response_alg": "RSA1_5",
// 		"userinfo_encrypted_response_enc": "A128CBC-HS256",
// 		"contacts": ["ve7jtb@example.org", "mary@example.org"],
// 		"request_uris": ["https://client.example.org/rf.txt #qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"]
// 	}`)
// 	actualBody := strings.TrimSpace(string(body))
// 	if actualBody != expectedBody {
// 		t.Errorf("got %v, want %v", actualBody, expectedBody)
// 	}

// 	expectedStatusCode := 201
// 	actualStatusCode := res.StatusCode
// 	if expectedStatusCode != actualStatusCode {
// 		t.Errorf("got %v, want %v", actualStatusCode, expectedStatusCode)
// 	}

// 	expectedContentType := "application/json"
// 	actualContentType := res.Header.Get("Content-Type")
// 	if expectedContentType != actualContentType {
// 		t.Errorf("got %v, want %v", actualContentType, expectedContentType)
// 	}

// 	expectedCacheControl := "no-store"
// 	actualCacheControl := res.Header.Get("cache-control")
// 	if expectedCacheControl != actualCacheControl {
// 		t.Errorf("got %v, want %v", actualCacheControl, expectedCacheControl)
// 	}

// 	expectedPragma := "no-cache"
// 	actualPragma := res.Header.Get("Pragma")
// 	if expectedPragma != actualPragma {
// 		t.Errorf("got %v, want %v", actualPragma, expectedPragma)
// 	}
// }
