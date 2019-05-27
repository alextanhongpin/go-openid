package openid

import (
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
	jwt "github.com/dgrijalva/jwt-go"
)

// IDToken is a security token that contains Claims about the Authentication of
// an End-User by and Authorization Server when using Client, and potentially
// other requested Claims.
type IDToken struct {
	jwt.StandardClaims
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

// NewIDToken returns a pointer to a new id token with empty fields.
func NewIDToken() *IDToken {
	return &IDToken{
		StandardClaims: jwt.StandardClaims{},
		Address:        &Address{},
		Email:          &Email{},
		Phone:          &Phone{},
		Profile:        &Profile{},
	}
}

func (i *IDToken) SignHS256(key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, i)
	return token.SignedString(key)
}

func (i *IDToken) ParseHS256(str string, key []byte) error {
	token, err := jwt.ParseWithClaims(str, i, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return err
	}
	var ok bool
	i, ok = token.Claims.(*IDToken)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid id_token")
}

// TODO: Check other libraries to see their validation.

// Validate performs validation on required fields.
func (i *IDToken) Validate(aud, iss, sub string, now int64) error {
	if err := i.StandardClaims.Valid(); err != nil {
		return err
	}
	if err := i.VerifyIssuer(iss); err != nil {
		return err
	}
	if err := i.VerifyAudience(aud); err != nil {
		return err
	}
	if err := i.VerifySubject(sub); err != nil {
		return err
	}
	if err := i.VerifyExpiresAt(now); err != nil {
		return err
	}
	if err := i.VerifyIssuedAt(now); err != nil {
		return err
	}
	return nil
}

func (i *IDToken) VerifyAudience(audience string) error {
	// REQUIRED. Audience(s) that this ID Token is intended for. It MUST
	// contain the OAuth 2.0 client_id of the Relying Party as an audience
	// value. It MAY also contain identifiers for other audiences. In the
	// general case, the aud value is an array of case sensitive strings.
	// In the common special case when there is one audience, the aud value
	// MAY be a single case sensitive string.
	claims := i.StandardClaims
	if !claims.VerifyAudience(audience, true) {
		return errors.New("aud does not match")
	}
	return nil
}

func (i *IDToken) VerifyIssuer(issuer string) error {
	// REQUIRED. Issuer Identifier for the Issuer of the response. The iss
	// value is a case sensitive URL using the https scheme that contains
	// scheme, host, and optionally, port number and path components and no
	// query or fragment components.
	claims := i.StandardClaims
	iss := claims.Issuer
	if err := validIss(iss); err != nil {
		return err
	}
	if !claims.VerifyIssuer(issuer, true) {
		return errors.New("iss does not match")
	}
	return nil
}

func (i *IDToken) VerifySubject(subject string) error {
	// REQUIRED. Subject Identifier. A locally unique and never reassigned
	// identifier within the Issuer for the End-User, which is intended to
	// be consumed by the Client, e.g., 24400320 or AItOawmwtWwcT0k51Bayew-
	// NvutrJUqsvl6qs7A4. It MUST NOT exceed 255 ASCII characters in length.
	// The sub value is a case sensitive string.
	sub := i.StandardClaims.Subject
	if err := validSub(sub); err != nil {
		return err
	}
	// if !cmpstr(sub, subject, true) {
	if eq := cmpstr(sub, subject, true); !eq {
		return errors.New("sub does not match")
	}
	return nil
}

func (i *IDToken) VerifyExpiresAt(now int64) error {
	// REQUIRED. Expiration time on or after which the ID Token MUST NOT
	// be accepted for processing. The processing of this parameter
	// requires that the current date/time MUST be before the expiration
	// date/time listed in the value. Implementers MAY provide for some
	// small leeway, usually no more than a few minutes, to account for
	// clock skew. Its value is a JSON number representing the number of
	// seconds from 1970-01-01T0:0:0Z as measured in UTC until the
	// date/time. See RFC 3339 [RFC3339] for details regarding date/times
	// in general and UTC in particular.

	claims := i.StandardClaims
	if ok := claims.VerifyExpiresAt(now, true); !ok {
		return errors.New("token expired")
	}
	return nil
}

func (i *IDToken) VerifyIssuedAt(now int64) error {
	// REQUIRED. Time at which the JWT was issued. Its value is a JSON
	// number representing the number of seconds from 1970-01-01T0:0:0Z as
	// measured in UTC until the date/time.
	claims := i.StandardClaims
	if ok := claims.VerifyIssuedAt(now, true); !ok {
		return errors.New("token used before issued")
	}
	return nil
}

func (i *IDToken) VerifyAuthTime() error {
	// Time when the End-User authentication occurred. Its value is a JSON
	// number representing the number of seconds from 1970-01-01T0:0:0Z as
	// measured in UTC until the date/time. When a max_age request is made
	// or when auth_time is requested as an Essential Claim, then this
	// Claim is REQUIRED; otherwise, its inclusion is OPTIONAL. (The
	// auth_time Claim semantically corresponds to the OpenID 2.0 PAPE
	// [OpenID.PAPE] auth_time response parameter.)
	return nil
}

func (i *IDToken) VerifyNonce(nonce string) error {
	// String value used to associate a Client session with an ID Token,
	// and to mitigate replay attacks. The value is passed through
	// unmodified from the Authentication Request to the ID Token. If
	// present in the ID Token, Clients MUST verify that the nonce Claim
	// Value is equal to the value of the nonce parameter sent in the
	// Authentication Request. If present in the Authentication Request,
	// Authorization Servers MUST include a nonce Claim in the ID Token
	// with the Claim Value being the nonce value sent in the
	// Authentication Request. Authorization Servers SHOULD perform no
	// other processing on nonce values used. The nonce value is a case
	// sensitive string.
	if eq := cmpstr(i.Nonce, nonce, true); !eq {
		return errors.New("nonce does not match")
	}
	return nil
}

func (i *IDToken) VerifyAuthenticationContextClassReference() error {
	// OPTIONAL. Authentication Context Class Reference. String specifying
	// an Authentication Context Class Reference value that identifies the
	// Authentication Context Class that the authentication performed
	// satisfied. The value "0" indicates the End-User authentication did
	// not meet the requirements of ISO/IEC 29115 [ISO29115] level 1.
	// Authentication using a long-lived browser cookie, for instance, is
	// one example where the use of "level 0" is appropriate.
	// Authentications with level 0 SHOULD NOT be used to authorize access
	// to any resource of any monetary value. (This corresponds to the
	// OpenID 2.0 PAPE [OpenID.PAPE] nist_auth_level 0.) An absolute URI or
	// an RFC 6711 [RFC6711] registered name SHOULD be used as the acr
	// value; registered names MUST NOT be used with a different meaning
	// than that which is registered. Parties using this claim will need to
	// agree upon the meanings of the values used, which may be
	// context-specific. The acr value is a case sensitive string.
	return nil
}

func (i *IDToken) VerifyAuthenticationMethodReferences() error {
	// OPTIONAL. Authentication Methods References. JSON array of strings
	// that are identifiers for authentication methods used in the
	// authentication. For instance, values might indicate that both
	// password and OTP authentication methods were used. The definition of
	// particular values to be used in the amr Claim is beyond the scope of
	// this specification. Parties using this claim will need to agree upon
	// the meanings of the values used, which may be context-specific. The
	// amr value is an array of case sensitive strings.
	return nil
}

func (i *IDToken) VerifyAuthorizedParty() error {
	// OPTIONAL. Authorized party - the party to which the ID Token was
	// issued. If present, it MUST contain the OAuth 2.0 Client ID of this
	// party. This Claim is only needed when the ID Token has a single
	// audience value and that audience is different than the authorized
	// party. It MAY be included even when the authorized party is the same
	// as the sole audience. The azp value is a case sensitive string
	// containing a StringOrURI value.
	return nil
}

// HasEmail returns true if the email scope is present.
func (i *IDToken) HasEmail() bool {
	return i.Email != nil
}

// HasAddress returns true if the address scope is present.
func (i *IDToken) HasAddress() bool {
	return i.Address != nil
}

// HasProfile returns true if the profile scope is present.
func (i *IDToken) HasProfile() bool {
	return i.Profile != nil
}

// HasPhone returns true if the phone scope is present.
func (i *IDToken) HasPhone() bool {
	return i.Phone != nil
}

// -- helpers

func validIss(iss string) error {
	if !govalidator.IsURL(iss) {
		return errors.New("issuer must be url")
	}
	if !strings.EqualFold(iss[0:5], "https") {
		return errors.New("issuer scheme must be https")
	}
	return nil
}

func validSub(sub string) error {
	if len(sub) > 255 {
		return errors.New("subject cannot be longer than 255 characters")
	}
	return nil
}
