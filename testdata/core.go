package testdata

import (
	"context"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/service"
	"github.com/stretchr/testify/mock"
)

type coreService struct {
	service.Core
	mock.Mock
}

func NewCoreService() coreService {
	return coreService{}
}

func (c *coreService) PreAuthenticate(req *openid.AuthenticationRequest) error {
	args := c.Called(req)
	return args.Error(0)
}

func (c *coreService) Authenticate(ctx context.Context, req *openid.AuthenticationRequest) (*openid.AuthenticationResponse, error) {
	args := c.Called(ctx, req)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*openid.AuthenticationResponse), args.Error(1)
}

func (c *coreService) Token(ctx context.Context, req *openid.AccessTokenRequest) (*openid.AccessTokenResponse, error) {
	args := c.Called(ctx, req)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*openid.AccessTokenResponse), args.Error(1)
}
