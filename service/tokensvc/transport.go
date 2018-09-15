package tokensvc

import (
	"log"

	"github.com/alextanhongpin/go-openid/app"
)

func FeatureToggle(isEnabled bool) func(app.Env) {
	return func(env app.Env) {
		if !isEnabled {
			return
		}
		log.Println("Token Service is enabled")

		r := env.Router
		s := MakeTokenService(env.DB, env.Cache)
		e := MakeServerEndpoints(s)

		r.POST("/token", e.PostToken())
	}
}
