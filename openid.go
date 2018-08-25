package openid

import (
	"errors"
	"net/http"
	"net/url"
)

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
	// if !a.Scope.Contain("openid") {
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

type AuthorizationRequest struct {
	ResponseType string `json:"response_type"`
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
}

func (r *AuthorizationRequest) Validate() error {
	// Required fields
	if r.ResponseType != "code" {
		return ErrUnsupportedResponseType
	}
	if r.ClientID == "" {
	}
	// Optional fields
	if r.RedirectURI == "" {
	}
	if r.Scope == "" {
		return ErrInvalidScope
	}
	if r.State == "" {
	}
	return nil
}

type AuthorizationResponse struct {
	Scope string `json:"scope"`
	State string `json:"state"`
}

// type AuthorizationError struct {
// 	Error            string `json:"error"`
// 	ErrorDescription string `json:"error_description"`
// 	ErrorURI         string `json:"error_uri"`
// }

func HandleAuthorizationRequest(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	req := AuthorizationRequest{
		ResponseType: q.Get("response_type"),
		ClientID:     q.Get("client_id"),
		RedirectURI:  q.Get("redirect_uri"),
		Scope:        q.Get("scope"),
		State:        q.Get("state"),
	}
	// Validate fields
	if err := req.Validate(); err != nil {
		// Return json error
	}
	// Create a new authorization code, store it to the cache and return the authorization code
	code := "newAuthorizationCodeLifetime10minutes"
	// Check if the client exist based on the client id
	// CheckClientExist(req.ClientID)
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
	}
	// u.Scheme = "https"
	// u.Host = "google"
	qq := u.Query()
	qq.Set("code", code)
	qq.Set("state", req.State)
	u.RawQuery = qq.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// Authorization errors
var (
	ErrInvalidRequest          = errors.New("invalid request")
	ErrUnauthorizedClient      = errors.New("unauthorized client")
	ErrAccessDenied            = errors.New("access denied")
	ErrUnsupportedResponseType = errors.New("unsupported response type")
	ErrInvalidScope            = errors.New("invalid scope")
	ErrServerError             = errors.New("server error")
	ErrTemporarilyUnavailable  = errors.New("temporarily unavailable")
)

type AuthorizationError struct {
	Error            string
	ErrorDescription string
	ErrorURI         string
	State            string
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

func validateContentType() {
	// Check if the contentType == "application/x-www-form-urlencoded"
}

type AccessTokenRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
	ClientID    string `json:"client_id"`
}

func (r *AccessTokenRequest) Sanitize() error {
	if r.GrantType != "authorization_code" {
		r.GrantType = "authorization_code"
	}
	return nil
}
func (r *AccessTokenRequest) Validate() error {
	if r.GrantType != "authorization_code" {
		return ErrInvalidRequest
	}
	// Check required field
	if r.Code == "" {
	}
	if r.RedirectURI == "" {
	}
	if r.ClientID == "" {
	}
	return nil
}

type AccessTokenResponse struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int64
	RefreshToken string
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
