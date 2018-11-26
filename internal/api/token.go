package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type TokenOptions struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
	TokenType            string
}

type TokenRequest struct {
	GrantType   string
	Code        string
	RedirectURI string
}

func (req *TokenRequest) IsGrantTypeAuthorizationCode() bool {
	return req.GrantType == "authorization_code"
}

func (req *TokenRequest) ValidateRequiredFields() error {
	// Validate required fields.
	if stringIsEmpty(req.Code) {
		return errors.New("code is required")
	}
	if stringIsEmpty(req.GrantType) {
		return errors.New("grant_type is required")
	}
	if stringIsEmpty(req.RedirectURI) {
		return errors.New("redirect_uri is required")
	}

	// Another option is to create the URI type with a validate method.
	if err := URI(req.RedirectURI).Validate(); err != nil {
		return fmt.Errorf(`"%s" is not a valid redirect_uri`, req.RedirectURI)
	}
	return nil
}

func (req *TokenRequest) ValidateClient(repo ClientRepository, clientID, clientSecret string) error {
	if clientID == "" {
		return errors.New("client_id is required")
	}
	if clientSecret == "" {
		return errors.New("client_secret is required")
	}
	client, err := repo.GetClientByCredentials(clientID, clientSecret)
	if err != nil {
		return err
	}
	if !client.HasRedirectURI(req.RedirectURI) {
		return errors.New("redirect_uri is invalid")
	}
	return nil
}

func (req *TokenRequest) ValidateCode(repo CodeRepository) error {
	code, err := repo.GetCodeByID(req.Code)
	if err != nil {
		return err
	}
	// TODO: Remove code if it exists.
	if err := repo.Delete(req.Code); err != nil {
		return err
	}
	if !code.IsValid() {
		return errors.New("code is invalid")
	}
	return nil
}

type TokenResponse struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	ExpiresIn    int64
	IDToken      string
}

func credentialsFromAuthHeader(str string) (string, string) {
	return "", ""
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

func Token(
	ctx context.Context,
	opts TokenOptions,
	clientRepo ClientRepository,
	codeRepo CodeRepository,
	signer Signer,
	req *TokenRequest,
) (*TokenResponse, error) {
	validate := func() error {
		// Validate required fields.
		if err := req.ValidateRequiredFields(); err != nil {
			return err
		}

		// Validate speficic business rules.
		if !req.IsGrantTypeAuthorizationCode() {
			return fmt.Errorf(`grant_type "%s" is invalid`, req.GrantType)
		}

		clientID, _ := ctx.Value(ContextKeyClientID).(string)
		clientSecret, _ := ctx.Value(ContextKeyClientSecret).(string)

		// Note the inversion of control. Rather than clientRepo.Has(clientID, clientSecret),
		// we make the request as the owner, not the Client.
		if err := req.ValidateClient(clientRepo, clientID, clientSecret); err != nil {
			return err
		}
		if err := req.ValidateCode(codeRepo); err != nil {
			return err
		}
		return nil
	}

	// Validate the requests.
	if err := validate(); err != nil {
		return nil, err
	}

	// Build the response.
	// Subject is the user we want to issue the token to.
	userID, _ := ctx.Value(ContextKeyUserID).(string)
	if stringIsEmpty(userID) {
		return nil, errors.New("subject is required")
	}

	// Allows us to provide the timestamp through the context to be mocked.
	now, ok := ctx.Value(ContextKeyTimestamp).(time.Time)
	if !ok {
		now = time.Now().UTC()
	}

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
