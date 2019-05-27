package openid

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
