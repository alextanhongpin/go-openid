package core_test

import (
	"context"
	"testing"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/core"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	assert := assert.New(t)
	model := core.NewModel()

	req := &oidc.AuthenticationRequest{
		ClientID:     "hello",
		RedirectURI:  "http://client.example.com/cb",
		ResponseType: "code",
		Scope:        "openid",
	}

	t.Run("ValidateAuthnRequest", func(t *testing.T) {
		err := model.ValidateAuthnRequest(req)
		assert.Nil(err)
	})

	t.Run("ValidateAuthnUser", func(t *testing.T) {
		ctx := context.Background()
		err := model.ValidateAuthnUser(ctx, req)
		assert.NotNil(err)
	})

	t.Run("ValidateAuthnClient", func(t *testing.T) {
		err := model.ValidateAuthnClient(req)
		assert.NotNil(err)
	})
}
