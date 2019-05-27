package openid

type GrantType string

func (g GrantType) String() string {
	return string(g)
}

func (g GrantType) Equal(grantType string) bool {
	return string(g) == grantType
}

var AuthorizationCode GrantType = "authorization_code"
