package openid

import "errors"

// IDToken is a security token that contains Claims about the Authentication of an End-User by and Authorization Server when using Client, and potentially other requested Claims.
type IDToken struct {
	Iss      string   `json:"iss"`
	Sub      string   `json:"sub"`
	Aud      string   `json:"aud"`
	Exp      int64    `json:"exp"`
	Iat      int64    `json:"iat"`
	AuthTime int64    `json:"auth_time"`
	Nonce    string   `json:"nonce"`
	Acr      string   `json:"acr"`
	Amr      []string `json:"amr"`
	Azp      string   `json:"azp"`
}

type Display int

const (
	Page Display = iota
	Popup
	Touch
	Wap
)

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

// Authentication Error Response
var (
	ErrInteractionRequired      = errors.New("interaction required")
	ErrLoginRequired            = errors.New("login required")
	ErrAccountSelectionRequired = errors.New("account selection required")
	ErrConsentRequired          = errors.New("consent required")
	ErrInvalidRequestURI        = errors.New("invalid request uri")
	ErrInvalidRequestObject     = errors.New("invalid request object")
	ErrRequestNotSupported      = errors.New("request not supported")
	ErrRequestURINotSupported   = errors.New("request uri not supported")
	ErrRegistrationNotSupported = errors.New("registration not supported")
)

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

type TokenRequest struct {
}

type TokenResponse struct {
}

type Address struct {
	Formatted     string
	StreetAddress string
	Locality      string
	Region        string
	PostalCode    string
	Country       string
}

type StandardClaims struct {
	Sub                 string
	Name                string
	GivenName           string
	FamilyName          string
	MiddleName          string
	Nickname            string
	PreferredUsername   string
	Profile             string
	Picture             string
	Website             string
	Email               string
	EmailVerified       bool
	Gender              string
	Birthdate           string
	ZoneInfo            string
	Locale              string
	PhoneNumber         string
	PhoneNumberVerified bool
	Address             Address
	UpdatedAt           int64
}

type UserInfoRequest struct{}
type UserInfoResponse struct{}

type Scope int

const (
	Profile Scope = iota
	Email
	Addr
	Phone
)

func (s Scope) String() string {
	return [...]string{"profile", "email", "address", "phone"}[s]
}

type RefreshRequest struct {
}
type RefreshResponse struct{}
