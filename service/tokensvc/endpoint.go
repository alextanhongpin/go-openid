package tokensvc

import (
	"errors"
	"log"
	"net/http"

	"strings"

	"encoding/base64"
	"encoding/json"

	"github.com/alextanhongpin/go-openid/utils/encoder"
	"github.com/julienschmidt/httprouter"
)

var (
	errInvalidGrantType          = errors.New("Invalid grant type")
	errInvalidCode               = errors.New("Code has expired")
	errInvalidRedirectURI        = errors.New("Invalid redirect uri")
	errInvalidContentType        = errors.New("Invalid content type")
	errUnauthorized              = errors.New("Unauthorized user")
	errInvalidBasicAuthorization = errors.New("Invalid basic authorization")
	errInvalidClientSecret       = errors.New("Invalid client secret")
)

// Endpoint is the interface
type Endpoint func(request interface{}) (response interface{}, err error)

type Endpoints struct {
	GetTokenEndpoint  Endpoint
	PostTokenEndpoint Endpoint
}

// MakeServerEndpoints creates the endpoints for the token services
func MakeServerEndpoints(s Service) *Endpoints {
	return &Endpoints{
		// GetTokenEndpoint is the endpoint
		GetTokenEndpoint:  MakeGetTokenEndpoint(s),
		PostTokenEndpoint: MakePostTokenEndpoint(s),
	}
}

// MakeGetTokenEndpoint creates and endpoint for getting the token
func MakeGetTokenEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		return nil, nil
	}
}

// MakePostTokenEndpoint creates and endpoint for creating the token
func MakePostTokenEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req := request.(postTokenRequest)
		log.Printf("MakePostTokenEndpoint request=%#v\n", req)
		res := postTokenResponse{}
		if req.GrantType != "authorization_code" {
			return res, errInvalidGrantType
		}

		if req.Code == "" {
			// Handle error
			return res, errInvalidCode
		}
		log.Println("MakePostTokenEndpoint message=checking code")
		valid, err := s.CheckCode(codeRequest{
			ClientID: req.ClientID,
			Code:     req.Code,
		})
		if err != nil {
			return nil, errInvalidCode
		}
		if !valid.Exist {
			return nil, errInvalidCode
		}

		// Decode access token to get client id and client secret
		// Query client from database
		client, err := s.CheckClient(clientRequest{ID: req.ClientID})
		if err != nil {
			return nil, err
		}
		if client.Data.ClientSecret != req.ClientSecret {
			return nil, errInvalidClientSecret
		}

		hasClient := false
		for _, v := range client.Data.ClientMetadata.RedirectURIs {
			if v == req.RedirectURI {
				hasClient = true
				break
			}
		}
		if !hasClient {
			return res, nil
		}

		// Generate access token
		// Generate refresh token
		// Generate id token

		out, err := s.PostToken(req)
		res = *out
		return res, nil
	}
}

func (e Endpoints) PostToken() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//  POST /token HTTP/1.1
		//   Host: server.example.com
		//   Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
		//   Content-Type: application/x-www-form-urlencoded

		// if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		// 	// Error occured
		// 	encoder.ErrorJSON(w, errInvalidContentType, http.StatusBadRequest)
		// 	return
		// }

		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		bearerType := authHeader[0]
		if strings.ToLower(bearerType) != "basic" {
			encoder.ErrorJSON(w, errUnauthorized, http.StatusBadRequest)
			return
		}

		accessToken := authHeader[1]
		decodedClientMetadata, err := base64.StdEncoding.DecodeString(accessToken)
		if err != nil {
			encoder.ErrorJSON(w, errInvalidBasicAuthorization, http.StatusBadRequest)
			return
		}
		req := postTokenRequest{}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		// composed of client secret and client id
		clientMetadata := strings.Split(string(decodedClientMetadata), ":")

		req.ClientID = clientMetadata[0]
		req.ClientSecret = clientMetadata[1]
		res, err := e.PostTokenEndpoint(req)

		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		w.Header().Set("Cache-Control", "no-cache, no-store")
		w.Header().Set("Pragma", "no-cache")
		encoder.JSON(w, res, http.StatusOK)
		// 	HTTP/1.1 200 OK
		//   Content-Type: application/json
		//   Cache-Control: no-cache, no-store
		//   Pragma: no-cache

		//   {
		//    "access_token":"SlAV32hkKG",
		//    "token_type":"Bearer",
		//    "expires_in":3600,
		//    "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
		//    "id_token":"eyJ0 ... NiJ9.eyJ1c ... I6IjIifX0.DeWt4Qu ... ZXso"
		//   }

	}
}
