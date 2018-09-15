package authsvc

import (
	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/middlewares"
)

// FeatureToggle allows you to toggle the feature
func FeatureToggle(isEnabled bool) func(app.Env) {
	return func(env app.Env) {
		// Disable the auth-service
		if !isEnabled {
			return
		}

		r := env.Router
		s := MakeAuthService(env.DB)
		e := MakeServerEndpoints(s)
		t := env.Tmpl
		o := env.Tracer

		// POST /api/users      get users
		// GET  /api/users/:id  get a user by id
		// GET  /register       renders the register view
		// POST /register       handle the register form submission

		r.GET("/api/users/:id", middlewares.ProtectAPI(e.GetUser()))
		r.GET("/api/users", middlewares.ProtectAPI(e.GetUsers()))
		r.DELETE("/api/users/:id", middlewares.ProtectAPI(e.DeleteUser()))
		r.PUT("/api/users/:id", middlewares.ProtectAPI(e.UpdateUser()))

		r.GET("/register", middlewares.ProtectRoute(e.GetRegister(t)))
		r.POST("/register", e.PostRegister(o))

		r.GET("/login", middlewares.ProtectRoute(e.GetLogin(t)))
		r.GET("/login/callback", e.GetLoginCallback())
		r.POST("/login", e.PostLogin())

		r.POST("/logout", e.PostLogout())

		r.GET("/users/:id", middlewares.ValidateAuth(e.GetUserView(t)))
		r.GET("/users", middlewares.ValidateAuth(e.GetUsersView(t)))
		r.GET("/users/:id/edit", middlewares.ValidateAuth(e.GetUserEditView(t)))
	}
}
