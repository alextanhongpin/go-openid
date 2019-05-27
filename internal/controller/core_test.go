package controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/controller"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/alextanhongpin/go-openid/service"
	"github.com/alextanhongpin/go-openid/testdata"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestPostAuthorize(t *testing.T) {
	assert := assert.New(t)

	var (
		req = &openid.AuthenticationRequest{
			ClientID:     "1",
			Scope:        "openid",
			ResponseType: "code",
			RedirectURI:  "http://client.example/cb",
			State:        "xyz",
		}
		res = &openid.AuthenticationResponse{
			Code:  "code123",
			State: "xyz",
		}
	)

	// Setup context with injected values.
	ctx := context.Background()
	ctx = openid.SetUserIDContextKey(ctx, "john.doe@mail.com")

	s := testdata.NewCoreService()
	t.Run("call with valid parameters and invalid session", func(t *testing.T) {
		s.On("Authenticate", context.Background(), req).Return(res, nil)
		u := querystring.Encode(url.Values{}, req)
		rr := corecurl(&s, false, "POST", "/authorize?"+u.Encode(), nil)

		var (
			code = http.StatusBadRequest
			msg  = "http: named cookie not present"
		)

		assert.Equal(code, rr.Code, "should equal response status found")

		var res openid.ErrorJSON
		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(err)
		assert.Equal(msg, res.Code)
	})

	t.Run("call with valid parameters", func(t *testing.T) {
		s.On("Authenticate", ctx, req).Return(res, nil)
		u := querystring.Encode(url.Values{}, req)
		rr := corecurl(&s, true, "POST", "/authorize?"+u.Encode(), nil)

		var (
			code     = http.StatusFound
			location = "http://client.example/cb?code=code123&state=xyz"
		)

		assert.Equal(code, rr.Code, "should equal response status found")
		assert.Equal(location, rr.Header().Get("Location"))
	})

	t.Run("call with empty requests", func(t *testing.T) {
		s.On("Authenticate", ctx, nil).Return(nil, errors.New("bad request"))
		rr := corecurl(&s, true, "POST", "/authorize", nil)

		var (
			code = http.StatusUnprocessableEntity
		)
		assert.Equal(code, rr.Code, "should return status 400 - Bad Request")

		var res openid.ErrorJSON
		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(err)
		assert.Equal("request is empty", res.Code)
	})

}

func corecurl(svc service.Core, enableSession bool, method, endpoint string, payload io.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	// Create a fake session with the following email.
	sess := session.NewManager()
	if enableSession {
		sess.SetSession(rr, "john.doe@mail.com")
	}

	ctl := controller.NewCore(
		controller.CoreSession(sess),
		controller.CoreService(svc),
	)

	router := httprouter.New()
	router.POST("/authorize", ctl.PostAuthorize)

	ctx := context.Background()
	req := httptest.NewRequest(method, endpoint, payload)
	req = req.WithContext(ctx)
	// Set the cookie.
	if enableSession {
		req.Header.Set("Cookie", rr.HeaderMap["Set-Cookie"][0])
	}

	router.ServeHTTP(rr, req)
	return rr
}
