package user_test

import (
	"testing"

	"github.com/alextanhongpin/go-openid/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	assert := assert.New(t)

	model := user.NewModel()

	tests := []struct {
		email    string
		password string
		message  string
	}{
		{"j", "", "password cannot be empty"},
		{"", "", "password cannot be empty"},
		{"a@b.com", "", "password cannot be empty"},
		{"john.doe@ mail", "", "password cannot be empty"},
		{"john.doe@mail.com", "", "password cannot be empty"},

		{"j", "xyz", "password cannot be empty"},
		{"", "xyz", "password cannot be empty"},
		// {"a@b.com", "xyz", "password cannot be empty"},
		// {"john.doe@ mail", "xyz", "password cannot be empty"},
	}

	var err error
	for _, tt := range tests {
		_, err = model.NewUser(tt.email, tt.password)
		assert.NotNil(err)

		// assert.Equal(tt.message, err.Error())
	}
}
