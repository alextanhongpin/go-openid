package models

// IDToken is a security token that contains Claims about the authentication of an End-User by an Authorization Server when using a Client, and potentially other requested Claims. The ID Token is represented as a JSON Web Token (JWT) [JWT].
type IDToken struct {
	ISS      string   `json:"iss" validate:"required"` // REQUIRED. "https://server.example.com"
	SUB      string   `json:"sub" validate:"required"` // REQUIRED. 24400320
	AUD      string   `json:"udd" validate:"required"` // REQUIRED. s6BhdRkqt3
	EXP      int      `json:"exp" validate:"required"` // REQUIRED. 1311281970
	IAT      int      `json:"iat" validate:"required"` // REQUIRED. 1311280970
	AuthTime int      `json:"auth_time"`               // OPTIONAL.
	Nonce    string   `json:"nonce"`                   // OPTIONAL.
	ACR      string   `json:"acr"`                     // OPTIONAL.
	AMR      []string `json:"amr"`                     // OPTIONAL.
	AZP      string   `json:"azp"`                     // OPTIONAL.
}
