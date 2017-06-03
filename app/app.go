package app

import (
  "github.com/julienschmidt/httprouter"
)

// Env is the app env
type Env struct {
  DB     *Database
  Router *httprouter.Router
  Tmpl   *Template
  // Log    zap.Logger
  // log
}
