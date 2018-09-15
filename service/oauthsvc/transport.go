package oauthsvc

import (
	"github.com/alextanhongpin/go-openid/app"
)

// FeatureToggle allows you to toggle the feature
func FeatureToggle(isEnabled bool) func(app.Env) {
	return func(env app.Env) {
		// Disable the auth-service
		if !isEnabled {
			return
		}

		r := env.Router
		s := MakeOAuthService(env.DB, env.Cache)
		e := MakeServerEndpoints(s)
		t := env.Tmpl

		// GET  /authorize      display the consent view for authorization (Authorization Code Flow)
		// POST /authorize  	user accept/decline the authorization (Authorization Code Flow)

		// TODO: Protect the path and only allow authenticated users to request to the endpoint
		r.GET("/authorize", e.GetAuthorizeView(t))
		r.POST("/authorize", e.PostAuthorize())
	}
}
