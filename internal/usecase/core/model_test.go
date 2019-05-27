package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/core"
	"github.com/alextanhongpin/go-openid/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestValidateAuthnRequest(t *testing.T) {
	assert := assert.New(t)
	model := core.NewModel()

	req := &openid.AuthenticationRequest{
		ClientID:     "hello",
		RedirectURI:  "http://client.example.com/cb",
		ResponseType: "code",
		Scope:        "openid",
	}

	t.Run("validate nil field", func(t *testing.T) {
		err := model.ValidateAuthnRequest(nil)
		assert.NotNil(err, "should handle nil arguments")
	})

	t.Run("validate required fields", func(t *testing.T) {
		err := model.ValidateAuthnRequest(req)
		assert.Nil(err, "should succeed with valid arguments")
	})

	t.Run("validate missing client_id", func(t *testing.T) {
		copy := *req
		copy.ClientID = ""
		err := model.ValidateAuthnRequest(&copy)
		if verr, ok := err.(*openid.ErrorJSON); ok {
			var (
				code = "invalid_request"
				desc = "client_id is required"
			)
			assert.Equal(code, verr.Code)
			assert.Equal(desc, verr.Description)
		} else {
			assert.True(ok, "should return custom error")
		}
	})

	t.Run("validate missing redirect_uri", func(t *testing.T) {
		copy := *req
		copy.RedirectURI = ""
		err := model.ValidateAuthnRequest(&copy)
		if verr, ok := err.(*openid.ErrorJSON); ok {
			var (
				code = "invalid_request"
				desc = "redirect_uri is required"
			)
			assert.Equal(code, verr.Code)
			assert.Equal(desc, verr.Description)
		} else {
			assert.True(ok, "should return custom error")
		}
	})

	t.Run("validate response_type", func(t *testing.T) {
		tests := []struct {
			responseType string
			valid        bool
		}{
			{"code", true},
			{"abc", false},
			{"", false},
		}
		for _, tt := range tests {
			copy := *req
			copy.ResponseType = tt.responseType
			err := model.ValidateAuthnRequest(&copy)
			if tt.valid {
				assert.Nil(err)
			} else {
				assert.NotNil(err)
			}
		}
	})

	t.Run("validate scope", func(t *testing.T) {
		copy := *req
		copy.Scope = ""
		err := model.ValidateAuthnRequest(&copy)
		if verr, ok := err.(*openid.ErrorJSON); ok {
			var (
				code = "invalid_request"
				desc = "scope is required"
			)
			assert.Equal(code, verr.Code)
			assert.Equal(desc, verr.Description)
		} else {
			assert.True(ok)
		}
	})
}

func TestClientValidation(t *testing.T) {
	assert := assert.New(t)

	// Setup repository
	client := database.NewClientKV()
	client.Put("app", &openid.Client{
		ClientID:     "app",
		RedirectURIs: []string{"http://client.example.com/cb"},
	})

	// Setup model.
	model := core.NewModel()
	model.SetClient(client)

	req := &openid.AuthenticationRequest{
		ClientID:     "app",
		RedirectURI:  "http://client.example.com/cb",
		ResponseType: "code",
		Scope:        "openid",
	}

	t.Run("validate existing client", func(t *testing.T) {
		err := model.ValidateAuthnClient(req)
		assert.Nil(err)
	})

	t.Run("validate non-existing client", func(t *testing.T) {
		copy := *req
		copy.ClientID = "null"
		err := model.ValidateAuthnClient(&copy)
		assert.NotNil(err)
		assert.Equal("client does not exist", err.Error())
	})

	t.Run("validate invalid client redirect_uri", func(t *testing.T) {
		copy := *req
		copy.RedirectURI = "http://unknown"
		err := model.ValidateAuthnClient(&copy)
		assert.NotNil(err)
		assert.Equal("redirect_uri incorrect", err.Error())
	})
}

func TestUserValidation(t *testing.T) {
	assert := assert.New(t)

	var (
		userID = "1"
	)

	// Setup repository.
	user := database.NewUserKV()
	user.Put(userID, &openid.User{
		Profile: openid.Profile{
			UpdatedAt: time.Now().UTC().Unix(),
		},
	})

	// Setup model.
	model := core.NewModel()
	model.SetUser(user)

	// Setup request.
	req := &openid.AuthenticationRequest{
		ClientID:     "hello",
		RedirectURI:  "http://client.example.com/cb",
		ResponseType: "code",
		Scope:        "openid",
		MaxAge:       int64((1 * time.Hour).Seconds()),
	}

	t.Run("validate existing user", func(t *testing.T) {
		ctx := context.Background()
		ctx = openid.SetUserIDContextKey(ctx, userID)
		err := model.ValidateAuthnUser(ctx, req)
		assert.Nil(err)
	})

	t.Run("validate non-existing user", func(t *testing.T) {
		ctx := context.Background()
		err := model.ValidateAuthnUser(ctx, req)
		assert.NotNil("user_id missing", err.Error())
	})
}
