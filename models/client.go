package models

import "gopkg.in/mgo.v2/bson"

// Email address string
type Email string

// ClientMetadata represents the data that is stored in the database
type ClientMetadata struct {
	RedirectURIs                 []string `json:"redirect_uris,omitempty" bson:"redirect_uris"`                                     // ["https://client.example.org/callback", "https://client.example.org/callback2"]
	ResponseTypes                []string `json:"response_types,omitempty" bson:"response_types"`                                   // defaults to `code`
	GrantTypes                   []string `json:"grant_types,omitempty" bson:"grant_types"`                                         // authorization_code, implicit, refresh_token
	ApplicationType              string   `json:"application_type,omitempty" bson:"application_type"`                               // web|native
	Contacts                     []Email  `json:"contacts,omitempty"`                                                               // ["ve7jtb@example.org", "mary@example.org"]
	ClientName                   string   `json:"client_name,omitempty" bson:"client_name"`                                         // "My Example"
	LogoURI                      string   `json:"logo_uri,omitempty" bson:"logo_uri"`                                               // "https://client.example.org/logo.png"
	ClientURI                    string   `json:"client_uri,omitempty" bson:"client_uri"`                                           //
	PolicyURI                    string   `json:"policy_uri,omitempty" bson:"policy_uri"`                                           //
	TosURI                       string   `json:"tos_uri,omitempty" bson:"tos_uri"`                                                 //
	JwksURI                      string   `json:"jwks_uri,omitempty" bson:"jwks_uri"`                                               // "https://client.example.org/my_public_keys.jwks"
	Jwks                         string   `json:"jwks,omitempty"`                                                                   //
	SectorIdentifierURI          string   `json:"sector_identifier_uri,omitempty" bson:"sector_identifier_uri"`                     // "https://other.example.net/file_of_redirect_uris.json"
	SubjectType                  string   `json:"subject_type,omitempty" bson:"subject_type"`                                       // "pairwise"
	IDTokenSignedResponseAlg     string   `json:"id_token_signed_response_alg,omitempty" bson:"id_token_signed_response_alg"`       //
	IDTokenEncryptedResponseAlg  string   `json:"id_token_encrypted_response_alg,omitempty" bson:"id_token_encrypted_response_alg"` //
	IDTokenEncryptedResponseEnc  string   `json:"id_token_encrypted_response_enc,omitempty" bson:"id_token_encrypted_response_enc"` //
	UserinfoSignedResponseAlg    string   `json:"userinfo_signed_response_alg,omitempty" bson:"userinfo_signed_response_alg"`       //
	UserinfoEncryptedResponseAlg string   `json:"userinfo_encrypted_response_alg,omitempty" bson:"userinfo_encrypted_response_alg"` // "RSA1_5"
	UserinfoEncryptedResponseEnc string   `json:"userinfo_encrypted_response_enc,omitempty" bson:"userinfo_encrypted_response_enc"` // "A128CBC-HS256"
	RequestObjectSigningAlg      string   `json:"request_object_signing_alg,omitempty" bson:"request_object_signing_alg"`           //
	RequestObjectEncryptionAlg   string   `json:"request_object_encryption_alg,omitempty" bson:"request_object_encryption_alg"`     //
	RequestObjectEncryptionEnc   string   `json:"request_object_encryption_enc,omitempty" bson:"request_object_encryption_enc"`     //
	TokenEndpointAuthMethod      string   `json:"token_endpoint_auth_method,omitempty" bson:"token_endpoint_auth_method"`           // client_secret_post|client_secret_basic|client_secret_jwt|private_key_jwt|none
	TokenEndpointAuthSigningAlg  string   `json:"token_endpoint_auth_signing_alg,omitempty" bson:"token_endpoint_auth_signing_alg"` //
	DefaultMaxAge                int      `json:"default_max_age,omitempty" bson:"default_max_age"`                                 //
	RequireAuthTime              int      `json:"require_auth_time,omitempty" bson:"require_auth_time"`                             //
	DefaultAcrValues             []string `json:"default_acr_values,omitempty" bson:"default_acr_values"`                           //
	InitiateLoginURI             string   `json:"initiate_login_uri,omitempty" bson:"initiate_login_uri"`                           //
	RequestURIs                  []string `json:"request_uris,omitempty" bson:"request_uris"`                                       // ["https://client.example.org/rf.txt#qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"]
}

// Client represents the client schema in the database
type Client struct {
	// Custom key-value pairs
	ID                      bson.ObjectId `json:"id" bson:"_id,omitempty"`                                              //
	IsPublished             bool          `json:"is_production" bson:"is_published"`                                    // Whether it fulfils the requirement to be dispatched
	Version                 string        `json:"version"`                                                              // The version of the api it supports
	UserID                  string        `json:"user_id,omitempty" bson:"user_id"`                                     //
	CreatedAt               *int64        `json:"created_at,omitempty" bson:"created_at"`                               //
	UpdatedAt               *int64        `json:"modified_at,omitempty" bson:"updated_at"`                              //
	DeletedAt               *int64        `json:"deleted_at,omitempty" bson:"deleted_at"`                               // Everything beneath this is the spec
	ClientID                string        `json:"client_id,omitempty" bson:"client_id"`                                 // Mandatory
	ClientSecret            string        `json:"client_secret,omitempty" bson:"client_secret"`                         //
	RegistrationAccessToken string        `json:"registration_access_token,omitempty" bson:"registration_access_token"` //
	RegistrationClientURI   string        `json:"registration_client_uri,omitempty" bson:"registration_client_uri"`     //
	ClientIDIssuedAt        *int64        `json:"client_id_issued_at,omitempty" bson:"client_id_issued_at"`             //
	ClientSecretExpiresAt   *int64        `json:"client_secret_expires_at,omitempty" bson:"client_secret_expires_at"`   //
	ClientMetadata          `bson:"client_metadata"`
}
