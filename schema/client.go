package schema

import (
	"fmt"

	jsonschema "github.com/xeipuuv/gojsonschema"
)

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
		return result, fmt.Errorf("%s: %s", err.Field(), err.Description())
	}
	return result, nil
}

// Validate validates the given client metadata and returns the corresponding
// errors.
func (c *ClientValidator) Validate(data interface{}) (*Result, error) {
	return validate(c.schema, data)
}

var clientJSON = `
{
	"$id": "http://server.example.com/schemas/client.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "Client Metadata according to the OIDC Client Registration Specification",
	"type": "object",
	"properties": {
		"redirect_uris": {
			"type": "array",
			"items": {
				"type": "string",
				"format": "uri"
			}
		},
		"response_types": {
			"type": "array",
			"items": {
				"type": "string",
				"enum": [
					"code",
					"token",
					"id_token"
				],
				"default": "code"
			}
		},
		"grant_types": {
			"type": "array",
			"items": {
				"type": "string",
				"enum": [
					"authorization_code",
					"implicit",
					"refresh_token"
				],
				"default": "authorization_code"
			}
		},
		"application_type": {
			"type": "string",
			"enum": [
				"web",
				"native"
			],
			"default": "web"
		},
		"contacts": {
			"type": "array",
			"items": {
				"type": "string",
				"format": "email"
			}
		},
		"client_name": {
			"type": "string"
		},
		"logo_uri": {
			"type": "string",
			"format": "uri"
		},
		"client_uri": {
			"type": "string",
			"format": "uri"
		},
		"policy_uri": {
			"type": "string",
			"format": "uri"
		},
		"tos_uri": {
			"type": "string",
			"format": "uri"
		},
		"jwks_uri": {
			"type": "string",
			"format": "uri"
		},
		"jwks": {
			"type": "string"
		},
		"sector_identifier_uri": {
			"type": "string",
			"format": "uri"
		},
		"subject_type": {
			"type": "string",
			"enum": [
				"pairwise",
				"public"
			]
		},
		"id_token_signed_response_alg": {
			"type": "string"
		},
		"id_token_encrypted_response_alg": {
			"type": "string"
		},
		"id_token_encrypted_response_enc": {
			"type": "string"
		},
		"userinfo_signed_response_alg": {
			"type": "string"
		},
		"userinfo_encrypted_response_alg": {
			"type": "string"
		},
		"userinfo_encrypted_response_enc": {
			"type": "string",
			"default": "A128CBC-HS256"
		},
		"request_object_signing_alg": {
			"type": "string"
		},
		"request_object_encryption_alg": {
			"type": "string"
		},
		"request_object_encryption_enc": {
			"type": "string",
			"default": "A128CBC-HS256"
		},
		"token_endpoint_auth_method": {
			"type": "string",
			"enum": [
				"client_secret_post",
				"client_secret_basic",
				"client_secret_jwt",
				"private_key_jwt",
				"none"
			]
		},
		"token_endpoint_auth_signing_alg": {
			"type": "string"
		},
		"default_max_age": {
			"type": "integer"
		},
		"require_auth_time": {
			"type": "boolean"
		},
		"default_acr_values": {
			"type": "string"
		},
		"initiate_login_uri": {
			"type": "string"
		},
		"request_uris": {
			"type": "array",
			"items": {
				"type": "string",
				"format": "uri"
			}
		}
	},
	"required": [
		"redirect_uris"
	],
	"additionalProperties": false
}`
