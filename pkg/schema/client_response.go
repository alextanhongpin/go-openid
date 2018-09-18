package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
)

//go:generate go run gen.go -file $PWD/json/client-response.json -var clientResponseJSON -out $PWD/gen_client_response.go

// ClientResponse represents the struct to validate client metadata.
type ClientResponse struct {
	schema *jsonschema.Schema
}

// NewClientResponseValidator returns a new pointer to the
// ClientValidator.
func NewClientResponseValidator() (*ClientResponse, error) {
	schema, err := loadSchema(clientResponseJSON)
	return &ClientResponse{schema}, err
}

// Validate validates the given client metadata and returns the corresponding
// errors.
func (c *ClientResponse) Validate(data interface{}) (*Result, error) {
	result, err := validate(c.schema, data)
	return result, err
}
