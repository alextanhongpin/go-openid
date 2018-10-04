


package repository

import "github.com/alextanhongpin/go-openid"

// Code represents the operations for the code repository.
type Code interface {
	Get(id string) (*openid.Code, bool)
	Put(id string, code *openid.Code)
	Delete(id string)
}
