package openid

import "github.com/alextanhongpin/passwd"

// User represents the user struct.
type User struct {
	hashedPassword string
	ID             string `json:"id,omitempty"`
	Address
	Email
	Phone
	Profile
}

func (u *User) SetPassword(password string) error {
	hash, err := passwd.Hash(password)
	if err != nil {
		return err
	}
	u.hashedPassword = hash
	return nil
}

func (u *User) ComparePassword(password string) error {
	return passwd.Verify(password, u.hashedPassword)
}

func (u *User) Clone() *User {
	copy := new(User)
	*copy = *u
	return copy
}

func (u *User) ToIDToken() *IDToken {
	user := u.Clone()

	idToken := NewIDToken()
	*idToken.Address = user.Address
	*idToken.Email = user.Email
	*idToken.Phone = user.Phone
	*idToken.Profile = user.Profile

	return idToken
}
