package controller_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/controller"
	"github.com/alextanhongpin/go-openid/pkg/appsensor"
	"github.com/alextanhongpin/go-openid/pkg/session"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestUserRegister(t *testing.T) {
	assert := assert.New(t)

	t.Run("register with invalid json", func(t *testing.T) {
		rr := curl("POST", "/register", bytes.NewBuffer([]byte(`hello world`)))
		assert.Equal(http.StatusBadRequest, rr.Code)

		var res openid.ErrorJSON
		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(err)
		assert.True(len(res.Code) > 0)
	})

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
			js, err := json.Marshal(creds)
			assert.Nil(err)
			rr := curl("POST", "/register", bytes.NewBuffer(js))
			assert.Equal(http.StatusBadRequest, rr.Code)

			var res openid.ErrorJSON
			err = json.NewDecoder(rr.Body).Decode(&res)
			assert.Nil(err)
			assert.Equal(tt.desc, res.Code)
		})
	}

	t.Run("register with valid email and password", func(t *testing.T) {
		creds := &controller.Credentials{
			Email:    "john.doe@mail.com",
			Password: "12345678",
		}
		js, err := json.Marshal(creds)
		assert.Nil(err)
		rr := curl("POST", "/register", bytes.NewBuffer(js))
		assert.Equal(http.StatusOK, rr.Code)
		type response struct {
			AccessToken string `json:"access_token"`
		}

		var res response
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(err)
		assert.True(len(res.AccessToken) > 0, "should return access_token")
	})
}

func TestUserLogin(t *testing.T) {
	assert := assert.New(t)

	// Register an account.
	creds := &controller.Credentials{
		Email:    "john.doe@mail.com",
		Password: "12345678",
	}

	js, err := json.Marshal(creds)
	assert.Nil(err)
	rr := curl("POST", "/register", bytes.NewBuffer(js))
	assert.Equal(http.StatusOK, rr.Code, "should register successfully")

	// Test different scenarios.
	tests := []struct {
		test, email, password, desc string
	}{
		{"login with invalid email", "john.doemail.com", "12345678", "invalid email"},
		{"login with incorrect/non-existing email", "jane@mail.com", "12345678", "email does not exist"},
		{"login with incorrect password", "john.doe@mail.com", "123456", "email does not exist"},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			creds := &controller.Credentials{
				Email:    tt.email,
				Password: tt.password,
			}
			js, err := json.Marshal(creds)
			assert.Nil(err)
			rr := curl("POST", "/login", bytes.NewBuffer(js))
			assert.Equal(http.StatusBadRequest, rr.Code)

			var res openid.ErrorJSON
			err = json.NewDecoder(rr.Body).Decode(&res)
			assert.Nil(err)
			assert.Equal(tt.desc, res.Code)
		})
	}
}

func newController() controller.User {
	c := controller.NewUser(
		controller.UserAppSensor(appsensor.NewLoginDetector()),
		controller.UserSession(session.NewManager()),
	)
	return c
}

func curl(method, endpoint string, payload io.Reader) *httptest.ResponseRecorder {
	ctl := newController()

	router := httprouter.New()
	router.POST("/register", ctl.PostRegister)
	router.POST("/login", ctl.PostLogin)

	req := httptest.NewRequest(method, endpoint, payload)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
