package schema

import (
	"fmt"

	jsonschema "github.com/xeipuuv/gojsonschema"
)

// Result represents an alias to the jsonschema.Result type.
type Result = jsonschema.Result

// Validator represents the interface for validation.
type Validator interface {
	Validate(data interface{}) (*Result, error)
}

type Validators interface {
	Validate(schema string, data interface{}, required bool) (*Result, error)
	Register(schema string, validator Validator)
}

type validators struct {
	schemas map[string]Validator
}

func NewValidators() *validators {
	return &validators{
		schemas: make(map[string]Validator),
	}
}

func (v *validators) Validate(schema string, data interface{}, required bool) (*Result, error) {
	if !required {
		return nil, nil
	}
	s, ok := v.schemas[schema]
	if !ok {
		if required {
			return nil, fmt.Errorf("schema %s is missing", schema)
		}
		return nil, nil
	}
	return s.Validate(data)
}

func (v *validators) Register(schema string, validator Validator) {
	v.schemas[schema] = validator
}
