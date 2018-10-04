package openid_test

import (
	"testing"

	openid "github.com/alextanhongpin/go-openid"

	"github.com/stretchr/testify/assert"
)

func TestAuthRequest(t *testing.T) {
	assert := assert.New(t)
	req := openid.AuthenticationRequest{}
	prompt := req.GetPrompt()
	assert.True(prompt.Is(openid.PromptNone))
}
