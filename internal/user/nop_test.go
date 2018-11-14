package user_test

import (
	"testing"

	"github.com/alextanhongpin/go-openid/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestNopModel(t *testing.T) {
	assert := assert.New(t)
	nop := new(user.NopModel)
	user, err := nop.NewUser("email", "password")
	assert.Nil(user, err)

	err = nop.ValidateEmail("email")
	assert.Nil(err)

	err = nop.ValidatePassword("password")
	assert.Nil(err)

	err = nop.ValidateLimit(-1000)
	assert.Nil(err)
}
