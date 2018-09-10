package passwd_test

import (
	"log"
	"testing"

	"github.com/alextanhongpin/go-openid/pkg/passwd"
	"github.com/stretchr/testify/assert"
)

func TestPasswordHashAndVerify(t *testing.T) {
	assert := assert.New(t)

	var (
		password = "secret"
	)
	hash, err := passwd.Hash(password)
	log.Println(hash)
	assert.Nil(err)

	err = passwd.Verify(password, hash)
	assert.Nil(err)
}

func TestEmptyPassword(t *testing.T) {
	assert := assert.New(t)

	_, err := passwd.Hash("")
	assert.NotNil(err)
}
