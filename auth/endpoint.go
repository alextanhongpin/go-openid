package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/alextanhongpin/go-openid/app"

	"github.com/julienschmidt/httprouter"
)

var (
	errorTemplateNotFound = errors.New("Template not found")
)

type endpoint struct {
	svc service
}

type Endpoint func(request interface{}) (response interface{}, err error)

// Endpoints  exposes
type Endpoints struct {
	// GetUserEndpoint does this
	GetUserEndpoint Endpoint //httprouter.Handle
}

// MakeServerEndpoints doe
func MakeServerEndpoints(s Service) *Endpoints {
	return &Endpoints{
		GetUserEndpoint: MakeGetUserEndpoint(s),
	}
}

// MakeGetUserEndpoint asd
func MakeGetUserEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req := request.(getUserRequest)
		user, err := s.GetUser(req.ID)
		if err != nil {
			return getUserResponse{Name: "something"}, err
		}
		return user, nil
	}
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	ps := r.Context().Value("params").(httprouter.Params)
	// 	var req getUserRequest
	// 	req.ID = ps.ByName("id")
	// 	user, err := s.GetUser(req.ID)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	// No user found
	// 	if user == nil {
	// 		w.Write([]byte("no user found"))
	// 		return
	// 	}

	// 	json.NewEncoder(w).Encode(user)
	// })
}

// Middleware example
// func (e Endpoints) GetUser(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// do stuff
// 		log.Println("At middleware one:start")
// 		params := r.Context().Value("params").(httprouter.Params)
// 		_, err := e.GetUserEndpoint(getUserRequest{ID: params.ByName("id")})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 		log.Println("At middleware one:end")
// 	})
// }

// GetUser returns
func (e Endpoints) GetUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := e.GetUserEndpoint(getUserRequest{ID: ps.ByName("id")})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		encodeResponse(w, user)
	}
}

// Encodes the struct and returns the data as json
func encodeResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ErrorJson returns the error
type ErrorJson struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

// encodeError returns a json error
func encodeError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(message)
	json.NewEncoder(w).Encode(ErrorJson{
		Error:   code,
		Message: message,
	})
}

type getUserRequest struct {
	ID string
}

type getUserResponse struct {
	Name string
}

// GET api/users/:id
func (e endpoint) getUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := e.svc.fetchOne(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// No user found
		if user == nil {
			w.Write([]byte("no user found"))
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

type loginResponse struct {
	OK          bool   `json:"ok"`
	RedirectURI string `json:"redirect_uri"`
}

func (e endpoint) loginHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var user *User

		// Decode the post payload
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the user exists
		user, err = e.svc.checkExist(user.Email)

		if err != nil {
			log.Printf("Error user not found err=%s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// User does not exist, throw the correct error
		if user == nil {
			log.Println("User does not exist")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("User found!")
		// User exists, redirect
		w.Header().Set("Content-Type", "application/json")

		response := loginResponse{
			OK:          true,
			RedirectURI: "/users/" + user.ID.Hex(),
		}
		json.NewEncoder(w).Encode(response)
	}
}

func (e endpoint) loginViewHandler(tmpl *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tmpl.Render(w, "login", nil)
	}
}

func (e endpoint) viewUserHandler(tmpl *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := e.svc.fetchOne(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// No user found
		if user == nil {
			w.Write([]byte("no user found"))
			return
		}
		log.Println("rendering user view")
		tmpl.Render(w, "user", user)
	}
}

// POST /register
// registerHandle register a user and return a redirect uri
type registerResponse struct {
	OK          bool   `json:"ok"`
	RedirectURI string `json:"redirect_uri"`
	AccessToken string `json:"access_token"`
}

func (e endpoint) registerHandler(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var u User

		// Decode the body payload
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("auth/endpoint.go -> registerHandler -> u=%v", u)

		// Check if the user exist
		user, err := e.svc.checkExist(u.Email)
		log.Printf("auth/*registerHandler user=%v", user)

		// Error occured
		if err != nil {
			log.Printf("auth/*registerHandler error=%s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// No user found, create new account
		if user == nil {
			log.Println("registerHandler: No user found. Creating one")
			userID, err := e.svc.create(u)
			if err != nil {
				log.Printf("registerHandler: Error creating user: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			// Successfully created, return payload
			log.Println("registerHandler: successfully created user")
			ctx := context.WithValue(r.Context(), "userID", userID)
			next(w, r.WithContext(ctx), ps)
			// CReate a new jsonweb token
			// w.WriteHeader(http.StatusCreated)
			// w.Header().Set("Content-Type", "application/json")

			// response := registerResponse{
			// 	OK:          true,
			// 	RedirectURI: "/users/" + userID,
			// }
			// json.NewEncoder(w).Encode(response)
			return
		}
		log.Printf("registerHandler: user exist %+v", user)
		// w.WriteHeader(http.StatusUnauthorized)
		// w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(`{"message": "user already exists"}`))

		http.Error(w, `{"success": false, "message": "user already exists"}`, http.StatusUnauthorized)
	}
}

// GET register/
// registerViewHandler renders the register view
func (e endpoint) registerViewHandler(tmpl *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tmpl.Render(w, "register", nil)
	}
}

// GET api/users
func (e endpoint) getUsersHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		users := e.svc.fetchMany()
		json.NewEncoder(w).Encode(users)
	}
}
