package openid

import "testing"

func TestIDToken(t *testing.T) {
	token := IDToken{
		Iss:      "https://server.example.com",
		Sub:      "24400320",
		Aud:      "s6BhdRkqt3",
		Nonce:    "n-0S6_WzA2Mj",
		Exp:      1311281970,
		Iat:      1311280970,
		AuthTime: 1311280969,
		Acr:      "urn:mace:incommon:iap:silver",
	}
}
