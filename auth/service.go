package auth

import (
  "time"

  "github.com/hashicorp/go-memdb"

  "github.com/alextanhongpin/go-openid/util"
)

type service interface {
  fetchOne(email string) (*User, error)
  // fetchMany() ([]User, error)
  create(email, password string) error
}

type userService struct {
  db *memdb.MemDB
}

func (svc userService) create(email, password string) error {
  txn := svc.db.Txn(true)
  defer txn.Abort()

  hashedPassword, err := util.HashPassword(password)
  if err != nil {
    return nil
  }

  u := &User{
    Email:     email,
    Name:      "john",
    Password:  hashedPassword,
    CreatedAt: time.Now(),
  }

  if err := txn.Insert("user", u); err != nil {
    return err
  }

  txn.Commit()

  return nil
}
func (svc userService) fetchOne(email string) (*User, error) {
  var user *User
  txn := svc.db.Txn(false)
  defer txn.Abort()

  // Lookup by id
  // raw, err := txn.First("user", "id", "john.doe@mail.com")

  raw, err := txn.First("user", "id", email)
  if err != nil {
    return user, err
  }

  // No user found
  if raw == nil {
    return user, nil
  }
  return raw.(*User), nil
  // var users []string
  // // for raw := result.Next(); raw != nil; {
  // //  users = append(users, raw.(*User).Name)
  // // }
  // for i := 0; i < 10; i++ {
  //   raw := result.Next()
  //   if raw == nil {
  //     break
  //   }
  //   users = append(users, raw.(*User).Name)
  // }

}
