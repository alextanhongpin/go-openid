package main

import jwt "github.com/dgrijalva/jwt-go"

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
	Address
	Email
	Phone
	Profile
}

func NewIDToken() *IDToken {
	return &IDToken{}
}

// Address represents the fields for the address scope.
type Address struct {
	Country       string `json:"country,omitempty"`
	Formatted     string `json:"formatted,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
}

// Email represents the fields for the email scope.
type Email struct {
	Email         string `json:"email,omitempty"` // Preferred e-mail address.
	EmailVerified bool   `json:"email_verified"`  // True if the e-mail address has been verified; otherwise false.
}

// Phone represents the fields for the scope phone.
type Phone struct {
	PhoneNumber         string `json:"phone_number,omitempty"` // Preferred telephone number.
	PhoneNumberVerified bool   `json:"phone_number_verified"`  // True if the phone number has been verified; otherwise false.
}

// Profile represents the fields for the scope profile.
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
