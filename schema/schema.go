package schema

import (
	"errors"
	"fmt"

	jsonschema "github.com/xeipuuv/gojsonschema"
)

type SchemaKey string

type Schema struct {
	loader map[SchemaKey]jsonschema.JSONLoader
	schema map[SchemaKey]*jsonschema.Schema
}

func New() (*Schema, error) {
	var (
		cm  = SchemaKey("client-metadata")
		crr = SchemaKey("client-registration-response")
	)

	loader := make(map[SchemaKey]jsonschema.JSONLoader)
	loader[cm] = jsonschema.NewStringLoader(clientMetadata)
	loader[crr] = jsonschema.NewStringLoader(clientRegistrationResponse)

	schema := make(map[SchemaKey]*jsonschema.Schema)

	s := Schema{
		loader: loader,
		schema: schema,
	}
	if err := s.loadOne(cm, clientMetadata); err != nil {
		return nil, err
	}
	if err := s.loadDependent(cm, crr); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Schema) Validate(key string, data interface{}) (*jsonschema.Result, error) {
	doc, ok := s.schema[SchemaKey(key)]
	if !ok {
		return nil, fmt.Errorf("json-schema doc missing for %s", key)
	}
	return doc.Validate(jsonschema.NewGoLoader(data))
}

func (s *Schema) loadOne(key SchemaKey, data string) (err error) {
	s.loader[key] = jsonschema.NewStringLoader(data)
	s.schema[key], err = jsonschema.NewSchema(s.loader[key])
	return
}

func (s *Schema) loadDependent(keys ...SchemaKey) error {
	if len(keys) < 2 {
		return errors.New("need at least 2 keys")
	}
	sl := jsonschema.NewSchemaLoader()
	idx := len(keys) - 1

	// Proceed when there are only more than 1 key.
	var loaders []jsonschema.JSONLoader
	for _, k := range keys[0:idx] {
		js, ok := s.loader[k]
		if !ok {
			return fmt.Errorf("schema missing for %s", k)
		}
		loaders = append(loaders, js)
	}
	// The json schema must be added in the correct order.
	if err := sl.AddSchemas(loaders...); err != nil {
		return err
	}

	var err error
	// Bind all the previous schemas to the "root" schema.
	s.schema[keys[idx]], err = sl.Compile(s.loader[keys[idx]])
	return err
}
