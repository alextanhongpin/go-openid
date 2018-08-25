package openid

import "testing"

type accessTokenRequestService struct {
}

func (s *accessTokenRequestService) Do(req AccessTokenRequest) (*AccessTokenResponse, error) {
	return nil, nil
}
func TestAccessTokenRequest(t *testing.T) {}
