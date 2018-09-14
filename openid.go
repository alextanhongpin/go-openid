package oidc

import (
	"errors"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// IDToken is a security token that contains Claims about the Authentication of
// an End-User by and Authorization Server when using Client, and potentially
// other requested Claims.
type IDToken struct {
	*jwt.StandardClaims
	AuthorizedParty                     string   `json:"azp,omitempty"`       // Authorized party - the party to which the ID Token was issued.
	Nonce                               string   `json:"nonce,omitempty"`     // Value used to associate a Client session with an ID Token.
	AuthTime                            int64    `json:"auth_time,omitempty"` // Time when the authentication occurred.
	AtHash                              string   `json:"at_hash,omitempty"`   // Access Token hash value.
	CodeHash                            string   `json:"c_hash,omitempty"`    // Code hash value.
	AuthenticationContextClassReference string   `json:"acr,omitempty"`       // Authentication context class reference.
	AuthenticationMethodReferences      []string `json:"amr,omitempty"`       // Authentication method references.
	SessionID                           string   `json:"sid,omitempty"`       // Session ID.
	SubJWK                              string   `json:"sub_jwk,omitempty"`   // Public key used to check the signature of an ID Token.
	Address                             *Address
	Email                               *Email
	Phone                               *Phone
	Profile                             *Profile
}

// TODO: Check other libraries to see their validation.
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
	now := time.Now().UTC().Unix()
	if i.ExpiresAt < now {
		return errors.New("time cannot be less then now")
	}
	if i.NotBefore < now {
		return errors.New("must be greater than now")
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
	AcrValues    string  `json:"acr_values,omitempty"`
	ClientID     string  `json:"client_id,omitempty"`
	Display      Display `json:"display,omitempty"`
	IDTokenHint  string  `json:"id_token_hint,omitempty"`
	LoginHint    string  `json:"login_hint,omitempty"`
	MaxAge       int64   `json:"max_age,omitempty"`
	Nonce        string  `json:"nonce,omitempty"`
	Prompt       Prompt  `json:"prompt,omitempty"`
	RedirectURI  string  `json:"redirect_uri,omitempty"`
	ResponseMode string  `json:"response_mode,omitempty"`
	ResponseType string  `json:"response_type,omitempty"`
	Scope        string  `json:"scope,omitempty"`
	State        string  `json:"state,omitempty"`
	UILocales    string  `json:"ui_locales,omitempty"`
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
	AccessToken string  `json:"access_token,omitempty"`
	ExpiresIn   int64   `json:"expires_in,omitempty"`
	IDToken     IDToken `json:"id_token,omitempty"`
	State       string  `json:"state,omitempty"`
	TokenType   string  `json:"token_type,omitempty"`
}

type Address struct {
	Country       string `json:"country,omitempty"`
	Formatted     string `json:"formatted,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
}

type Email struct {
	Email         string `json:"email,omitempty"` // Preferred e-mail address.
	EmailVerified bool   `json:"email_verified"`  // True if the e-mail address has been verified; otherwise false.
}

type Phone struct {
	PhoneNumber         string `json:"phone_number,omitempty"` // Preferred telephone number.
	PhoneNumberVerified bool   `json:"phone_number_verified"`  // True if the phone number has been verified; otherwise false.
}

type Profile struct {
	Birthdate         string `json:"birth_date,omitempty"`         // Birthday.
	FamilyName        string `json:"family_name,omitempty"`        // Surname(s) or first name(s).
	Gender            string `json:"gender,omitempty"`             // Gender.
	GivenName         string `json:"given_name,omitempty"`         // Given name(s) or first name(s).
	Locale            string `json:"locale,omitempty"`             // Locale.
	MiddleName        string `json:"middle_name,omitempty"`        // Middle name(s).
	Name              string `json:"name,omitempty"`               // Full name.
	Nickname          string `json:"nickname,omitempty"`           // Casual name.
	Picture           string `json:"picture,omitempty"`            // Profile picture URL.
	PreferredUsername string `json:"preferred_username,omitempty"` // Shorthand name by which the End-User wishes to be referred to.
	Profile           string `json:"profile,omitempty"`            // Profile page URL.
	UpdatedAt         int64  `json:"updated_at,omitempty"`         // Time the information was last updated.
	ZoneInfo          string `json:"zone_info,omitempty"`          // Time zone.
	Website           string `json:"website,omitempty"`            // Web page or blog URL.
}

type Scope int

const (
	ProfileScope Scope = iota
	EmailScope
	AddrScope
	PhoneScope
	OpenIDScope
)

var scopes = [...]string{"profile", "email", "address", "phone", "openid"}

func (s Scope) String() string {
	return scopes[s]
}

func (s Scope) Contains(scope string) bool {
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
