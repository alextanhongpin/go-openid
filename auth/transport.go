package auth

import (
  "html/template"

  "github.com/hashicorp/go-memdb"
  "github.com/julienschmidt/httprouter"
)

// SetupAuth is the route
func SetupAuth(db *memdb.MemDB, r *httprouter.Router, tmpl map[string]*template.Template) {
  e := endpoint{userService{db: db}}
  // API
  r.GET("/api/users/:id", e.getUserHandler())

  // Static & Forms
  r.GET("/users/:id", e.viewUserHandler(tmpl))
  r.POST("/login", e.createUserHandler())
  r.GET("/login", e.loginHandler(tmpl))
}
