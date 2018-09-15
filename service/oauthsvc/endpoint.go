package oauthsvc

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/utils/encoder"
	"github.com/julienschmidt/httprouter"
)

type Endpoint func(request interface{}) (response interface{}, err error)

type Endpoints struct {
	GetAuthorizeViewEndpoint Endpoint
	PostAuthorizeEndpoint    Endpoint
}

var (
	errInteractionRequired      = errors.New("The Authorization Server requires End-User interaction of some form to proceed")
	errLoginRequired            = errors.New("The Authorization Server requires End-User authentication")
	errAccountSelectionRequired = errors.New("The End-User is REQUIRED to select a session at the Authorization Server")
	errConsentRequired          = errors.New("The Authorization Server requires End-User consent")
	errInvalidRequestURI        = errors.New("The request_uri in the Authorization Request returns an error or contains invalid data")
	errInvalidRequestObject     = errors.New("The request parameter contains an invalid Request Object")
	errRequestNotSupported      = errors.New("The OP does not support use of the request parameter")
	errRequestURINotSupported   = errors.New("The OP does not support use of the request_uri parameter")
	errRegistrationNotSupported = errors.New("The OP does not support use of the registration parameter")
	// Others
	errNoRedirectURIs      = errors.New("The redirect_uri is not present in client")
	errInvalidRedirectURI  = errors.New("The redirect_uri does not match the client redirect_uri")
	errInvalidScope        = errors.New("One or more of the scopes provided is invalid")
	errInvalidResponseType = errors.New("The response_type is invalid")
)

func MakeServerEndpoints(s Service) *Endpoints {
	return &Endpoints{
		GetAuthorizeViewEndpoint: MakeGetAuthorizeViewEndpoint(s),
		PostAuthorizeEndpoint:    MakePostAuthorizeEndpoint(s),
		// GetAuthorizeEndpoint:     MakeGetAuthorizeEndpoint(s),
	}
}

// MakeGetAuthorizeViewEndpoint creates an endpoint for displaying the authorize endpoint view
func MakeGetAuthorizeViewEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		// 1. Check if the user is logged in or not through cookie
		// 2. If not, redirect them to login
		// 3. Check the redirect URI
		// 4. Check the scopes
		req := request.(getAuthorizeRequest)
		res := getAuthorizeResponse{}

		// oid, err := utils.ValidateId(req.ClientID)
		// if err != nil {
		// 	return &res, err
		// }

		// TODO: Do a checking and change the implementation based on different response_type
		// response_type can be code (Authorization Code Flow), implicit (Mobile) or hybrid (Hybrid Flow)
		if req.ResponseType != "code" {
			return &res, errInvalidResponseType
		}

		client, err := s.GetClient(getClientRequest{ClientID: req.ClientID})
		if err != nil {
			return &res, err
		}
		log.Printf("MakeGetAuthorizeViewEndpoint type=make_endpoint client=%v\n", client)
		// 3
		redirectURI, err := validateRedirectURI(client.Data.RedirectURIs, req.RedirectURI)
		if err != nil {
			return &res, err
		}
		log.Printf("MakeGetAuthorizeViewEndpoint type=make_endpoint redirect_uri=%v\n", redirectURI)
		// 4
		scopes, err := validateScopes(req.Scope)
		if err != nil {
			return &res, err
		}
		log.Printf("MakeGetAuthorizeViewEndpoint type=make_endpoint valid_scopes=%v\n", scopes)

		return getAuthorizeResponse{
			Authorize: Authorize{
				Scope:        req.Scope,
				ResponseType: req.ResponseType,
				ClientID:     req.ClientID,
				State:        req.State,
				RedirectURI:  redirectURI,
			},
			Scopes: scopes,
			URL:    req.URL,
		}, nil
	}
}

func MakePostAuthorizeEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req := request.(postAuthorizeRequest)
		res := postAuthorizeResponse{}
		// oid, err := utils.ValidateId(req.ClientID)
		// if err != nil {
		// 	return &res, err
		// }

		if req.ResponseType != "code" {
			return &res, errInvalidResponseType
		}

		client, err := s.GetClient(getClientRequest{ClientID: req.ClientID})
		if err != nil {
			return &res, err
		}
		log.Printf("MakePostAuthorizeEndpoint type=make_endpoint client=%v\n", client)
		// 3
		redirectURI, err := validateRedirectURI(client.Data.RedirectURIs, req.RedirectURI)
		if err != nil {
			return &res, err
		}
		log.Printf("MakePostAuthorizeEndpoint type=make_endpoint redirect_uri=%v\n", redirectURI)
		// 4
		scopes, err := validateScopes(req.Scope)
		if err != nil {
			return &res, err
		}
		log.Printf("MakePostAuthorizeEndpoint type=make_endpoint valid_scopes=%v\n", scopes)
		response, err := s.PostAuthorize(postAuthorizeRequest{
			Authorize: Authorize{
				ClientID:    req.ClientID,
				RedirectURI: redirectURI,
				State:       req.State,
			},
		})
		if err != nil {
			return response, err
		}
		log.Printf("MakePostAuthorizeEndpoint type=make_endpoint response=%#v", response)
		return response, nil
	}
}

func (e Endpoints) GetAuthorizeView(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log.Printf("GetAuthorizeView type=endpoint url=%v\n", r.URL.String())
		req := getAuthorizeRequest{
			Authorize{
				Scope:        r.URL.Query().Get("scope"),
				ResponseType: r.URL.Query().Get("response_type"),
				ClientID:     r.URL.Query().Get("client_id"),
				RedirectURI:  r.URL.Query().Get("redirect_uri"),
				State:        r.URL.Query().Get("state"),
			},
			r.URL.String(),
		}
		res, err := e.GetAuthorizeViewEndpoint(req)
		if err != nil {
			log.Printf("GetAuthorizeView type=endpoint error=%v\n", err)
			// If there is an error, should redirect the user instead
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		t.Render(w, "consent", res)
		// If successful, redirect the user to the UI
		// Then load the consent screen
	}
}

func (e Endpoints) PostAuthorize() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// returns code and state if success
		// Code will expire after 5 minutes/or after it has been used once
		log.Printf("PostAuthorize type=endpoint url=%v\n", r.URL.String())
		req := postAuthorizeRequest{
			Authorize: Authorize{
				Scope:        r.URL.Query().Get("scope"),
				ResponseType: r.URL.Query().Get("response_type"),
				ClientID:     r.URL.Query().Get("client_id"),
				RedirectURI:  r.URL.Query().Get("redirect_uri"),
				State:        r.URL.Query().Get("state"),
			},
		}
		log.Printf("PostAuthorize type=endpoint req=%v\n", req)
		response, err := e.PostAuthorizeEndpoint(req)
		if err != nil {
			// If there is an error, should redirect the user instead
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
		}
		res := response.(*postAuthorizeResponse)
		urlStr, err := res.genURL(r.URL.Query().Get("redirect_uri"))
		res.RedirectURI = urlStr
		if err != nil {
			// If there is an error, should redirect the user instead
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
		}
		// http.Redirect(w, r, urlStr, http.StatusFound)
		// Redirect to callback url with authorization code and state
		encoder.JSON(w, res, http.StatusOK)
	}
}

func validateScopes(scopes string) ([]string, error) {
	splitScopes := strings.Split(scopes, " ")

	allowedScopes := make(map[string]int)
	allowedScopes["openid"] = 2
	allowedScopes["profile"] = 2
	allowedScopes["address"] = 2
	allowedScopes["phone"] = 2

	for _, scope := range splitScopes {
		value, ok := allowedScopes[scope]
		if !ok {
			allowedScopes[scope] = 0
		}
		allowedScopes[scope] = value + 1
	}

	validScopes := []string{}
	for s, score := range allowedScopes {
		if score == 1 {
			// Bad, there are some wrong one
			return []string{}, errInvalidScope
		} else if score == 3 {
			// This are the ones specified
			validScopes = append(validScopes, s)
		}
	}
	return validScopes, nil
}

func validateRedirectURI(redirectURIs []string, redirectURI string) (string, error) {
	// No redirect uris, return error
	if len(redirectURIs) == 0 {
		return "", errNoRedirectURIs
	}
	for _, ruri := range redirectURIs {
		if ruri == redirectURI {
			return redirectURI, nil
		}
	}
	return "", errInvalidRedirectURI
}
