package openid

import "strings"

// -- scopes

// By using bits, we save on a lot of checking. For example, we would not face
// the issue of duplicate scopes "email email". Also, it's possible to check if
// the scopes are an exact match, or contains one of the scope in the list.

// Scope represents the openid scope.
type Scope int

const (
	ScopeNone Scope = 1 << iota
	ScopeAddress
	ScopeEmail
	ScopeOpenID
	ScopePhone
	ScopeProfile
)

// Has returns true if the given scope is in the list of scope. Can have multiple scope.
func (s Scope) Has(ss Scope) bool {
	return s&ss != 0
}

// Is returns true if the scope is exactly one of the defined scope.
func (s Scope) Is(ss Scope) bool {
	return s&ss == s|ss
}

var scopemap = map[string]Scope{
	"address": ScopeAddress,
	"email":   ScopeEmail,
	"openid":  ScopeOpenID,
	"phone":   ScopePhone,
	"profile": ScopeProfile,
}

func NewScope(scope string) (i Scope) {
	scopes := strings.Split(scope, " ")
	for _, s := range scopes {
		if v, exist := scopemap[s]; exist {
			i |= v
		}
	}
	if i == 0 {
		i |= ScopeNone
	}
	return
}
