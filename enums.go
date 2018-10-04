package openid

import (
	"strings"
)

// Prompt represents the enums for prompting.
type Prompt int

const (
	PromptNone Prompt = 1 << iota
	PromptConsent
	PromptLogin
	PromptSelectAccount
)

// Has returns true if the prompt value is one the prompt enums.
func (p Prompt) Has(pp Prompt) bool {
	return p&pp != 0
}

// Is returns true if the prompt valus is exactly one of the prompt enums.
func (p Prompt) Is(pp Prompt) bool {
	return p&pp == p|pp
}

var promptmap = map[string]Prompt{
	"none":           PromptNone,
	"login":          PromptLogin,
	"consent":        PromptConsent,
	"select_account": PromptSelectAccount,
}

func parsePrompt(prompt string) (i Prompt) {
	ps := strings.Split(prompt, " ")
	for _, p := range ps {
		if v, exist := promptmap[p]; exist {
			i |= v
		}
	}
	if i == 0 {
		i |= PromptNone
	}
	return
}

// ResponseType represents the enum for response type.
type ResponseType int

const (
	ResponseTypeNone ResponseType = 1 << iota
	ResponseTypeCode
	ResponseTypeIDToken
	ResponseTypeToken
)

// Has returns true if the response type belong to one of the enum.
func (r ResponseType) Has(rr ResponseType) bool {
	return r&rr != 0
}

// Is returns true if the enum matches exactly one of the enum.
func (r ResponseType) Is(rr ResponseType) bool {
	return r&rr == r|rr
}

var responsetypemap = map[string]ResponseType{
	"code":     ResponseTypeCode,
	"id_token": ResponseTypeIDToken,
	"token":    ResponseTypeToken,
}

func parseResponseType(responseType string) (i ResponseType) {
	rs := strings.Split(responseType, " ")
	for _, r := range rs {
		if v, exist := responsetypemap[r]; exist {
			i |= v
		}
	}
	if i == 0 {
		i |= ResponseTypeNone
	}
	return
}

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

// parseScope parses the scope from a string into the Scope enum.
func parseScope(scope string) (i Scope) {
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

// struct memory allocation = 0, interface{} = 8, bool = 1
var displaymap = map[string]struct{}{
	"page":  struct{}{},
	"popup": struct{}{},
	"touch": struct{}{},
	"wap":   struct{}{},
}

// -- flow

// CheckFlow returns the current openid registration flow.
func CheckFlow(enum ResponseType) string {
	var (
		code    = ResponseTypeCode
		token   = ResponseTypeToken
		idToken = ResponseTypeIDToken
	)
	if enum.Is(code) {
		return "authorization_code"
	}
	if enum.Is(idToken) || enum.Is(idToken|token) {
		return "implicit"
	}
	if enum.Is(code|idToken) || enum.Is(code|idToken|token) {
		return "hybrid"
	}
	return ""
}
