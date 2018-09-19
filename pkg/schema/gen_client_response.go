// Code generated by github.com/alextanhongpin/go-openid/schema/gen; DO NOT EDIT.
// This file was generated by robots at
// 2018-09-18T23:31:08+08:00

package schema 

var clientResponseJSON = `{
	"$id": "client-response.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "Client Metadata according to the OIDC Client Registration Specification",
	"type": "object",
	"allOf": [{
		"$ref": "#/definitions/client-response"
	}],
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
					"type": "integer",
					"default": 0
				}
			},
			"required": [
				"client_id"
			]
		}
	}
}
`