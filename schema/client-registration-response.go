package schema

import jsonschema "github.com/xeipuuv/gojsonschema"

//go:generate go run gen.go -file json/client-registration-response.json -var clientRegistrationResponseJSON -out gen_client-registration-response.go

type ClientRegistrationResponse struct {
	schema *jsonschema.Schema
}

func NewClientRegistrationResponseValidator() (*ClientRegistrationResponse, error) {
	var (
		cl  = jsonschema.NewStringLoader(clientJSON)
		crr = jsonschema.NewStringLoader(clientRegistrationResponseJSON)
		sl  = jsonschema.NewSchemaLoader()
	)
	if err := sl.AddSchemas(cl); err != nil {
		return nil, err
	}
	schema, err := sl.Compile(crr)
	return &ClientRegistrationResponse{schema}, err
}

func (c *ClientRegistrationResponse) Validate(data interface{}) (*Result, error) {
	return validate(c.schema, data)
}
