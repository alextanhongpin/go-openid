package auth

import (
  "github.com/alextanhongpin/go-openid/app"
)

// FeatureToggle allows you to toggle the feature
func FeatureToggle(isEnabled bool) func(app.Env) {
  return func(env app.Env) {
    e := endpoint{userService{db: env.Db}}
    r := env.Router
    t := env.Tmpl

    // API
    r.GET("/api/users/:id", e.getUserHandler())

    // Static & Forms
    r.GET("/users/:id", e.viewUserHandler(t))
    r.POST("/login", e.createUserHandler())
    r.GET("/login", e.loginHandler(t))
  }
}
