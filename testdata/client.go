package testdata

import (
	openid "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/service"
	"github.com/stretchr/testify/mock"
)

type clientService struct {
	service.Client
	mock.Mock
}

func NewClientService() *clientService {
	return &clientService{}
}

func (c *clientService) Register(req *openid.Client) (*openid.Client, error) {
	args := c.Called(req)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*openid.Client), args.Error(1)
}

func (c *clientService) Read(clientID string) (*openid.Client, error) {
	args := c.Called(clientID)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*openid.Client), args.Error(1)
}
