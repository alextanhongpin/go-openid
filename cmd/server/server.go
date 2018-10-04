package main

import (
	"flag"
	"log"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid/controller"
	"github.com/alextanhongpin/go-openid/internal/client"
	"github.com/alextanhongpin/go-openid/pkg/appsensor"
	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/session"
)

func main() {
	var (
		port   = flag.Int("port", 8080, "the port of the application")
		tplDir = flag.String("tpldir", "templates", "the datadir of the html templates")
	)
	flag.Parse()

	// Create new router.
	r := httprouter.New()

	// Load templates.
	tpl := html5.New(*tplDir)
	tpl.Load("login", "register", "client-register", "consent", "index")

	sessMgr := session.NewManager()
	sessMgr.Start()
	defer sessMgr.Stop()

	// TODO: Run a cron job that handles deletion of unused data for a
	// certain period of time.
	aps := appsensor.NewLoginDetector()

	// -- endpoints
	{
		c := controller.NewIndex(
			controller.IndexTemplate(tpl),
			controller.IndexSession(sessMgr),
		)
		r.GET("/", c.GetIndex)
	}
	{
		c := controller.NewUser(
			controller.UserSession(sessMgr),
			controller.UserAppSensor(aps),
			controller.UserTemplate(tpl),
		)
		r.POST("/logout", c.PostLogout)
		r.GET("/register", c.GetRegister)
		r.GET("/login", c.GetLogin)
		r.POST("/login", c.PostLogin)
		r.POST("/register", c.PostRegister)
	}
	{
		s, err := client.NewService()
		if err != nil {
			log.Fatal(err)
		}
		c := controller.NewClient(
			controller.ClientTemplate(tpl),
			controller.ClientService(s),
		)
		r.GET("/connect/register", c.GetClientRegister)
		r.POST("/connect/register", c.PostClientRegister)
	}
	{
		c := controller.NewCore(
			controller.CoreSession(sessMgr),
			controller.CoreTemplate(tpl),
		)
		r.GET("/authorize", c.GetAuthorize)
		r.POST("/authorize", c.PostAuthorize)
		r.POST("/token", c.PostToken)
	}
	srv := gsrv.New(*port, r)
	<-srv
	log.Println("Gracefully shutdown HTTP server.")
}
