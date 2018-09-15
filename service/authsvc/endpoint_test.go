package authsvc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

var mockDB map[string]User
var mockauthsvc = MakeServerEndpoints(mocksvc{})

func createDB() map[string]User {
	if mockDB == nil {
		mockDB = make(map[string]User)
		mockDB["123456"] = User{
			Email:    &Email{Email: "john.doe@mail.com"},
			Password: "123456",
		}
	}
	return mockDB
}

func TestPostRegisterExistingUserEndpoint(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"email":"john.doe@mail.com", "password": "123456"}`)))

	router := httprouter.New()
	router.Handle("POST", "/register", mockauthsvc.PostRegister())
	router.ServeHTTP(w, r)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)

	log.Print(string(body))

	expectedBody := `{"ok":false,"error":"User with the email already exists"}`
	gotBody := strings.TrimSpace(string(body))

	if expectedBody != gotBody {
		t.Errorf("got %v, want %v", expectedBody, gotBody)
	}

	expectedStatusCode := 200
	gotStatusCode := res.StatusCode
	if expectedStatusCode != gotStatusCode {
		t.Errorf("got %v, want %v", gotStatusCode, expectedStatusCode)
	}
}
func TestPostRegisterNewUserEndpoint(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"email":"jane.doe@mail.com", "password": "123456"}`)))

	router := httprouter.New()
	router.Handle("POST", "/register", mockauthsvc.PostRegister())
	router.ServeHTTP(w, r)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)

	expectedBody := make(map[string]interface{})
	expectedBody["ok"] = "true"
	expectedBody["user_id"] = "123456"
	expectedBody["redirect_uri"] = "/users/123456"

	gotBody := make(map[string]interface{})
	_ = json.Unmarshal(body, &gotBody)
	delete(gotBody, "access_token")

	if reflect.DeepEqual(expectedBody, gotBody) {
		t.Errorf("got %v, want %v", expectedBody, gotBody)
	}

	expectedStatusCode := 200
	gotStatusCode := res.StatusCode
	if expectedStatusCode != gotStatusCode {
		t.Errorf("got %d, want %d", gotStatusCode, expectedStatusCode)
	}
}

func TestPostLoginExistingUserEndpoint(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"email":"john.doe@mail.com", "password": "123456"}`)))

	router := httprouter.New()
	router.Handle("POST", "/login", mockauthsvc.PostLogin())
	router.ServeHTTP(w, r)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)

	log.Print(string(body))

	expectedBody := `{"ok":true,"redirect_uri":"/login/callback?user_id="}`
	gotBody := strings.TrimSpace(string(body))

	if expectedBody != gotBody {
		t.Errorf("got %v, want %v", gotBody, expectedBody)
	}

	expectedStatusCode := 200
	gotStatusCode := res.StatusCode
	if expectedStatusCode != gotStatusCode {
		t.Errorf("got %v, want %v", gotStatusCode, expectedStatusCode)
	}
}

func TestPostLoginNewUserEndpoint(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"email":"jane.doe@mail.com", "password": "123456"}`)))

	router := httprouter.New()
	router.Handle("POST", "/login", mockauthsvc.PostLogin())
	router.ServeHTTP(w, r)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)

	expectedBody := `{"error":400,"message":"No user with the email found"}`
	gotBody := strings.TrimSpace(string(body))

	if expectedBody != gotBody {
		t.Errorf("got %v, want %v", gotBody, expectedBody)
	}

	expectedStatusCode := 400
	gotStatusCode := res.StatusCode
	if expectedStatusCode != gotStatusCode {
		t.Errorf("got %v, want %v", gotStatusCode, expectedStatusCode)
	}
}

// Setup mocksvc for mocking the authsvc
type mocksvc struct{}

func (s mocksvc) GetUser(getUserRequest) (*getUserResponse, error) {
	return &getUserResponse{
		Data: User{},
	}, nil
}

func (s mocksvc) GetUsers(getUsersRequest) (*getUsersResponse, error) {
	return &getUsersResponse{
		Data: []User{},
	}, nil
}

func (s mocksvc) DeleteUser(deleteUserRequest) (*deleteUserResponse, error) {
	return &deleteUserResponse{
		Ok: true,
	}, nil
}

func (s mocksvc) CreateUser(createUserRequest) (*createUserResponse, error) {
	return &createUserResponse{
		ID: "123456",
	}, nil
}

func (s mocksvc) CheckUser(email string) (*User, error) {
	db := createDB()
	for _, user := range db {
		if user.Email.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (s mocksvc) UpdateUser(req updateUserRequest) (*updateUserResponse, error) {
	return &updateUserResponse{}, nil
}
