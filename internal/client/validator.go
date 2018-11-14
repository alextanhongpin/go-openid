package client

import (
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

type Validator struct {
	client         *schema.Client
	clientResponse *schema.ClientResponse
}

func NewValidator() (*Validator, error) {
	client, err := schema.NewClientValidator()
	if err != nil {
		return nil, err
	}

	clientResponse, err := schema.NewClientResponseValidator()
	if err != nil {
		return nil, err
	}

	return &Validator{
		client:         client,
		clientResponse: clientResponse,
	}, nil
}
