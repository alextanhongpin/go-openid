package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/core"
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/querystring"
	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/alextanhongpin/go-openid/service"

	"github.com/julienschmidt/httprouter"
)

type Core struct {
	service  service.Core
	template *html5.Template
	session  *session.Manager
}

func NewCore() Core {
	return Core{
		service: core.New(),
	}
}

func (c *Core) SetService(s service.Core) {
	c.service = s
}

func (c *Core) SetTemplate(h *html5.Template) {
	c.template = h
}

func (c *Core) SetSession(s *session.Manager) {
	c.session = s
}

func (c *Core) GetAuthorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()

	var req oidc.AuthenticationRequest
	if err := querystring.Decode(q, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.service.PreAuthenticate(&req); err != nil {
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

	isAuthorized := c.session.HasSession(r)
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
	c.template.Render(w, "consent", res)
}

func (c *Core) PostAuthorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	// User needs to have a session in order to call the post
	// authorize endpoint.
	sess, err := c.session.GetSession(r)
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
	res, err := c.service.Authenticate(ctx, &req)
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

func (c *Core) PostToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	// Checks if the user has an active session. If the session is
	// not active, user does not have the right credentials to
	// access the token endpoint.
	sess, err := c.session.GetSession(r)
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

	res, err := c.service.Token(ctx, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	json.NewEncoder(w).Encode(res)
}