package schema

import (
	"errors"

	jsonschema "github.com/xeipuuv/gojsonschema"
)

func loadSchema(source string) (*jsonschema.Schema, error) {
	loader := jsonschema.NewStringLoader(source)
	return jsonschema.NewSchema(loader)
}

func validate(schema *jsonschema.Schema, data interface{}) (*Result, error) {
	result, err := schema.Validate(jsonschema.NewGoLoader(data))
	if err != nil {
		return nil, err
	}
	if !result.Valid() {
		err := result.Errors()[0]
		return result, errors.New(err.Description())
	}
	return result, nil
}
