package openid

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

func TestResponseType(t *testing.T) {
	assert := assert.New(t)

	var (
		code    = ResponseTypeCode
		token   = ResponseTypeToken
		idToken = ResponseTypeIDToken
	)

	t.Run("all code", func(t *testing.T) {
		s := "code id_token token code id_token"
		p := parseResponseType(s)

		// Check individuals.
		assert.True(p.Has(code))
		assert.True(p.Has(token))
		assert.True(p.Has(idToken))

		// Check combinations.
		assert.True(p.Has(code | token))
		assert.True(p.Has(token | idToken))
		assert.True(p.Has(code | token | idToken))

		// Check exact match.
		assert.True(!p.Is(code))
		assert.True(!p.Is(code | token))
		assert.True(p.Is(code | token | idToken))

		q := parseResponseType(s)
		assert.True(p.Is(q))
		assert.True(p.Has(q))
		assert.True(q.Is(p))
		assert.True(q.Has(p))
	})
	t.Run("only one", func(t *testing.T) {
		s := " code "
		p := parseResponseType(s)

		// Check individuals.
		assert.True(p.Has(code))
		assert.True(!p.Has(token))
		assert.True(!p.Has(idToken))

		// Check combinations.
		assert.True(p.Has(code | token))
		assert.True(!p.Has(token | idToken))
		assert.True(p.Has(code | token | idToken))

		// Check exact match.
		assert.True(p.Is(code))
		assert.True(!p.Is(code | token))
		assert.True(!p.Is(code | token | idToken))
	})
}
