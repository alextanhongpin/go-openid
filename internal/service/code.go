package service

import "github.com/alextanhongpin/go-openid/domain/code"

type Code struct {
	codes code.Repository
}

func (c *Code) Validate(in code.Code) error {
	code, err := c.codes.WithCode(in.Code)
	if err != nil {
		return err
	}
	if err := c.codes.Delete(in.Code); err != nil {
		return err
	}
	return code.Validate()
}
