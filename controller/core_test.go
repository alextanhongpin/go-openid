package controller_test

import (
	"context"
	"io"
	"log"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/controller"
	"github.com/alextanhongpin/go-openid/internal/core"
	"github.com/alextanhongpin/go-openid/model"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/alextanhongpin/go-openid/testdata"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestPostAuthorize(t *testing.T) {
	assert := assert.New(t)
	req := &openid.AuthenticationRequest{
		ClientID:     "1",
		Scope:        "openid",
		ResponseType: "code",
		RedirectURI:  "http://client.example/cb",
		State:        "xyz",
	}

	// Setup context with injected values.
	ctx := openid.SetUserIDContextKey(context.Background(), "john.doe@mail.com")

	// Setup model behaviours.
	model := testdata.NewCoreModel()
	model.On("ValidateAuthnUser", ctx, req).Return(nil)
	model.On("ValidateAuthnRequest", req).Return(nil)
	model.On("ValidateAuthnClient", req).Return(nil)

	u := querystring.Encode(url.Values{}, req)
	rr := corecurl(&model, "POST", "/authorize?"+u.Encode(), nil)

	assert.Equal(302, rr.Code, "should equal response status found")
	log.Println(rr.Body.String(), rr.Header().Get("Location"))
}

func corecurl(model model.Core, method, endpoint string, payload io.Reader) *httptest.ResponseRecorder {

	ctx := context.Background()
	rr := httptest.NewRecorder()

	// Setup services.
	svc := core.NewService(model)

	// Create a fake session with the following email.
	sess := session.NewManager()
	sess.SetSession(rr, "john.doe@mail.com")

	ctl := controller.NewCore(
		controller.CoreSession(sess),
		controller.CoreService(&svc),
	)

	router := httprouter.New()
	router.POST("/authorize", ctl.PostAuthorize)

	req := httptest.NewRequest(method, endpoint, payload)
	// Set the cookie.
	req.Header.Set("Cookie", rr.HeaderMap["Set-Cookie"][0])
	req = req.WithContext(ctx)

	router.ServeHTTP(rr, req)
	return rr
}
