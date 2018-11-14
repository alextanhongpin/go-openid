package user_test

import (
	"testing"

	"github.com/alextanhongpin/go-openid/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestService_RegistrationError(t *testing.T) {
	assert := assert.New(t)

	// Creates a new user service.
	svc := user.New()

	tests := []struct {
		email, password, errMsg string
	}{
		{"", "", "email is required"},
		{"john", "", `"john" is not a valid email`},
		{"a.b@mail.com", "", "password is required"},
		{"a.b@mail.com", "12345", "password length must be larger than or equal to 6"},
	}
	var err error
	for _, tt := range tests {
		_, err = svc.Register(tt.email, tt.password)
		assert.NotNil(err)
		assert.Equal(tt.errMsg, err.Error())
	}
}
func TestService_RegisterTwice(t *testing.T) {
	assert := assert.New(t)

	// Creates a new user service.
	svc := user.New()
	_, err := svc.Register("john.doe@mail.com", "123456")
	assert.Nil(err)

	_, err = svc.Register("john.doe@mail.com", "123456")
	assert.NotNil(err)
	assert.Equal("user already exist", err.Error())
}
