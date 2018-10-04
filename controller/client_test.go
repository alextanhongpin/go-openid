package controller_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/controller"
	"github.com/alextanhongpin/go-openid/service"
	"github.com/alextanhongpin/go-openid/testdata"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	assert := assert.New(t)

	var (
		clientID     = "1"
		clientSecret = "secret"
		client       = openid.Client{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}
	)
	s := testdata.NewClientService()
	s.On("Read", clientID).Return(&client, nil)
	rr := curlClient(s, "GET", "/clients/"+clientID, nil)

	assert.Equal(http.StatusOK, rr.Code)

	var res openid.Client
	err := json.NewDecoder(rr.Body).Decode(&res)
	assert.Nil(err)
	assert.Equal(clientID, res.ClientID)
	assert.Equal(clientSecret, res.ClientSecret)
}

func curlClient(service service.Client, method, endpoint string, payload io.Reader) *httptest.ResponseRecorder {
	ctl := controller.NewClient(controller.ClientService(service))

	router := httprouter.New()
	router.GET("/clients/:client_id", ctl.GetClient)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(method, endpoint, payload)
	router.ServeHTTP(rr, req)
	return rr
}
