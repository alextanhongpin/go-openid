package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Home struct {
}

func NewHome() Home {
	return Home{}
}

// GetLogin renders the login html.
func (h *Home) GetLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// h.template.Render(w, "login", nil)
}
