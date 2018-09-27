package controller

import (
	"net/http"

	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/julienschmidt/httprouter"
)

// Index represents the index controller.
type Index struct {
	template *html5.Template
	session  *session.Manager
}

// NewIndex returns a new index.
func NewIndex() Index {
	return Index{}
}

// -- setters

// SetTemplate sets the existing template.
func (i *Index) SetTemplate(t *html5.Template) {
	i.template = t
}

// SetSession sets the existing session.
func (i *Index) SetSession(s *session.Manager) {
	i.session = s
}

// GetIndex represents the index endpoint.
func (i *Index) GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type data struct {
		IsLoggedIn bool
	}
	var res data
	sess, err := i.session.GetSession(r)
	if err != nil {
		res.IsLoggedIn = false
	}
	if sess != nil {
		res.IsLoggedIn = true
	}
	i.template.Render(w, "index", res)
}
