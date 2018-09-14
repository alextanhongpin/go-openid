package oidc

import "strings"

type Scope int

const (
	ScopeNone Scope = 1 << iota
	ScopeAddress
	ScopeEmail
	ScopeOpenID
	ScopePhone
	ScopeProfile
)

func (s Scope) Has(ss Scope) bool {
	return s&ss != 0
}

var scopemap = map[string]Scope{
	"address": ScopeAddress,
	"email":   ScopeEmail,
	"openid":  ScopeOpenID,
	"phone":   ScopePhone,
	"profile": ScopeProfile,
}

func checkScope(scope string) (i Scope) {
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

//
// type Scopes [5]string
//
// var scopes Scopes = [...]string{"profile", "email", "address", "phone", "openid"}
//
// func (s Scopes) Contains(scope string) bool {
//         for _, ss := range s {
//                 if ss == scope {
//                         return true
//                 }
//         }
//         return false
// }
