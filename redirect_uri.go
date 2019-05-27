package openid

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type RedirectURI string

func (r RedirectURI) Validate() error {
	if !govalidator.IsURL(string(r)) {
		return errors.New("redirect_uri is invalid")
	}
	return nil
}

// RedirectURIs represents a slice of valid redirect uris.
type RedirectURIs []RedirectURI

// Contains checks if the redirect uri is present in the slice.
func (r RedirectURIs) Contains(uri string) bool {
	for _, u := range r {
		if u == uri {
			return true
		}
	}
	return false
}
