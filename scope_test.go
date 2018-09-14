package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope(t *testing.T) {
	assert := assert.New(t)

	s := "address email"
	si := checkScope(s)

	assert.True(si.Has(ScopeAddress))
	assert.True(si.Has(ScopeEmail))
	assert.True(!si.Has(ScopeOpenID))

	s = "  "
	si = checkScope(s)

	assert.True(si.Has(ScopeNone))

	s = " profile address eml openid"
	si = checkScope(s)

	assert.True(si.Has(ScopeAddress))
	assert.True(si.Has(ScopeProfile))
	assert.True(si.Has(ScopeOpenID))
}
