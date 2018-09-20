package authorization

import (
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

func Authorize(req *oidc.AuthenticationRequest) error {
	if err := req.Validate("", nil); err != nil {
		return err
	}
	if prompt := req.GetPrompt(); prompt.Is(oidc.PromptNone) {
		// OR: LoginRequired
		return oidc.ErrInteractionRequired
	}
	return nil
}

func Authenticate(user *oidc.User) (*oidc.AuthenticationResponse, error) {
	// Get id_token of user.
	idToken := user.ToIDToken()

	// Generate a new id_token. TODO: Store the token in envvars.
	token, err := idToken.SignHS256([]byte("id_token_secret"))
	if err != nil {
		return nil, err
	}

	var (
		aud = "https://server.example.com"
		sub = user.ID
		iss = user.ID
		iat = time.Now().UTC()
		exp = iat.Add(2 * time.Hour)
		key = []byte("access_token_secret")
	)

	// Generate an access token with the user_id.
	claims := crypto.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
	accessToken, err := crypto.NewJWT(key, claims)
	if err != nil {
		return nil, err
	}

	return &oidc.AuthenticationResponse{
		AccessToken: accessToken,
		ExpiresIn:   exp.Unix(),
		IDToken:     token,
		State:       "",
		TokenType:   "bearer",
	}, nil
}