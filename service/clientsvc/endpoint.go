package clientsvc

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/utils/encoder"
	"github.com/alextanhongpin/go-openid/utils/validator"
	"github.com/julienschmidt/httprouter"
)

// TODO: Check client registration full logic
// TODO: Add validation - only registered users with auth token can create clients
// TODO: Access token creation and validation
var (
	errTypeAssertion         = errors.New("Unhandled type assertion")
	errInvalidContentType    = errors.New("Content-Type must be application/json")
	errInvalidAcceptHeader   = errors.New("Accept must be application/json")
	errInvalidRedirectURI    = errors.New("One or more redirect_uri values are invalid")
	errInvalidClientMetadata = errors.New("The value of one of the Client Metadata fields is invalid and the server has rejected this request. Note that an Authorization Server MAY choose to substitute a valid value for any requested parameter of a Client's Metadata")
)

// Endpoint implements the service
type Endpoint func(request interface{}) (response interface{}, err error)

// GetClientEndpoint  		 	Returns a client by id
// GetClientsEndpoint 		 	Returns a list of clients
// GetClientRegisterEndpoint 	Display the registration page for client
// PostClientRegisterEndpoint 	Create a new client
// UpdateClientEndpoint			Update a client
// DeleteClientEndpoint			Delete a client

// Endpoints defines the individual endpoints for our dynamic client registration service
type Endpoints struct {
	GetClientEndpoint          Endpoint
	GetClientsEndpoint         Endpoint
	PostClientRegisterEndpoint Endpoint
	UpdateClientEndpoint       Endpoint
	DeleteClientEndpoint       Endpoint
	// GetClientMetadataEndpoint Endpoint
}

// MakeServerEndpoints creates an endpoint that exposes the route through an api endpoint
func MakeServerEndpoints(s Service) *Endpoints {
	return &Endpoints{
		GetClientEndpoint:          MakeGetClientEndpoint(s),
		GetClientsEndpoint:         MakeGetClientsEndpoint(s),
		PostClientRegisterEndpoint: MakePostClientRegisterEndpoint(s),
		UpdateClientEndpoint:       MakeUpdateClientEndpoint(s),
		DeleteClientEndpoint:       MakeDeleteClientEndpoint(s),
	}
}

// MakeGetClientEndpoint create the endpoint that returns a client
func MakeGetClientEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(getClientMetadataRequest)
		if !ok {
			return req, errTypeAssertion
		}
		return s.GetClientMetadata(req)
	}
}

// MakeGetClientsEndpoint create the endpoint that returns a list of clients
func MakeGetClientsEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(getClientsMetadataRequest)
		if !ok {
			return req, errTypeAssertion
		}
		return s.GetClientsMetadata(req)
	}
}

// MakePostClientRegisterEndpoint create the post registration endpoint
func MakePostClientRegisterEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(postClientMetadataRequest)
		if !ok {
			return req, errTypeAssertion
		}
		return s.PostClientMetadata(req)
	}
}

// MakeUpdateClientEndpoint create an endpoint to update client
func MakeUpdateClientEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(updateClientMetadataRequest)
		if !ok {
			return req, errTypeAssertion
		}
		return s.UpdateClientMetadata(req)
	}
}

// MakeDeleteClientEndpoint create an endpoint to delete a client
func MakeDeleteClientEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(deleteClientMetadataRequest)
		if !ok {
			return req, errTypeAssertion
		}

		return s.DeleteClientMetadata(req)
	}
}

// GetClientConnectView will render the registration page
func (e Endpoints) GetClientConnectView(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Check if authorization cookie is present through a middleware
		// Only authorized user can create a new client
		t.Render(w, "client_register", nil)
	}
}

// GetClientView will render the client page
func (e Endpoints) GetClientView(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := getClientMetadataRequest{ID: ps.ByName("id")}
		res, err := e.GetClientEndpoint(req)

		log.Printf("GetClientView %v", res)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		// t.Render(w, "client", res)
		t.Render(w, "client", res)
	}
}

// GetClientsView will render the clients page
func (e Endpoints) GetClientsView(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		req := getClientsMetadataRequest{}
		res, err := e.GetClientsEndpoint(req)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		t.Render(w, "clients", res)
	}
}

// GetClient will return a client by id
func (e Endpoints) GetClient() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := getClientMetadataRequest{ID: ps.ByName("id")}
		res, err := e.GetClientEndpoint(req)
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		encoder.JSON(w, res, http.StatusOK)
	}
}

// GetClients will return a list of clients
func (e Endpoints) GetClients() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		req := getClientsMetadataRequest{}
		res, err := e.GetClientsEndpoint(req)
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		encoder.JSON(w, res, http.StatusOK)
	}
}

// PostClient will create a new client
func (e Endpoints) PostClient() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req postClientMetadataRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		log.Printf("\nPostClient type=endpoint req=%v", req)
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			log.Println("PostClient type=endpoint error=invalid_content_type")
			encoder.ErrorJSON(w, errInvalidContentType, http.StatusBadRequest)
			return
		}

		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			log.Println("PostClient type=endpoint error=invalid_accept_header")
			encoder.ErrorJSON(w, errInvalidAcceptHeader, http.StatusBadRequest)
			return
		}
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		bearerType := authHeader[0]

		log.Printf("PostClient type=endpoint auth_header=%v bearer=%v \n", authHeader, bearerType)

		if bearerType != "Bearer" {
			// Handle error
		}
		// accessToken := authHeader[1]
		// accessToken is a jwt token, decode it and check if the user exists
		// userID, err := deocode(accessToken)
		//   Authorization: Bearer eyJhbGciOiJSUzI1NiJ9.eyJ ...

		// Invalid redirect url?
		// Invalid client metadata?
		// Check if name exists?
		// errInvalidClientMetadata

		// if err != nil {
		// 	encoder.ErrorJSON(w, err, http.StatusNotFound)
		// 	return
		// }

		res, err := e.PostClientRegisterEndpoint(req)
		log.Printf("PostClient type=endpoint error=%v \n", err)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusNotFound)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		encoder.JSON(w, res, http.StatusCreated)
	}
}

// UpdateClient will update a client
func (e Endpoints) UpdateClient() httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Println("UpdateClient type=endpoint event=start")
		var req updateClientMetadataRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			log.Printf("UpdateClient type=endpoint event=error error=decoding error %v \n", err)
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		req.ID, err = validator.ValidateID(ps.ByName("id"))
		if err != nil {
			log.Printf("UpdateClient type=endpoint event=error error=decoding error %v \n", err)
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		log.Printf("UpdateClient type=endpoint event=decode_success params=%v \n", req)
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			log.Println("UpdateClient type=endpoint error=invalid_content_type")
			encoder.ErrorJSON(w, errInvalidContentType, http.StatusBadRequest)
			return
		}

		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			log.Println("UpdateClient type=endpoint error=invalid_accept_header")
			encoder.ErrorJSON(w, errInvalidAcceptHeader, http.StatusBadRequest)
			return
		}
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		bearerType := authHeader[0]

		log.Printf("UpdateClient type=endpoint auth_header=%v bearer=%v \n", authHeader, bearerType)

		if bearerType != "Bearer" {
			// Handle error
		}

		res, err := e.UpdateClientEndpoint(req)
		log.Printf("UpdateClient type=endpoint error=%v \n", err)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusNotFound)
			return
		}
		log.Printf("UpdateClient type=endpoint event=update_success res=%v \n", res)

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		encoder.JSON(w, res, http.StatusCreated)
	})
}

// DeleteClient will delete a client
func (e Endpoints) DeleteClient() httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := deleteClientMetadataRequest{
			ID: bson.ObjectIdHex(ps.ByName("id")),
		}
		log.Printf("DeleteClient type=endpoint req=%v \n", req)
		res, err := e.DeleteClientEndpoint(req)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		encoder.JSON(w, res, http.StatusOK)
	})
}
