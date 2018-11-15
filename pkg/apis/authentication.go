package apis

import (
	"errors"
	"strings"

	"github.com/alextanhongpin/go-openid/repository"
)

type AuthenticationRequest struct{}

func Authenticate(req *AuthenticationRequest) error {

	// fetch - get data from somewhere
	// save - store data to repo
	// provide - create a new data
	// update - overwrite an existing data
	// compute - do some work to get data
	validations := []func(req *AuthnRequest) error{
		ValidateScope,
		ValidateClient(repo),
	}

	// Advantage of function is you can pass in value diractly. They are also good if your functiosn are pure. Advantage of struct is large initialization.
}

func Consent() {}

func ValidateScope(req *AuthnRequest) error {
	if req.Scope == "" {
		return errors.New("scope required")
	}
	if !strings.Contains(req.Scope, "openid") {
		return errors.New(`scope "openid" is required`)
	}
}
func ValidateClient(repo repository.Client) func(req *AuthnRequest) error {
	return func(req *AuthnRequest) error {

		repo.Has(req.ClientID)
	}
}
