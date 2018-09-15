package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope(t *testing.T) {
	assert := assert.New(t)

	s := "address email"
	si := parseScope(s)

	assert.True(si.Has(ScopeAddress))
	assert.True(si.Has(ScopeEmail))
	assert.True(!si.Has(ScopeOpenID))

	s = "  "
	si = parseScope(s)
	assert.True(si.Has(ScopeNone))

	s = ""
	si = parseScope(s)
	assert.True(si.Has(ScopeNone))

	s = "a b c xy y z"
	si = parseScope(s)
	assert.True(si.Has(ScopeNone))

	// Can handle duplicates!
	s = "profile profile address eml openid"
	si = parseScope(s)

	assert.True(!si.Has(ScopeNone), "should not have scope none")
	assert.True(si.Has(ScopeAddress))
	assert.True(si.Has(ScopeProfile))
	assert.True(si.Has(ScopeOpenID))
}
