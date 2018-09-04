package oidc

import (
	"errors"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claims represents the OpenIDConnect claims.
type Claims struct {
	*jwt.StandardClaims
	UserID string `json:"user_id,omitempty"`
}

// IDToken is a security token that contains Claims about the Authentication of
// an End-User by and Authorization Server when using Client, and potentially
// other requested Claims.
type IDToken struct {
	Issuer          string   `json:"iss"`
	Subject         string   `json:"sub"`
	Audience        string   `json:"aud"`
	ExpiresIn       int64    `json:"exp"`
	IssuedAt        int64    `json:"iat"`
	AuthTime        int64    `json:"auth_time"`
	Nonce           string   `json:"nonce"`
	AuthCtxClassRef string   `json:"acr"`
	AuthMethodRefs  []string `json:"amr"`
	AuthorizedParty string   `json:"azp"`
}

// Validate performs validation on required fields.
func (i *IDToken) Validate() error {
	if strings.TrimSpace(i.Issuer) == "" {
		return errors.New("issuer cannot be empty")
	}
	if strings.TrimSpace(i.Subject) == "" {
		return errors.New("subject cannot be empty")
	}
	if strings.TrimSpace(i.Audience) == "" {
		return errors.New("audience cannot be empty")
	}
	if i.ExpiresIn < 1 {
		return errors.New("exp cannot be zero")
	}
	if i.IssuedAt < 1 {
		return errors.New("issued at date cannot be zero")
	}
	return nil
}

// Display represents the authentication display options.
type Display int

const (
	Page Display = iota
	Popup
	Touch
	Wap
)

// String fulfils the Stringer interface for Display.
func (d Display) String() string {
	return [...]string{"page", "popup", "touch", "wap"}[d]
}

type Prompt int

const (
	None Prompt = iota
	Login
	Consent
	SelectAccount
)

func (p Prompt) String() string {
	return [...]string{"none", "login", "consent", "select_account"}[p]
}

// AuthenticationRequest is an OAuth 2.0 Authorization Request that requests that the End User be authenticated by the Authorization Server.
type AuthenticationRequest struct {
	Scope        string
	ResponseType string
	ClientID     string
	RedirectURI  string
	State        string
	ResponseMode string
	Nonce        string
	Display      Display
	Prompt       Prompt
	MaxAge       int64
	UILocales    string
	IDTokenHint  string
	LoginHint    string
	AcrValues    string
}

func (a *AuthenticationRequest) Validate() error {
	if a.Scope == "" {
		return errors.New("invalid scope")
	}
	// if !a.Scope.Contain("oidc") {
	// 	return errors.New("invalid scope")
	// }
	// Check other required values
	return nil
}

// AuthenticationResponse is an OAuth 2.0 Authorization Response message returned from the OP's Authorization Endpoint in response to the Authorization Request message sent by the RP.
type AuthenticationResponse struct {
	AccessToken string
	TokenType   string
	IDToken     IDToken
	State       string
	ExpiresIn   int64
}

// ErrorResponse represents the error response parameters
type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
	State            string `json:"state"`
}
type Address struct {
	Formatted     string
	StreetAddress string
	Locality      string
	Region        string
	PostalCode    string
	Country       string
}

type Email struct {
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}

type Phone struct {
	PhoneNumber         string `json:"phone_number,omitempty"`
	PhoneNumberVerified bool   `json:"phone_number_verified,omitempty"`
}

type Profile struct {
	Sub               string
	Name              string
	GivenName         string
	FamilyName        string
	MiddleName        string
	Nickname          string
	PreferredUsername string
	Profile           string
	Picture           string
	Website           string
	Gender            string
	Birthdate         string
	ZoneInfo          string
	Locale            string
	UpdatedAt         int64
}

type StandardClaims struct {
	Profile *Profile
	Email   *Email
	Address *Address
	Phone   *Phone
}

type UserInfoRequest struct{}
type UserInfoResponse struct{}

type Scope int

const (
	ProfileScope Scope = iota
	EmailScope
	AddrScope
	PhoneScope
	OpenIDScope
)

var scopes =  [...]string{"profile", "email", "address", "phone", "openid"}
func (s Scope) String() string {
	return scopes[s] 
}

func (s Scope) Contains (scope string) bool {
	for _, ss := range scopes {
		if ss == scope {
			return true
		}
	} 
	return false
}


// func makeAccessTokenRequest () {
// 	t := &http.Transport{
// 		Dial: (&net.Dialer{
// 			Timeout: 5 *time.Second,
// 			KeepAlive: 5 *time.Second,
// 		}).Dial,
// 		TLSHandshakeTimeout: 5 * time.Second,
// 		ResponseHeaderTimeout: 5 * time.Second,
// 		ExpectContinueTimeout: 1  * time.Second,
// 	}
// 	client := &http.Client {
// 		Timeout: 10 * time.Second,
// 		Transport: t,
// 	}
// 	// res, err := client.Get(url)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	req, err := http.NewRequest("GET", "url", nil)
// 	if err != nil {
// 	log.Fatal(err)
// 	}
// 	req = req.WithContext(ctx)
// 	res, err := client.Do(req)
// 	if err != nil {
//
// 	}
// 	defer res.Body.Close()
//
// }
//
// func HandleAccessTokenRequest (w http.ResponseWriter, r *http.Request) {
//
// 	r.Header().Get("Authorization")
// 	r.Header().Get("Content-Type") == "application/x-www-form-urlencoded"
//
// 	var req AccessTokenRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 	}
// 	if err := req.Validate(); err != nil {
// 	}
// 	if err := FindCodeInCache(req.Code, req.ClientID) {
//
// 	}
// 	cdb, err := FindClient(req.ClientID)
// 	if err != nil {
//
// 	}
// 	if cdb.RedirectURI != req.RedirectURI {}
//
// 	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
// 	w.Header().Set("Cache-Control", "no-store")
// 	w.Header().Set("Pragma", "no-cache")
// }
