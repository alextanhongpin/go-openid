package app

import (
  "github.com/hashicorp/go-memdb"
  "github.com/julienschmidt/httprouter"
)

// Env is the app env
type Env struct {
  Db     *memdb.MemDB
  Router *httprouter.Router
  Tmpl   *Template
  // Log    zap.Logger
  // log
}
