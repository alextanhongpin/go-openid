package auth

import (
  "log"

  "github.com/alextanhongpin/go-openid/app"
  "github.com/alextanhongpin/go-openid/util"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

// type Service interface {
//   GetUser(string) (User, error)
// }

// type authsvc struct{}

// func (s authsvc) GetUser(id string) {
//   var user User
//   user.ID = id
//   return user, nil
// }

// func NewAuthService() Service {
//   return &authsvc{}
// }

type Service interface {
  GetUser(string) (*User, error)
}

type authservice struct {
  db *app.Database
}

func MakeAuthService(db *app.Database) Service {
  return &authservice{db}
}

func (s authservice) GetUser(id string) (*User, error) {
  // var response User

  // session := s.db.CopySession()
  // defer session.Close()
  // c := s.db.Collection("user", session)

  // err := c.FindId(bson.ObjectIdHex(id)).One(&response)
  // if err != nil {
  //   if err == mgo.ErrNotFound {
  //     var u *User
  //     return u, nil
  //   }
  //   return &response, err
  // }
  // return &response, nil
  var user User
  return &user, nil
}

type service interface {
  checkExist(string) (*User, error)
  fetchOne(string) (*User, error)
  // fetchMany() ([]User, error)
  create(User) (string, error)
  fetchMany() []User
}

type userService struct {
  db *app.Database
}

// create a new user
func (svc userService) create(user User) (string, error) {
  log.Printf("auth/service/userService.create email=%s, password=%s", user.Email, user.Password)
  hashedPassword, err := util.HashPassword(user.Password)
  if err != nil {
    return "", err
  }

  session := svc.db.CopySession()
  defer session.Close()
  c := svc.db.Collection("user", session)

  id := bson.NewObjectId()
  err = c.Insert(&User{
    ID:       id,
    Email:    user.Email,
    Password: hashedPassword,
  })
  if err != nil {
    log.Printf("auth/service/userService.create error=%s", err.Error())
    return "", err
  }
  return id.Hex(), nil
}

func (svc userService) fetchMany() []User {
  var users []User
  session := svc.db.CopySession()
  defer session.Close()
  c := svc.db.Collection("user", session)

  iter := c.Find(bson.M{}).Iter()
  user := User{}
  for iter.Next(&user) {
    users = append(users, user)
  }
  return users
}

func (svc userService) checkExist(email string) (*User, error) {
  var result User
  log.Printf("auth/*.fetchOne email=%s", email)

  session := svc.db.CopySession()
  defer session.Close()
  c := svc.db.Collection("user", session)

  log.Printf("auth/*.fetchOne found collection user %v", c)
  err := c.Find(bson.M{"email": email}).One(&result)
  if err != nil {
    if err == mgo.ErrNotFound {
      log.Printf("No user %v, %v", result, &result)
      var u *User
      return u, nil
    }

    log.Printf("auth/*.fetchOne error=%s", err.Error())
    return &result, err

  }

  return &result, nil
}

func (svc userService) fetchOne(id string) (*User, error) {
  var result User
  log.Printf("auth/*.fetchOne id=%s", id)

  session := svc.db.CopySession()
  defer session.Close()
  c := svc.db.Collection("user", session)

  log.Printf("auth/*.fetchOne found collection user %v", c)
  err := c.FindId(bson.ObjectIdHex(id)).One(&result)
  if err != nil {
    if err == mgo.ErrNotFound {
      log.Printf("No user %v, %v", result, &result)
      var u *User
      return u, nil
    }

    log.Printf("auth/*.fetchOne error=%s", err.Error())
    return &result, err

  }
  log.Printf("Found user %+v, %+v", result, &result)
  return &result, nil
}
