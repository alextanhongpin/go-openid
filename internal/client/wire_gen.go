// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package client

import (
	database "github.com/alextanhongpin/go-openid/internal/database"
	schema "github.com/alextanhongpin/go-openid/pkg/schema"
)

// Injectors from wire.go:

func NewService() (*clientServiceImpl, error) {
	client, err := schema.NewClientValidator()
	if err != nil {
		return nil, err
	}
	clientResponse, err := schema.NewClientResponseValidator()
	if err != nil {
		return nil, err
	}
	clientKV := database.NewClientKV()
	clientClientModelImpl := NewClientModelImpl(clientKV)
	clientClientValidatorImpl := &clientValidatorImpl{
		client:         client,
		clientResponse: clientResponse,
		model:          clientClientModelImpl,
	}
	clientClientServiceImpl := NewClientServiceImpl(clientClientValidatorImpl)
	return clientClientServiceImpl, nil
}
