package client

import (
	"encoding/json"
	"net/http"

	"github.com/alextanhongpin/go-openid/pkg/querystring"
)

// Make a http request to the authorization server.
func Authorize(w http.ResponseWriter, r *http.Request) {
	req := oidc.AuthenticationRequest{
		ResponseType: "code",
		Scope:        "openid profile email",
		ClientID:     "xyz",
		State:        "abc",
		RedirectURI:  "https://client.example.com",
	}
	req.Validate()

	// It is basically a redirect to the consent page.
	qs := querystring.Encode(uri.Values{}, req)
	http.Redirect(w, qs.Encode(), http.StatusFound)

	// If the user is logged-in, use the credentials.
	// Else, request user to log in.
	// Once the user is logged in, ask for their consent on the scopes.
	// Make a post request.
}

func AuthorizeCallback(w http.ResponseWriter, r *http.Request) {
	// Authorization response.
	var req oidc.AuthenticationResponse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// Do something.
	}

	// Exchange with the token.

}
