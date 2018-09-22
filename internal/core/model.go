package core

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/repository"
)

type modelImpl struct {
	code   repository.Code
	client repository.Client
}

func (m *modelImpl) ValidateClient(clientID, redirectURI string) error {
	client, exist := m.client.Get(clientID)
	if !exist || client == nil {
		return errors.New("client does not exist")
	}
	if client.ClientID != clientID {
		return errors.New("client_id does not match")
	}
	if ok := client.RedirectURIs.Contains(redirectURI); !ok {
		return errors.New("redirect_uri does not match")
	}
	return nil
}

func (m *modelImpl) NewCode() string {
	c := crypto.NewXID()
	code := oidc.NewCode(c)
	m.code.Put(c, code)
	return c
}
