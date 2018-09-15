package authsvc

import (
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/go-playground/validator.v9"
)

// Address represents the address schema in the database
type Address struct {
	Formatted     string `json:"formatted,omitempty" bson:"formatted,omitempty"`
	StreetAddress string `json:"street_address,omitempty" bson:"street_address,omitempty"`
	Locality      string `json:"locality,omitempty" bson:"locality,omitempty"`
	Region        string `json:"region,omitempty" bson:"region,omitempty"`
	PostalCode    string `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	Country       string `json:"country,omitempty" bson:"country,omitempty"`
	Latitude      int64  `json:"latitude,omitempty" bson:"latitude,omitempty" validate:"latitude"`
	Longitude     int64  `json:"longitude,omitempty" bson:"longitude,omitempty" validate:"longitude"`
}

// Email represents the email schema in the database
type Email struct {
	Email         string `json:"email" bson:"email,omitempty" validate:"email,required"`
	EmailVerified bool   `json:"email_verified" bson:"email_verified,omitempty"`
}

// Phone represents the phone schema in the database
type Phone struct {
	PhoneNumber         string `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	PhoneNumberVerified bool   `json:"phone_number_verified,omitempty" bson:"phone_number_verified,omitempty"`
}

// CustomTime implements a custom unmarshalling for the time
type CustomTime struct {
	time.Time
}

// UnmarshalJSON takes a javascript timestamp and convert it to a go time
func (m *CustomTime) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}
	*m = CustomTime{Time: time.Unix(i/1000, (i%1000)*1000*1000)}
	return nil
}

// Profile represents the profile schema in the database
type Profile struct {
	Name              string      `json:"name,omitempty" bson:"name,omitempty"`
	FamilyName        string      `json:"family_name,omitempty" bson:"family_name,omitempty"`
	GivenName         string      `json:"given_name,omitempty" bson:"given_name,omitempty"`
	MiddleName        string      `json:"middle_name,omitempty" bson:"middle_name,omitempty"`
	Nickname          string      `json:"nickname,omitempty" bson:"nickname,omitempty"`
	PreferredUsername string      `json:"preferred_username,omitempty" bson:"preferred_username,omitempty" validate:"required"`
	Profile           string      `json:"profile,omitempty" bson:"profile,omitempty"`
	Picture           string      `json:"picture,omitempty" bson:"picture,omitempty"`
	Website           string      `json:"website,omitempty" bson:"website,omitempty" validate:"url"`
	Gender            string      `json:"gender,omitempty" bson:"gender,omitempty" validate:"eq=male|eq=female"`
	BirthDate         *CustomTime `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
	ZoneInfo          string      `json:"zone_info,omitempty" bson:"zone_info,omitempty"`
	Locale            string      `json:"locale,omitempty" bson:"locale,omitempty"`
	UpdatedAt         *time.Time  `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// User represents the user schema in the database
type User struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email     *Email        `json:"email,omitempty" bson:"email,omitempty" validate:"email,required"`
	Phone     *Phone        `json:"phone,omitempty" bson:"phone,omitempty"`
	Profile   *Profile      `json:"profile,omitempty" bson:"profile,omitempty"`
	Address   *Address      `json:"address,omitempty" bson:"address,omitempty"`
	Password  string        `json:"password,omitempty" bson:"password,omitempty" validate:"min=6,required"`
	Role      []string      `json:"role,omitempty" bson:"role,omitempty" validate:"dive,eq=admin|eq=user,required"`
	Sub       string        `json:"sub,omitempty" bson:"sub,omitempty"`
	CreatedAt *time.Time    `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// NOTE: We use pointers for time.Time to ensure it is not returned in the json field
// with the sane default

// // Claims represents the custom jwt standard claims
// type Claims struct {
// 	UserID string `json:"user_id"`
// 	jwt.StandardClaims
// }

// TODO: Look at the json-api specification
type getUserRequest struct {
	ID string
}

type getUserResponse struct {
	Data User `json:"data"`
}

type getUsersRequest struct{}

type getUsersResponse struct {
	Data []User `json:"data"`
}

type deleteUserRequest struct {
	ID string `json:"id" bson:"_id,omitempty"`
}

type deleteUserResponse struct {
	Ok bool `json:"ok"`
}

// createUserRequest should be the schema for user creation.
type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type createUserResponse struct {
	ID string `json:"id"`
}

type updateUserRequest struct {
	User
}

type updateUserResponse struct {
	Ok bool `json:"ok"`
}

type postRegisterRequest struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required,length(6|50)"`
}

type postRegisterResponse struct {
	Ok          bool   `json:"ok"`
	Error       string `json:"error,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

type getLoginCallbackRequest struct {
	UserID string `json:"user_id,omitempty"`
}

type getLoginCallbackResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

type postLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type postLoginResponse struct {
	Ok          bool   `json:"ok"`
	Error       string `json:"error,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

type getUserViewRequest struct {
	ID string `json:"id"`
}

type getUserViewResponse struct {
	User
}

type getUserEditViewRequest struct {
	ID string `json:"id"`
}

type getUserEditViewResponse struct {
	User
}

type getUsersViewRequest struct {
}

//
type getUsersViewResponse struct {
	Data  []User `json:"data,omitempty"`
	Count int    `json:"count,omitempty"`
}
