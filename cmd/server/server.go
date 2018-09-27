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

type Controllers struct {
	User   controller.User
	Core   controller.Core
	Index  controller.Index
	Client controller.Client
}

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

	ctl := Controllers{
		User:   makeUserController(aps, tpl, sessMgr),
		Index:  makeIndexController(tpl, sessMgr),
		Client: makeClient(tpl),
		Core:   makeCoreController(tpl, sessMgr),
	}

	// -- endpoints
	// Index endpoints.
	r.GET("/", ctl.Index.GetIndex)

	// User endpoints.
	r.POST("/logout", ctl.User.PostLogout)
	r.GET("/register", ctl.User.GetRegister)
	r.GET("/login", ctl.User.GetLogin)
	r.POST("/login", ctl.User.PostLogin)
	r.POST("/register", ctl.User.PostRegister)

	// Client endpoints.
	r.GET("/connect/register", ctl.Client.GetClientRegister)
	r.POST("/connect/register", ctl.Client.PostClientRegister)

	// OpenID Connect endpoints.
	r.GET("/authorize", ctl.Core.GetAuthorize)
	r.POST("/authorize", ctl.Core.PostAuthorize)
	r.POST("/token", ctl.Core.PostToken)

	srv := gsrv.New(*port, r)
	<-srv
	log.Println("Gracefully shutdown HTTP server.")
}

func makeUserController(a appsensor.LoginDetector, h *html5.Template, s *session.Manager) controller.User {
	ctl := controller.NewUser()
	ctl.SetAppSensor(a)
	ctl.SetTemplate(h)
	ctl.SetSession(s)
	return ctl
}

func makeIndexController(h *html5.Template, s *session.Manager) controller.Index {
	c := controller.NewIndex()
	c.SetTemplate(h)
	c.SetSession(s)
	return c
}

func makeClient(h *html5.Template) controller.Client {
	s, err := client.NewService()
	if err != nil {
		log.Fatal(err)
	}
	c := controller.NewClient()
	c.SetTemplate(h)
	c.SetService(s)
	return c
}

func makeCoreController(h *html5.Template, s *session.Manager) controller.Core {
	c := controller.NewCore()
	c.SetTemplate(h)
	c.SetSession(s)
	return c
}
