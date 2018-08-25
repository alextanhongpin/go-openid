package openid

import (
	"log"
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

func TestAuthorizationRequest(t *testing.T) {
	payload := map[string]string{
		"response_type": "code",
		"client_id":     "abc123",
		"redirect_uri":  "http://client.com",
		"scope":         "profile",
		"state":         "xyz",
	}
	req, err := http.NewRequest("GET", "/authorize", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	for k, v := range payload {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleAuthorizationRequest)
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
	qq := u.Query()

	if state := qq.Get("state"); state != payload["state"] {
		t.Errorf("handler returned wrong state: got %v want %v", state, payload["state"])
	}

	// Find a way to mock the code generation
	if code := qq.Get("code"); code != payload["code"] {
		t.Errorf("handler returned wrong code: got %v want %v", code, payload["code"])
	}
	log.Printf("got %#v", rr.Header().Get("Location"))
}

// Testing POST
// js, _ := json.Marshal(struct)
// req, err := http.NewRequest("POST", "/endpoint", bytes.NewBuffer(js))
