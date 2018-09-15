package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/api"
)

func TestClientError(t *testing.T) {
	assert := assert.New(t)

	_, err := api.RegisterClient(nil)

	// Perform type assertion back to the JSON error.
	cerr, ok := err.(*oidc.ErrorJSON)
	assert.True(ok)

	var (
		code = "invalid_client_metadata"
		desc = "client metadata cannot be empty"
	)
	assert.Equal(code, cerr.Code)
	assert.Equal(desc, cerr.Description)
}
