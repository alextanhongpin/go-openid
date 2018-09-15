package app

import (
	"html/template"
	"net/http"
	"strings"
)

// Template contains all the page template.
type Template struct {
	templates map[string]*template.Template
}

// Render will perform a lookup in existing templates and render it
func (tmp Template) Render(w http.ResponseWriter, name string, data interface{}) {
	t, ok := tmp.templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
	}
	t.Execute(w, data)
}

// NewTemplate returns a new template
func NewTemplate() *Template {
	templates := make(map[string]*template.Template)

	// GET /					 renders the index page
	// GET /login                renders the login page
	// GET /register             renders the register page
	// GET /authorize            renders the consent page
	// GET /users      		     renders the user list page
	// GET /users/:id            renders the user profile page
	// GET /users/:id/edit       renders the edit user page
	// GET /clients/register     renders the client registration page
	// GET /clients 		     renders the client list page
	// GET /clients/:id			 renders the client page

	funcsMap := template.FuncMap{"StringsJoin": strings.Join}

	templates["index"] = template.Must(template.ParseFiles("templates/index.html"))
	templates["login"] = template.Must(template.ParseFiles("templates/login.html"))
	templates["register"] = template.Must(template.ParseFiles("templates/register.v2.html"))
	templates["consent"] = template.Must(template.ParseFiles("templates/consent.html"))
	templates["users"] = template.Must(template.ParseFiles("templates/users.html"))
	templates["user"] = template.Must(template.ParseFiles("templates/user.v2.html"))
	templates["user-edit"] = template.Must(template.ParseFiles("templates/user-edit.html"))
	templates["client_register"] = template.Must(template.ParseFiles("templates/client_register.html"))
	templates["clients"] = template.Must(template.ParseFiles("templates/clients.html"))
	templates["client"] = template.Must(template.New("client.html").Funcs(funcsMap).ParseFiles("templates/client.html"))

	return &Template{templates: templates}
}
