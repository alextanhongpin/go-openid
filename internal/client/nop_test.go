package client_test

import (
	"testing"

	"github.com/alextanhongpin/go-openid/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNopService(t *testing.T) {
	assert := assert.New(t)
	nop := new(client.NopService)
	client, err := nop.Register(nil)
	assert.Nil(client)
	assert.Nil(err)
}

func TestNopModel(t *testing.T) {
	assert := assert.New(t)
	nop := new(client.NopModel)
	assert.Equal("", nop.GenerateClientID())

	clientSecret, err := nop.GenerateClientSecret()
	assert.Nil(err)
	assert.Equal(clientSecret, "")

	registrationAccessToken, err := nop.GenerateRegistrationAccessToken("")
	assert.Nil(err)
	assert.Equal(registrationAccessToken, "")

	assert.Equal(int64(0), nop.GenerateClientIDIssuedAt())
	assert.Equal("", nop.GenerateRegistrationClientURI())

	assert.Nil(nop.ValidateClient(nil))
	assert.Nil(nop.ValidateClientResponse(nil))
}

func TestNopRepository(t *testing.T) {
	assert := assert.New(t)
	nop := new(client.NopRepository)
	c, err := nop.Get("")
	assert.Nil(c, err)

	err = nop.Put("", nil)
	assert.Nil(err)

	has := nop.Has("")
	assert.Equal(false, has)

	clients := nop.List(1000)
	assert.Equal(0, len(clients))

	client, err := nop.GetByCredentials("", "")
	assert.Nil(client)
	assert.Nil(err)
}
