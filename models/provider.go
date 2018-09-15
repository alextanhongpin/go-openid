package models

type Provider struct {
	Issuer                                           string   // "https://server.example.com"
	AuthorizationEndpoint                            string   // "https://server.example.com/connect/authorize"
	TokenEndpoint                                    string   // "https://server.example.com/connect/token"
	UserinfoEndpoint                                 string   // "https://server.example.com/connect/userinfo"
	JwksURI                                          string   // "https://server.example.com/jwks.json",
	RegistrationEndpoint                             string   //   "https://server.example.com/connect/register",
	ScopesSupported                                  []string //  ["openid", "profile", "email", "address","phone", "offline_access"],
	ResponseTypesSupported                           []string //  ["code", "code id_token", "id_token", "token id_token"],
	ResponseModesSupported                           []string //  Defaults to  ["query", "fragment"].
	GrantTypesSupported                              []string // Defaults to ["authorization_code", "implicit"]
	AcrValuesSupported                               []string // ["urn:mace:incommon:iap:silver","urn:mace:incommon:iap:bronze"],
	SubjectTypesSupported                            []string // ["public", "pairwise"],
	IDTokenSigningAlgValuesSupported                 []string // ["RS256", "ES256", "HS256"],
	id_token_encryption_alg_values_supported         []string //  ["RSA1_5", "A128KW"],
	id_token_encryption_enc_values_supported         []string // ["A128CBC-HS256", "A128GCM"],
	userinfo_signing_alg_values_supported            []string // ["RS256", "ES256", "HS256"],
	userinfo_encryption_alg_values_supported         []string // ["RSA1_5", "A128KW"],
	userinfo_encryption_enc_values_supported         []string // ["A128CBC-HS256", "A128GCM"],
	request_object_signing_alg_values_supported      []string
	request_object_encryption_alg_values_supported   []string // ["none", "RS256", "ES256"],
	request_object_encryption_enc_values_supported   []string
	token_endpoint_auth_methods_supported            []string // ["client_secret_basic", "private_key_jwt"],
	token_endpoint_auth_signing_alg_values_supported []string //  ["RS256", "ES256"],
	display_values_supported                         []string // ["page", "popup"],
	claim_types_supported                            []string //  ["normal", "distributed"],
	claims_supported                                 []string // ["sub", "iss", "auth_time", "acr", "name", "given_name", "family_name", "nickname", "profile", "picture", "website", "email", "email_verified", "locale", "zoneinfo", "http://example.info/claims/groups"],
	service_documentation                            string   //  "http://server.example.com/connect/service_documentation.html",
	claims_locales_supported                         []string
	ui_locales_supported                             []string // ["en-US", "en-GB", "en-CA", "fr-FR", "fr-CA"]
	claims_parameter_supported                       bool     // true
	request_parameter_supported                      string
	request_uri_parameter_supported                  string
	require_request_uri_registration                 string
	op_policy_uri                                    string
	op_tos_uri                                       string
	check_session_iframe                             string // "https://server.example.com/connect/check_session",
	end_session_endpoint                             string // "https://server.example.com/connect/end_session",
}
