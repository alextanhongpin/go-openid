package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
)

// Result represents an alias to the jsonschema.Result type.
type Result = jsonschema.Result

// Validator represents the interface for validation.
type Validator interface {
	Validate(data interface{}) (*Result, error)
}
