package auth

import (
	"time"

	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/go-playground/validator.v9"
)

// Address is the schema of the address in the database
type Address struct {
	Formatted     string `json:"formatted"`
	StreetAddress string `json:"street_address"`
	Locality      string `json:"locality"`
	Region        string `json:"region"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	Latitude      int64  `json:"latitude" validate:"latitude"`
	Longitude     int64  `json:"longitude" validate:"longitude"`
}

// User is the schema of the user in the database
type User struct {
	ID                  bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email               string        `json:"email" validate="email"`
	EmailVerified       bool          `json:"email_verified"`
	Password            string        `json:"password" validate="min=6,required"`
	Role                []string      `json:"role" validate="dive,eq=admin|eq=user,required"`
	Sub                 string        `json:"sub"`
	Name                string        `json:"name"`
	GivenName           string        `json:"given_name"`
	FamilyName          string        `json:"family_name"`
	MiddleName          string        `json:"middle_name"`
	Nickname            string        `json:"nickname"`
	PreferredUsername   string        `json:"preferred_username"`
	Profile             string        `json:"profile"`
	Picture             string        `json:"picture"`
	Website             string        `json:"website" validate="url"`
	Gender              string        `json:"gender" validate:eq=male|eq=female`
	BirthDate           time.Time     `json:"birth_date"`
	ZoneInfo            string        `json:"zone_info"`
	Locale              string        `json:"locale"`
	PhoneNumber         string        `json:"phone_number"`
	PhoneNumberVerified bool          `json:"phone_number_verified"`
	Address             Address       `json:"address"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}
