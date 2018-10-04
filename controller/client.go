package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/service"

	"github.com/julienschmidt/httprouter"
)

type Client struct {
	service  service.Client
	template *html5.Template
}

func NewClient(opts ...clientOption) Client {
	c := Client{}
	for _, o := range opts {
		o(&c)
	}
	return c
}

func (c *Client) GetClientRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// sess, err := c.session.GetSession(r)
	// if err != nil {
	//         writeError(w, http.StatusBadRequest, err)
	//         return
	// }

	// TODO: Check if the user is authorized to read the client
	// details.
	id := r.URL.Query().Get("client_id")
	if id == "" {
		c.template.Render(w, "client-register", nil)
		return
	}
	client, err := c.service.Read(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	json.NewEncoder(w).Encode(client)
}

func (c *Client) PostClientRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	// TODO: Check if the user is authorized to perform client
	// registration.

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")

	buildClient := func(r *http.Request) *openid.Client {
		var (
			clientName   = r.FormValue("client_name")
			redirectURIs = strings.Split(r.FormValue("redirect_uris"), " ")
		)
		client := openid.NewClient()
		client.ClientName = clientName
		client.RedirectURIs = redirectURIs
		return client
	}

	client := buildClient(r)

	newClient, err := c.service.Register(client)
	if err != nil {
		v, ok := err.(*openid.ErrorJSON)
		if ok {
			json.NewEncoder(w).Encode(v)
		} else {
			json.NewEncoder(w).Encode(M{"error": err.Error()})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newClient)
}

// -- options

type clientOption func(c *Client)

func ClientService(s service.Client) clientOption {
	return func(c *Client) {
		c.service = s
	}
}

func ClientTemplate(h *html5.Template) clientOption {
	return func(c *Client) {
		c.template = h
	}
}
