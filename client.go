package openid

type Client struct {
	RedirectURIs                 string
	ResponseTypes                []string
	GrantTypes                   []string
	ApplicationType              string
	Contacts                     []string
	ClientName                   string
	LogoURI                      string
	ClientURI                    string
	PolicyURI                    string
	TosURI                       string
	JwksURI                      string
	Jwks                         string
	SectorIdentifierURI          string
	SubjectType                  string
	IDTokenSignedResponseAlg     string
	IDTokenEncryptedResponseAlg  string
	IDTokenEncryptedResponseEnc  string
	UserinfoSignedResponseAlg    string
	UserinfoEncryptedResponseAlg string
	UserinfoEncryptedResponseEnc string
	RequestObjectSigningAlg      string
	RequestObjectEncryptionAlg   string
	RequestObjectEncryptionEnc   string
	TokenEndpointAuthMethod      string
	TokenEndpointAuthSigningAlg  string
	DefaultMaxAge                int64
	RequireAuthTime              int64
	DefaultAcrValues             string
	InitiateLoginURI             string
	RequestURIs                  []string
}
