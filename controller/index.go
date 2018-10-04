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
func NewIndex(opts ...indexOption) Index {
	i := Index{}
	for _, o := range opts {
		o(&i)
	}
	return i
}

type indexOption func(i *Index)

// -- setters

// SetTemplate sets the existing template.
func IndexTemplate(t *html5.Template) indexOption {
	return func(i *Index) {
		i.template = t
	}
}

// SetSession sets the existing session.
func IndexSession(s *session.Manager) indexOption {
	return func(i *Index) {
		i.session = s
	}
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
