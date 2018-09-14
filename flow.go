package oidc

type Flow int

const (
	AuthorizationCodeFlow Flow = 1 << iota
	ImplicitFlow
	HybridFlow
	UnknownFlow
)

var flows = map[string]Flow{
	"code":                AuthorizationCodeFlow,
	"id_token":            ImplicitFlow,
	"id_token token":      ImplicitFlow,
	"code id_token":       HybridFlow,
	"code id_token token": HybridFlow,
}

// CheckFlow returns the current OIDC registration flow.
func CheckFlow(responseType string) Flow {
	t := sortstr(responseType)
	if f, ok := flows[t]; ok {
		return f
	}
	return UnknownFlow
}
