package oidc

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

var UserContextKey = ContextKey("user_id")
