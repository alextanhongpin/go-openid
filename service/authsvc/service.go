package authsvc

import (
	"log"
	"strings"
	"time"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/utils/encrypt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Service implements all the methods.
type Service interface {
	GetUser(getUserRequest) (*getUserResponse, error)
	GetUsers(getUsersRequest) (*getUsersResponse, error)
	DeleteUser(deleteUserRequest) (*deleteUserResponse, error)
	CreateUser(createUserRequest) (*createUserResponse, error)
	UpdateUser(updateUserRequest) (*updateUserResponse, error)
	CheckUser(string) (*User, error)
}

type authservice struct {
	db *app.Database
}

// MakeAuthService creates a new authentication service.
func MakeAuthService(db *app.Database) Service {
	return &authservice{db}
}

// GetUser return a user by id.
func (s authservice) GetUser(req getUserRequest) (*getUserResponse, error) {
	session := s.db.NewSession()
	defer session.Close()

	c := s.db.Collection("user", session)

	var user User
	err := c.FindId(bson.ObjectIdHex(req.ID)).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &getUserResponse{}, nil
		}
		return nil, nil
	}

	return &getUserResponse{
		Data: user,
	}, nil
}

// GetUsers returns a list of users from the database.
func (s authservice) GetUsers(req getUsersRequest) (*getUsersResponse, error) {
	session := s.db.NewSession()
	defer session.Close()
	c := s.db.Collection("user", session)

	var users []User

	err := c.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	log.Printf("GetUsers type=service event=find_all users_sum=%v users=%v", len(users), users)
	return &getUsersResponse{
		Data: users,
	}, nil
}

func (s authservice) UpdateUser(req updateUserRequest) (*updateUserResponse, error) {
	session := s.db.NewSession()
	defer session.Close()

	c := s.db.Collection("user", session)

	log.Printf("UpdateUser type=service request=%#v", req.Phone)

	err := c.UpdateId(req.ID, bson.M{
		"$set": bson.M{
			"phone": bson.M{
				"phone_number": req.Phone.PhoneNumber,
			},
			"profile": bson.M{
				"name":               req.Profile.Name,
				"family_name":        req.Profile.FamilyName,
				"given_name":         req.Profile.GivenName,
				"middle_name":        req.Profile.MiddleName,
				"nickname":           req.Profile.Nickname,
				"preferred_username": req.Profile.PreferredUsername,
				"profile":            req.Profile.Profile,
				"picture":            req.Profile.Picture,
				"website":            req.Profile.Website,
				"gender":             req.Profile.Gender,
				"birth_date":         req.Profile.BirthDate,
				"zone_info":          req.Profile.ZoneInfo,
				"locale":             req.Profile.Locale,
				"updated_at":         req.Profile.UpdatedAt,
			},
			"address": bson.M{
				"formatted":      req.Address.Formatted,
				"street_address": req.Address.StreetAddress,
				"locality":       req.Address.Locality,
				"region":         req.Address.Region,
				"postal_code":    req.Address.PostalCode,
				"country":        req.Address.Country,
				"latitude":       req.Address.Latitude,
				"longitude":      req.Address.Longitude,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &updateUserResponse{Ok: true}, nil
}

// CreateUser create a new user in the database.
func (s authservice) CreateUser(user createUserRequest) (*createUserResponse, error) {

	session := s.db.NewSession()
	defer session.Close()
	c := s.db.Collection("user", session)

	hashedPassword, err := encrypt.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	id := bson.NewObjectId()
	createdAt := time.Now()

	err = c.Insert(&User{
		ID: id,
		Email: &Email{
			Email: user.Email,
		},
		Profile: &Profile{
			Nickname: strings.Split(user.Email, "@")[0],
		},
		Password:  hashedPassword,
		CreatedAt: &createdAt,
	})

	if err != nil {
		return nil, err
	}

	return &createUserResponse{
		ID: id.Hex(),
	}, nil
}

func (s authservice) DeleteUser(user deleteUserRequest) (*deleteUserResponse, error) {
	log.Printf("DeleteUser type=service message=start params=%v", user.ID)
	session := s.db.NewSession()
	defer session.Close()

	c := s.db.Collection("user", session)
	// bson.ObjectId(user.ID)
	err := c.RemoveId(bson.ObjectIdHex(user.ID))

	// err := c.Remove(bson.M{"_id": user.ID})
	log.Printf("DeleteUser type=service err=%v \n", err)
	if err != nil {
		return nil, err
	}

	return &deleteUserResponse{
		Ok: true,
	}, nil
}

// CheckUser validates if the email already exists in the database.
func (s authservice) CheckUser(email string) (*User, error) {
	session := s.db.NewSession()
	defer session.Close()
	c := s.db.Collection("user", session)

	var user User
	err := c.Find(bson.M{"email.email": email}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
