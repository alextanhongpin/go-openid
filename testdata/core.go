package testdata

import (
	"context"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/model"
	"github.com/stretchr/testify/mock"
)

type coreModel struct {
	model.Core
	mock.Mock
}

func NewCoreModel() coreModel {
	return coreModel{}
}

func (c *coreModel) ValidateAuthnRequest(req *openid.AuthenticationRequest) error {
	args := c.Called(req)
	return args.Error(0)
}

func (c *coreModel) ValidateAuthnUser(ctx context.Context, req *openid.AuthenticationRequest) error {
	args := c.Called(ctx, req)
	return args.Error(0)
}

func (c *coreModel) ValidateAuthnClient(req *openid.AuthenticationRequest) error {
	args := c.Called(req)
	return args.Error(0)
}

func (c *coreModel) NewCode() string {
	return "new_code"
}
