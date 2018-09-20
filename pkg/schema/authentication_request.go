package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
)

//go:generate go run gen.go -file $PWD/json/authentication-request.json -var authenticationRequestJSON -out $PWD/gen_authentication_request.go

// AuthenticationRequest represents the struct to validate client metadata.
type AuthenticationRequest struct {
	schema *jsonschema.Schema
}

// NewAuthenticationRequestValidator returns a new pointer to the
// ClientValidator.
func NewAuthenticationRequestValidator() (*AuthenticationRequest, error) {
	schema, err := loadSchema(authenticationRequestJSON)
	return &AuthenticationRequest{schema}, err
}

// Validate validates the given client metadata and returns the corresponding
// errors.
func (a *AuthenticationRequest) Validate(data interface{}) (*Result, error) {
	result, err := validate(a.schema, data)
	return result, err
}
