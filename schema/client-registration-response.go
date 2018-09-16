package schema

import jsonschema "github.com/xeipuuv/gojsonschema"

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

var clientRegistrationResponseJSON = `
{
	"$id": "http://server.example.com/schemas/client-registration-response.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "Client Metadata according to the OIDC Client Registration Specification",
	"type": "object",
	"allOf": [
		{
			"$ref": "#/definitions/client"
		},
		{
			"$ref": "http://server.example.com/schemas/client.json"
		}
	],
	"definitions": {
		"client": {
			"type": "object",
			"properties": {
				"client_id": {
					"type": "string"
				},
				"client_secret": {
					"type": "string"
				},
				"registration_access_token": {
					"type": "string"
				},
				"registration_client_uri": {
					"type": "string",
					"format": "uri"
				},
				"client_id_issued_at": {
					"type": "integer"
				},
				"client_secret_expires_at": {
					"type": "integer"
				}
			},
			"required": [
				"client_id",
				"client_secret_expires_at"
			]
		}
	}
}`
