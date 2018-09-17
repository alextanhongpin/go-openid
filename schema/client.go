package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
)

//go:generate go run gen.go -file json/client.json -var clientJSON -out gen_client.go

// ClientValidator represents the struct to validate client metadata.
type ClientValidator struct {
	schema *jsonschema.Schema
}

// NewClientValidator returns a new pointer to the
// ClientValidator.
func NewClientValidator() (*ClientValidator, error) {
	schema, err := loadSchema(clientJSON)
	return &ClientValidator{schema}, err
}

// Validate validates the given client metadata and returns the corresponding
// errors.
func (c *ClientValidator) Validate(data interface{}) (*Result, error) {
	return validate(c.schema, data)
}
