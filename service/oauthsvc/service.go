package oauthsvc

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/models"
	"github.com/alextanhongpin/go-openid/utils/crypto"
)

// Service is the interface for the OAuth Service
type Service interface {
	GetAuthorize(getAuthorizeRequest) (*getAuthorizeResponse, error)
	PostAuthorize(postAuthorizeRequest) (*postAuthorizeResponse, error)
	// External service
	GetClient(getClientRequest) (*getClientResponse, error)
}

type oauthsvc struct {
	db    *app.Database
	cache *app.Cache
}

// exported func MakeOAuthService returns unexported type *oauthsvc.oauthsvc, which can be annoying to use
// MakeOAuthService func(db *app.Database, cache *app.Cache) *oauthsvc

// MakeOAuthService is a helper to initialize the OAuth Service with the database through Dependency Injection
func MakeOAuthService(db *app.Database, cache *app.Cache) *oauthsvc {
	return &oauthsvc{db, cache}
}

func (s oauthsvc) GetAuthorize(req getAuthorizeRequest) (*getAuthorizeResponse, error) {
	// session := s.db.NewSession()
	// defer session.Close()

	// // authorize collection
	// c := s.db.Collection("oauth", session)
	return &getAuthorizeResponse{}, nil
}

func (s oauthsvc) PostAuthorize(req postAuthorizeRequest) (*postAuthorizeResponse, error) {
	token, err := crypto.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	// If req.ClientID, remove the token
	// Else,
	log.Printf("PostAuthorize type=service crypto=%v\n client_id=%v", token, req.ClientID)

	// val, err := client.Get(req.ClientID).Result()
	// if err == redis.Nil {
	// 	// Key does not exist
	// } else if err != nil {
	// 	return nil, err
	// } else {
	// 	// key exists, remove it
	// 	n, err := client.Del(req.ClientID).Result()
	// }
	// Cache the client id to the token
	err = s.cache.Client.Set(req.ClientID, token, time.Second*300).Err()
	if err != nil {
		return nil, err
	}
	// Store the token in redis, which will expire in 5 minutes time,
	// or once it has been exchanged with access token
	return &postAuthorizeResponse{
		State:       req.State,
		Code:        token,
		RedirectURI: req.RedirectURI,
	}, nil
}

func (s oauthsvc) GetClient(req getClientRequest) (*getClientResponse, error) {
	session := s.db.NewSession()
	defer session.Close()

	// client collection
	c := s.db.Collection("client", session)

	var client models.Client
	err := c.Find(bson.M{"client_id": req.ClientID}).One(&client)
	// What if the client does not exist/has been deleted?
	if err != nil {
		return nil, err
	}

	return &getClientResponse{Data: client}, nil
}
