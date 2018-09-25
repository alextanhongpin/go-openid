package main

import (
	"net/http"
	"net/url"

	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/julienschmidt/httprouter"
)

// Endpoints represent the endpoints for the OpenIDConnect.
type Endpoints struct {
	service        *serviceImpl
	sessionManager *session.Manager
	template       *html5.Template
}

// NewEndpoints returns a pointer to new endpoints.
func NewEndpoints(s *serviceImpl) *Endpoints {
	return &Endpoints{
		service: s,
	}
}

func (e *Endpoints) GetLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type data struct {
		ReturnURL string
	}
	if ok := e.sessionManager.HasSession(r); ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	parseURI := func(u url.Values) (string, error) {
		base64uri := u.Get("return_url")
		if base64uri == "" {
			return "/", nil
		}
		return decodeBase64(base64uri)
	}

	uri, err := parseURI(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d := data{uri}
	e.template.Render(w, "login", d)
}
