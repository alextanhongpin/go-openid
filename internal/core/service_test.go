package core_test

import (
	"testing"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/core"
	database "github.com/alextanhongpin/go-openid/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestServiceAuthenticate(t *testing.T) {
	assert := assert.New(t)

	// Setup repository.
	client := database.NewClientKV()
	client.Put("hello", &oidc.Client{
		ClientID:     "hello",
		RedirectURIs: []string{"http://client.example.com/cb"},
	})

	// Setup model.
	model := core.NewModel()
	model.SetClient(client)

	// Setup service.
	service := core.NewService(&model)

	req := &oidc.AuthenticationRequest{
		ClientID:     "hello",
		RedirectURI:  "http://client.example.com/cb",
		ResponseType: "code",
		Scope:        "openid",
	}

	err := service.PreAuthenticate(req)
	assert.Nil(err)
}
