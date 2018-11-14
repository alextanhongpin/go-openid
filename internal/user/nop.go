package user

import openid "github.com/alextanhongpin/go-openid"

type NopModel struct{}

func (n *NopModel) NewUser(email, password string) (*openid.User, error) {
	return nil, nil
}

func (n *NopModel) ValidateEmail(email string) error {
	return nil
}

func (n *NopModel) ValidatePassword(password string) error {
	return nil
}

func (n *NopModel) ValidateLimit(limit int) error {
	return nil
}
