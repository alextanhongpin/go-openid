package user

import (
	"errors"
	"strings"
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/asaskevich/govalidator"
)

type Model struct {
}

func NewModel() *Model {
	return &Model{}
}

// NewUser stores the username and hashed password into the storage.
func (m *Model) NewUser(email, password string) (*openid.User, error) {
	user := new(openid.User)
	user.ID = crypto.NewXID()
	user.Profile = openid.Profile{}
	user.Profile.UpdatedAt = time.Now().UTC().Unix()
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	user.Email = openid.Email{
		Email: email,
	}
	return user, nil
}

func (m *Model) ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if ok := govalidator.IsEmail(email); !ok {
		return errors.New("invalid email")
	}
	return nil
}

func (m *Model) ValidatePassword(password string) error {
	if len(password) == 0 {
		return errors.New("arguments cannot be empty")
	}
	if len(password) < 8 {
		return errors.New("password cannot be less than 8 characters")
	}
	return nil
}

func (m *Model) ValidateLimit(limit int) error {
	if limit < 0 {
		return errors.New("limit cannot be less then 0")
	}
	if limit > 100 {
		return errors.New("limit cannot be more than 100")
	}
	return nil
}

func IsEmptyString(str string) bool {
	// Handle empty space.
	return len(strings.TrimSpace(str)) == 0
}
