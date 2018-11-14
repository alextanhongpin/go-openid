package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
)

func ValidateString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func LowerLimit(i, lower int) error {
	if i < lower {
		return fmt.Errorf(`"%d" must be greater than %d`, i, lower)
	}
	return nil
}

func UpperLimit(i, upper int) error {
	if i > upper {
		return fmt.Errorf(`"%d" must be less than %d`, i, upper)
	}
	return nil
}

func ValidateRange(i, lower, upper int) error {
	err := LowerLimit(i, lower)
	if err != nil {
		return err
	}
	return UpperLimit(i, upper)
}

func ValidatePassword(password string) error {
	if ValidateString(password) {
		return errors.New("password is required")
	}
	if len(password) < 6 {
		return errors.New("password length must be larger than or equal to 6")
	}
	return nil
}

func ValidateEmail(email string) error {
	if ValidateString(email) {
		return errors.New("email is required")
	}
	if !govalidator.IsEmail(email) {
		return fmt.Errorf(`"%s" is not a valid email`, email)
	}
	return nil
}
