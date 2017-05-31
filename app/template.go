package app

import (
  "html/template"
)

// Template is the cached
func Template() map[string]*template.Template {
  templates := make(map[string]*template.Template)
  templates["login"] = template.Must(template.ParseFiles("static/login.html"))
  templates["user"] = template.Must(template.ParseFiles("static/user.html"))
  return templates
}
