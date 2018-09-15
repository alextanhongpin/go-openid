package clientsvc

import "github.com/alextanhongpin/go-openid/app"

func FeatureToggle(isEnabled bool) func(app.Env) {
	return func(env app.Env) {
		if !isEnabled {
			return
		}

		r := env.Router
		e := MakeServerEndpoints(MakeClientService(env.DB))
		t := env.Tmpl

		// GET  	/connect/register		View the client registration form
		// POST 	/connect/register		Create a new client
		// PATCH	/connect/register/:id	Update an existing client
		// DELETE	/connect/register/:id	Delete an existing client
		// GET 	   	/clients/:id			View a client by id
		// GET		/clients				View a list of clients
		// GET		/api/clients/:id		Get a client by id
		// GET		/api/clients			Get a list of clients

		// TODO: Protect the endpoints
		r.GET("/connect/register", e.GetClientConnectView(t))
		r.POST("/connect/register", e.PostClient())
		// TODO: Should be connect/register?client_id=client_id
		r.PATCH("/connect/register/:id", e.UpdateClient())
		// TODO: Should be connect/register?client_id=client_id
		r.DELETE("/connect/register/:id", e.DeleteClient())

		// Should be /connect/register?client_id=123456
		r.GET("/clients/:id", e.GetClientView(t))
		r.GET("/clients", e.GetClientsView(t))

		// API endpoints has to be protected
		r.GET("/api/clients/:id", e.GetClient())
		r.GET("/api/clients", e.GetClients())
	}
}
