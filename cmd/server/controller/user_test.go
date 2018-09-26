package controller_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/alextanhongpin/go-openid/cmd/server/controller"
	"github.com/alextanhongpin/go-openid/pkg/appsensor"
	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	assert := assert.New(t)

	sessMgr := session.NewManager()
	aps := appsensor.NewLoginDetector()
	// Setup Controller.

	userController := controller.NewUser()
	userController.SetAppSensor(aps)
	// userController.SetTemplate(tpl)
	userController.SetSession(sessMgr)

	t.Run("call POST /register endpoint", func(t *testing.T) {
		credentials := &controller.Credentials{
			Email:    "john.doe@mail.com",
			Password: "123456",
		}
		rr := postRegisterEndpoint(&userController, credentials)
		assert.Equal(200, rr.Code)
		log.Println(rr.Body.String())
	})

	t.Run("call POST /login endpoint", func(t *testing.T) {
		credentials := &controller.Credentials{
			Email:    "john.doe@mail.com",
			Password: "123456",
		}
		rr := postLoginEndpoint(&userController, credentials)
		assert.Equal(200, rr.Code)
		log.Println(rr.Body.String())
	})
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
