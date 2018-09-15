package schema

var clientRegistrationResponse = `
{
	"$id": "client-registration-response.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "Client Metadata according to the OIDC Client Registration Specification",
	"type": "object",
	"allOf": [
		{
			"$ref": "#/definitions/client-response"
		},
		{
			"$ref": "client-metadata.json"
		}
	],
	"definitions": {
		"client-response": {
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
			],
			"additionalProperties": false
		}
	}
}`