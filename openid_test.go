package openid

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "hello world")
// }
//
// func TestAccessTokenRequest(t *testing.T) {
// 	req, err := http.NewRequest("POST", "/token", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(helloWorldHandler)
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, rr.Code)
// 	}
// 	expected := `{"alive": true}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned wrong body: got %v want %v", rr.Body.String(), expected)
// 	}
// }
type mockAuthService struct{}

const authorizationCode = "secret_code"

func (s *mockAuthService) GenerateCode() string {
	return authorizationCode
}

func TestAuthorizationRequest(t *testing.T) {
	payload := AuthorizationRequest{
		ResponseType: "code",
		ClientID:     "abc123",
		RedirectURI:  "http://client.com",
		Scope:        "profile",
		State:        "xyz",
	}
	q, err := EncodeAuthorizationRequest(&payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/authorize", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()

	h := HandleAuthorizationRequest(&mockAuthService{})
	handler := http.HandlerFunc(h)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Check redirect URL
	location := rr.Header().Get("Location")
	u, err := url.Parse(location)
	if err != nil {
		t.Fatal(err)
	}
	authRes := DecodeAuthorizationResponse(u.Query())

	if state := authRes.State; state != payload.State {
		t.Errorf("handler returned wrong state: got %v want %v", state, payload.State)
	}

	// Find a way to mock the code generation
	if code := authRes.Code; code != authorizationCode {
		t.Errorf("handler returned wrong code: got %v want %v", code, authorizationCode)
	}
}

// Testing POST
// js, _ := json.Marshal(struct)
// req, err := http.NewRequest("POST", "/endpoint", bytes.NewBuffer(js))
