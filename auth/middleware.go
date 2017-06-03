package auth

import (
  "context"
  "fmt"
  "net/http"
  "time"

  "github.com/julienschmidt/httprouter"
)

func chain(next http.Handler) httprouter.Handle {
  return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ctx := context.WithValue(r.Context(), "params", ps)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

// A middleware to check if the user is already registered
// func checkEmailExist () {

// }

// A middleware to check if the user is logged in
func isLoggedIn(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("access_token")
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }
    fmt.Println(cookie.Value)
    next.ServeHTTP(w, r)
  })
}

func mockCookie(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    expiration := time.Now().Add(30 * 24 * time.Hour)
    cookie := http.Cookie{Name: "access_token", Value: "this_is_your_access_token", Expires: expiration}
    http.SetCookie(w, &cookie)
    next.ServeHTTP(w, r)
  })
}

func cookieHandler(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("hello cookie!"))
}
