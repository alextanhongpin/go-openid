package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/cmd/server/controller"
	"github.com/alextanhongpin/go-openid/internal/client"
	"github.com/alextanhongpin/go-openid/pkg/appsensor"
	"github.com/alextanhongpin/go-openid/pkg/gsrv"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
	"github.com/alextanhongpin/go-openid/pkg/session"
)

// M represents simple map interface.
type M map[string]interface{}

// Credentials represent the user credentials for the application.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	svc := NewService()

	// Setup user controller.
	userController := controller.NewUser()
	userController.SetAppSensor(aps)
	userController.SetTemplate(tpl)
	userController.SetSession(sessMgr)

	// Setup index controller.
	indexController := controller.NewIndex()
	indexController.SetTemplate(tpl)
	indexController.SetSession(sessMgr)

	// Setup client controller.

	clientService, err := client.NewService()
	if err != nil {
		log.Fatal(err)
	}
	clientController := controller.NewClient()
	clientController.SetTemplate(tpl)
	clientController.SetService(clientService)

	// -- endpoints

	getAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()

		var req oidc.AuthenticationRequest
		if err := querystring.Decode(q, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := svc.core.PreAuthenticate(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		redirectToLogin := func() {
			redirectURI := getHost(r)
			redirectURI.RawQuery = q.Encode()
			base64uri := encodeBase64(redirectURI.String())
			u := fmt.Sprintf(`http://localhost:8080/login?return_url=%s`, base64uri)
			http.Redirect(w, r, u, http.StatusFound)
		}

		isAuthorized := sessMgr.HasSession(r)
		prompt := req.GetPrompt()

		// If the prompt is set to none, but the user is unauthorized,
		// an error should be returned indicating that login is
		// required.
		if prompt.Is(oidc.PromptNone) && !isAuthorized {
			http.Error(w, oidc.ErrLoginRequired.Error(), http.StatusBadRequest)
			return
		}

		// If the user is not authorized, login them first.
		if !isAuthorized {
			redirectToLogin()
			return
		}

		type response struct {
			QueryString string
		}
		res := response{q.Encode()}
		tpl.Render(w, "consent", res)
	}

	postAuthorize := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := r.Context()

		// User needs to have a session in order to call the post
		// authorize endpoint.
		sess, err := sessMgr.GetSession(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// Attach the user_id to the context.
		ctx = oidc.SetUserIDContextKey(ctx, sess.UserID)

		// Construct the request payload from the querystring.
		var req oidc.AuthenticationRequest
		if err := req.FromQueryString(r.URL.Query()); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// Attempt to authenticate the user.
		res, err := svc.core.Authenticate(ctx, &req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		u, err := buildURL(req.RedirectURI, res.ToQueryString())
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		http.Redirect(w, r, u, http.StatusFound)
	}
	postToken := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := r.Context()

		// Checks if the user has an active session. If the session is
		// not active, user does not have the right credentials to
		// access the token endpoint.
		sess, err := sessMgr.GetSession(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		ctx = oidc.SetUserIDContextKey(ctx, sess.UserID)

		// Put the extra data in the context to be validated by the
		// service.
		authorization := r.Header.Get("Authorization")
		ctx = oidc.SetAuthContextKey(ctx, authorization)

		// Parse request body.
		var req oidc.AccessTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		res, err := svc.core.Token(ctx, &req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		json.NewEncoder(w).Encode(res)
	}

	// Index endpoints.
	r.GET("/", indexController.GetIndex)

	// User endpoints.
	r.POST("/logout", userController.PostLogout)
	r.GET("/register", userController.GetRegister)
	r.GET("/login", userController.GetLogin)
	r.POST("/login", userController.PostLogin)
	r.POST("/register", userController.PostRegister)

	// Client endpoints.
	r.GET("/connect/register", clientController.GetClientRegister)
	r.POST("/connect/register", clientController.PostClientRegister)

	// OpenID Connect endpoints.
	r.GET("/authorize", getAuthorize)
	r.POST("/authorize", postAuthorize)
	r.POST("/token", postToken)

	srv := gsrv.New(*port, r)
	<-srv
	log.Println("Gracefully shutdown HTTP server.")
}

func decodeBase64(in string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(in)
	return string(b), err
}

func encodeBase64(in string) string {
	return base64.URLEncoding.EncodeToString([]byte(in))
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(M{
		"error": err.Error(),
	})
}

func buildURL(uri string, q url.Values) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// getHost tries its best to return the request host.
func getHost(r *http.Request) *url.URL {
	u := r.URL

	// The scheme is http because that's the only protocol your server handles.
	u.Scheme = "http"

	// If client specified a host header, then use it for the full URL.
	u.Host = r.Host

	// Otherwise, use your server's host name.
	if u.Host == "" {
		u.Host = "your-host-name.com"
	}
	// if r.URL.IsAbs() {
	//         host := r.Host
	//         // Slice off any port information.
	//         if i := strings.Index(host, ":"); i != -1 {
	//                 host = host[:i]
	//         }
	//         u.Host = host
	// }

	return u
}
