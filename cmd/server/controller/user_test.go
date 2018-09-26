package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/cmd/server/controller"
	"github.com/alextanhongpin/go-openid/pkg/appsensor"
	"github.com/alextanhongpin/go-openid/pkg/session"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestUserRegister(t *testing.T) {
	assert := assert.New(t)

	sessMgr := session.NewManager()
	aps := appsensor.NewLoginDetector()

	// Setup Controller.
	userController := controller.NewUser()
	userController.SetAppSensor(aps)
	userController.SetSession(sessMgr)

	tests := []struct {
		test, email, password, desc string
	}{
		{"register with empty params", "", "", "email cannot be empty"},
		{"register with invalid email and no password", "x", "", "invalid email"},
		{"register with valid email and no password", "john.doe@mail.com", "", "arguments cannot be empty"},
		{"register with no email and with password", "", "1234", "email cannot be empty"},
		{"register with valid email and invalid password", "john.doe@mail.com", "1234", "password cannot be less than 8 characters"},
		{"register with invalid email and valid password", "john.doemail.com", "12345678", "invalid email"},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			creds := &controller.Credentials{
				Email:    tt.email,
				Password: tt.password,
			}
			rr := postRegisterEndpoint(&userController, creds)
			assert.Equal(400, rr.Code)

			var res oidc.ErrorJSON
			err := json.NewDecoder(rr.Body).Decode(&res)
			assert.Nil(err)
			assert.Equal(tt.desc, res.Code)
		})
	}

	t.Run("register with valid email and password", func(t *testing.T) {
		credentials := &controller.Credentials{
			Email:    "john.doe@mail.com",
			Password: "12345678",
		}
		rr := postRegisterEndpoint(&userController, credentials)
		assert.Equal(200, rr.Code)
		type response struct {
			AccessToken string `json:"access_token"`
		}
		var res response
		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(err)
		assert.True(len(res.AccessToken) > 0, "should return access_token")
	})
}

func TestUserLogin(t *testing.T) {
	assert := assert.New(t)

	sessMgr := session.NewManager()
	aps := appsensor.NewLoginDetector()

	// Setup Controller.
	userController := controller.NewUser()
	userController.SetAppSensor(aps)
	userController.SetSession(sessMgr)

	// Register an account.
	credentials := &controller.Credentials{
		Email:    "john.doe@mail.com",
		Password: "12345678",
	}
	rr := postRegisterEndpoint(&userController, credentials)
	assert.Equal(200, rr.Code, "should register successfully")

	// Test different scenarios.
	tests := []struct {
		test, email, password, desc string
	}{
		{"login with invalid email", "john.doemail.com", "12345678", "invalid email"},
		{"login with non-existing email", "jane.doe@mail.com", "12345678", "email does not exist"},
		{"login with incorrect email", "john.doe@mail.com", "123456", "password do not match"},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			credentials := &controller.Credentials{
				Email:    tt.email,
				Password: tt.password,
			}
			rr := postLoginEndpoint(&userController, credentials)
			assert.Equal(400, rr.Code)

			var res oidc.ErrorJSON
			err := json.NewDecoder(rr.Body).Decode(&res)
			assert.Nil(err)
			assert.Equal(tt.desc, res.Code)
		})
	}
}

func postRegisterEndpoint(u *controller.User, r *controller.Credentials) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.POST("/register", u.PostRegister)

	jsonBody, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func postLoginEndpoint(u *controller.User, r *controller.Credentials) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.POST("/login", u.PostLogin)

	jsonBody, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
