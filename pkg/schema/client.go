package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
)

//go:generate go run gen.go -file $PWD/json/client.json -var clientJSON -out $PWD/gen_client.go

// Client represents the struct to validate client metadata.
type Client struct {
	schema *jsonschema.Schema
}

// NewClientValidator returns a new pointer to the
// ClientValidator.
func NewClientValidator() (*Client, error) {
	schema, err := loadSchema(clientJSON)
	return &Client{schema}, err
}

// Validate validates the given client metadata and returns the corresponding
// errors.
func (c *Client) Validate(data interface{}) (*Result, error) {
	result, err := validate(c.schema, data)
	return result, err
}
