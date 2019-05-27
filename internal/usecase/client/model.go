package client

import (
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/rs/xid"
)

func NewClient() openid.Client {
	return openid.Client{
		ClientID:                xid.New().String(),
		ClientSecret:            ranstr.RandomString(32),
		RegistrationAccessToken: "",
		ClientIDIssuedAt:        time.Now().UTC().Unix(),
		ClientSecretExpiresAt:   0,
		RegistrationClientURI:   "https://server.example.com/c2id/clients",
	}
}

// func (m *Model) GenerateRegistrationAccessToken(clientID string) (string, error) {
//         var (
//                 aud = "https://server.example.com/c2id/clients"
//                 sub = clientID
//                 iss = clientID
//
//                 iat = time.Now().UTC()
//                 day = time.Hour * 24
//                 exp = iat.Add(7 * day)
//                 key = []byte("client_token_secret")
//         )
//         claims := crypto.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
//         accessToken, err := crypto.NewJWT(key, claims)
//         return accessToken, err
// }
