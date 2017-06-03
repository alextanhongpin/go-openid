package app

import (
  "html/template"
  "net/http"
)

// Template is the cached
// func Template() map[string]*template.Template {
//   templates := make(map[string]*template.Template)
//   templates["login"] = template.Must(template.ParseFiles("static/login.html"))
//   templates["user"] = template.Must(template.ParseFiles("static/user.html"))
//   return templates
// }
type Template struct {
  templates map[string]*template.Template
}

// Render the template by writing the found template to standard output
func (template Template) Render(w http.ResponseWriter, name string, data interface{}) {
  t, ok := template.templates[name]
  if !ok {
    http.Error(w, "Template not found", http.StatusInternalServerError)
  }
  t.Execute(w, data)
}

// NewTemplate returns a new template
func NewTemplate() *Template {
  templates := make(map[string]*template.Template)

  // Register templates here
  // GET /login
  templates["login"] = template.Must(template.ParseFiles("static/login.html"))

  // GET /register
  templates["register"] = template.Must(template.ParseFiles("static/register.html"))

  // GET /user
  templates["user"] = template.Must(template.ParseFiles("static/user.html"))
  return &Template{templates: templates}
}
