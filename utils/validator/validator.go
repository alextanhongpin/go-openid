package validator

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

var (
	errIDRequired = errors.New("Id field is missing")
	errIDInvalid  = errors.New("Id provided is invalid")
)

// ValidateID returns a valid bsonObject id or an error if
// the string is not a valid mongoid
func ValidateID(id string) (bson.ObjectId, error) {
	// Cannot return nil
	oid := bson.NewObjectId()
	// String cannot be empty
	if id == "" {
		return oid, errIDRequired
	}
	// Not a valid object id
	if !bson.IsObjectIdHex(id) {
		return oid, errIDInvalid
	}
	// Valid, parse id
	oid = bson.ObjectIdHex(id)
	return oid, nil
}
