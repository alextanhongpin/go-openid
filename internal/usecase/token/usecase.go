package token

import (
	"context"
	"errors"
	"time"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/domain/client"
	"github.com/alextanhongpin/go-openid/pkg/signer"
	"github.com/alextanhongpin/go-openid/usecase"
	jwt "github.com/dgrijalva/jwt-go"
)

type UseCase struct {
	clientService client.Service
}

// func TokenEndpoint(w http.ResponseWriter, r *http.Request) {
//         ctx := r.Context()
//         clientID, clientSecret := credentialsFromAuthHeader(r.Header.Get("Authorization"))
//         ctx = context.WithValue(ctx, ContextKeyClientID, clientId)
//         ctx = context.WithValue(ctx, ContextKeyClientSecret, clientSecret)
//
//         // TODO: Get the subject from the current session.
//         // session := sessionManager.GetSession(r)
//         // ctx = context.WithValue(ctx, ContextKeyUserID, session.ID)
//
//         var request Client
//         if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
//                 // Handle error here.
//                 // Wrap the error in a custom error format here.
//                 // json.Encoder(w).Encode(Error{Message: err.Error()})
//                 return
//         }
//         response, err := service.Token(ctx, request)
//         if err != nil {
//                 // Handle error here.
//                 return
//         }
//         // Handle response here.
// }

/*
type Service struct {

}
func (s *Service) Token (ctx context.Context, req *AuthenticateRequest) {
	// Be explicit on what to expect.
	opts := TokenOptions {
		AccessTokenDuration: 2 * time.Hour,
		RefreshTokenDuration: 24 * time.Hour,
		Issuer: "openid",
		TokenType: "Bearer",
	}
	return Token(ctx, opts, s.clientRepo, s.codeRepo, s.signer, req)
}
*/

func (u *UseCase) Token(ctx context.Context, req usecase.TokenRequest) (*usecase.TokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// Authenticate the client credentials.
	{
		clientID, clientSecret := func() (string, string) {
			ctx := client.Context{}
			return ctx.ClientID(), ctx.ClientSecret()
		}()
		// ClientID and ClientSecret must be valid.
		if err := u.clientService.ValidateCredentials(clientID, clientSecret); err != nil {
			return nil, err
		}
	}
	if !openid.AuthorizationCode.Equal(req.GrantType) {
		return nil, errors.New("grant_type is invalid")
	}
	// Code must be valid.
	if err := u.codeService.Validate(req.Code); err != nil {
		return nil, err
	}
	// UserID must exist in session.
	userID, ok := user.Context{}.Value()
	if !ok {
		return nil, errors.New("unauthorized")
	}
	// This will always return the server time, unless mocked.
	ts, _ := timestamp.Context{}.Value()
	// if !ok {
	//         ts = time.Now()
	// }
	// if ts.After(time.Now()) {
	//         ts = time.Now()
	// }

	claims := Claims{
		defaults: jwt.StandardClaims{
			Subject:  userID,
			IssuedAt: now.Unix(),
			Issuer:   opts.Issuer,
		},
		signer: signer,
	}
	tokenDuration := func(expiresIn time.Duration) ClaimsModifier {
		return func(claims *jwt.StandardClaims) {
			claims.ExpiresAt = now.Add(expiresIn).Unix()
		}
	}
	accessToken, err := claims.Sign(tokenDuration(opts.AccessTokenDuration))
	if err != nil {
		return nil, err
	}
	refreshToken, err := claims.Sign(tokenDuration(opts.RefreshTokenDuration))
	if err != nil {
		return nil, err
	}
	idToken, err := signer.Sign(NewIDToken())
	if err != nil {
		return nil, err
	}
	// Return the response.
	return &TokenResponse{
		ExpiresIn:    int64(opts.AccessTokenDuration.Seconds()),
		TokenType:    opts.TokenType,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}, nil
}

type Claims struct {
	defaults jwt.StandardClaims
	signer   Signer
}

type ClaimsModifier func(claims *jwt.StandardClaims)

func (c *Claims) Sign(modifiers ...ClaimsModifier) (string, error) {
	claims := c.defaults
	for _, modifier := range modifiers {
		modifier(&claims)
	}
	return c.signer.Sign(claims)
}
