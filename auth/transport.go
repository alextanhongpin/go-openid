package auth

import (
  "context"
  "fmt"
  "net/http"

  "github.com/alextanhongpin/go-openid/app"
  "github.com/julienschmidt/httprouter"
)

// FeatureToggle allows you to toggle the feature
func FeatureToggle(isEnabled bool) func(app.Env) {
  return func(env app.Env) {
    // Don't run this
    if !isEnabled {
      return
    }
    e := endpoint{userService{db: env.Db}}
    r := env.Router
    t := env.Tmpl

    // API
    r.GET("/api/users/:id", e.getUserHandler())

    // Static & Forms
    r.GET("/users/:id", e.viewUserHandler(t))

    r.POST("/login", e.createUserHandler())
    r.GET("/login", e.loginViewHandler(t))

    r.GET("/register", e.registerViewHandler(t))
    r.POST("/register", e.registerHandler())
    r.GET("/context", middleware(handler))
  }
}

type key string

const ctxName key = "id"

func middleware(next httprouter.Handle) httprouter.Handle {
  return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ctx := context.WithValue(r.Context(), ctxName, 12345)
    next(w, r.WithContext(ctx), ps)
  })
}
func handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  reqID := r.Context().Value(ctxName).(int)
  fmt.Fprintf(w, "hello request id: %d", reqID)
}
