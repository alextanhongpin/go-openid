package oauthsvc

import (
	"log"
	"net/url"

	"github.com/alextanhongpin/go-openid/models"
)

// Authorize the schema that is stored in the database
type Authorize struct {
	Scope        string `json:"scope" valid="required,matches([profile|email|address|phone|offline_access])"` // REQUIRED. profile|email|address|phone|offline_access
	ResponseType string `json:"response_type" valid="required"`                                               // REQUIRED.
	ClientID     string `json:"client_id" valid="required"`                                                   // REQUIRED.
	RedirectURI  string `json:"redirect_uri" valid="email,required"`                                          // REQUIRED.
	State        string `json:"state"`                                                                        // RECOMMENDED.
	ResponseMode string `json:"response_mode"`
	Nonce        string `json:"nonce"`
	Display      string `json:"display" valid="matches([page|popup|touch|wrap])"`            // page|popup|touch|wrap
	Prompt       string `json:"prompt" valid="matches([none|login|consent|select_account])"` // none|login|consent|select_account
	MaxAge       int    `json:"max_age"`
	UILocales    string `json:"ui_locales"`
	IDTokenHint  string `json:"id_token_hint"`
	LoginHint    string `json:"login_hint"`
	AcrValues    string `json:"acr_values"`
}

type authenticationError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
	State            string `json:"state"`
}

func (a authenticationError) RedirectLink(redirectURI string) string {

	u, err := url.Parse("")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "http"
	u.Host = "google.com"
	q := u.Query()
	q.Set("error", a.Error)
	q.Set("error_description", a.ErrorDescription)
	q.Set("error_uri", a.ErrorURI)
	q.Set("state", a.State)
	u.RawQuery = q.Encode()
	return u.String()
}

type getAuthorizeRequest struct {
	Authorize
	URL string `json:"url"`
}
type getAuthorizeResponse struct {
	Authorize
	Scopes []string `json:"scopes"`
	URL    string   `json:"url"`
}
type postAuthorizeRequest struct {
	Authorize
}
type postAuthorizeResponse struct {
	Code        string `json:"code,omitempty"`
	State       string `json:"state,omitempty"`
	RedirectURI string `json:"redirect_uri"`
}

func (p postAuthorizeResponse) genURL(redirectURI string) (string, error) {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("code", p.Code)
	q.Set("state", p.State)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

type getClientRequest struct {
	ClientID string
}
type getClientResponse struct {
	Data models.Client
}
