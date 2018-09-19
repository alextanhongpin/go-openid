package service

type User interface {
	Register(email, password string) error
	UserInfo()
	Authenticate()
	Authorize()
}
