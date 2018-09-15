package openidsvc


import 	"github.com/alextanhongpin/go-openid/app"

func FeatureToggle (isEnabled bool) app.Env {
	return func (env app.Env) {
		r := env.Router

		r.GET("/.well-known/openid-configuration")
		r.GET("/userinfo")

// 		  GET /userinfo HTTP/1.1
//   Host: server.example.com
//   Authorization: Bearer SlAV32hkKG

	}
}