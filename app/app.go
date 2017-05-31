package app

import (
  "html/template"

  "github.com/hashicorp/go-memdb"
  "github.com/julienschmidt/httprouter"
)

// Env is the app env
type Env struct {
  Db     *memdb.MemDB
  Router *httprouter.Router
  Tmpl   map[string]*template.Template
  // log
}
