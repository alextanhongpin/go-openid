package tokensvc

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/models"
	"github.com/alextanhongpin/go-openid/utils/crypto"
	"github.com/alextanhongpin/go-openid/utils/jwt"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type Service interface {
	// GetToken(getTokenRequest) (getTokenResponse, error)
	PostToken(postTokenRequest) (*postTokenResponse, error)
	CheckClient(clientRequest) (*clientResponse, error)
	CheckCode(codeRequest) (*codeResponse, error)
}

type tokenservice struct {
	db    *app.Database
	cache *app.Cache
}

func MakeTokenService(db *app.Database, cache *app.Cache) *tokenservice {
	return &tokenservice{
		db:    db,
		cache: cache,
	}
}

func (s tokenservice) PostToken(req postTokenRequest) (*postTokenResponse, error) {
	//   HTTP/1.1 200 OK
	//   Content-Type: application/json
	//   Cache-Control: no-cache, no-store
	//   Pragma: no-cache

	refreshToken, err := crypto.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	tokenClaims := jwt.Claims{
		"123456",
		jwtgo.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "www.openid",
		},
	}
	idToken, err := tokenClaims.Sign()
	if err != nil {
		return nil, err
	}

	accessTokenClaims := jwt.Claims{
		"123456",
		jwtgo.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "www.openid",
		},
	}
	accessToken, err := accessTokenClaims.Sign()
	if err != nil {
		return nil, err
	}

	return &postTokenResponse{
		AccessToken:  accessToken, // "SlAV32hkKG",
		TokenType:    "Bearer",
		ExpiresIn:    time.Now().Add(1 * time.Hour),
		RefreshToken: refreshToken, //"tGzv3JOkF0XG5Qx2TlKWIA",
		IDToken:      idToken,      // "eyJ0 ... NiJ9.eyJ1c ... I6IjIifX0.DeWt4Qu ... ZXso",
	}, nil
}

func (s tokenservice) CheckClient(req clientRequest) (*clientResponse, error) {
	session := s.db.NewSession()
	defer session.Close()
	var client models.Client

	c := s.db.Collection("client", session)
	err := c.Find(bson.M{"client_id": req.ID}).One(&client)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &clientResponse{
		Data: client,
	}, nil
}

func (s tokenservice) CheckCode(req codeRequest) (*codeResponse, error) {
	code, err := s.cache.Client.Get(req.ClientID).Result()
	if err != nil {
		return nil, err
	}

	// One-time usage, delete it after calling
	n, err := s.cache.Client.Del(req.ClientID).Result()
	if err != nil {
		return nil, err
	}

	log.Printf("tokenservice method=CheckCode code=%s\n deleted=%d", code, n)
	return &codeResponse{
		Exist: code == req.Code,
	}, nil
}
