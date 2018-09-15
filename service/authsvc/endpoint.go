package authsvc

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	tracelog "github.com/opentracing/opentracing-go/log"
	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/models"
	"github.com/alextanhongpin/go-openid/utils/cookie"
	"github.com/alextanhongpin/go-openid/utils/encoder"
	"github.com/alextanhongpin/go-openid/utils/validator"
)

// Endpoint represents the route that will be called by the server.
type Endpoint func(request interface{}) (response interface{}, err error)

// Endpoints exposes a list of endpoints that will be implemented by the server.
type Endpoints struct {
	GetUserEndpoint          Endpoint
	GetUsersEndpoint         Endpoint
	DeleteUserEndpoint       Endpoint
	UpdateUserEndpoint       Endpoint
	PostRegisterEndpoint     Endpoint
	GetLoginCallbackEndpoint Endpoint
	PostLoginEndpoint        Endpoint
	GetUserViewEndpoint      Endpoint
	GetUserEditViewEndpoint  Endpoint
	GetUsersViewEndpoint     Endpoint
}

// MakeServerEndpoints is a helper to create all the endpoints and
// insert the service through dependency injection.
func MakeServerEndpoints(s Service) *Endpoints {
	return &Endpoints{
		GetUserEndpoint:          MakeGetUserEndpoint(s),
		GetUsersEndpoint:         MakeGetUsersEndpoint(s),
		DeleteUserEndpoint:       MakeDeleteUserEndpoint(s),
		UpdateUserEndpoint:       MakeUpdateUserEndpoint(s),
		PostRegisterEndpoint:     MakePostRegisterEndpoint(s),
		GetLoginCallbackEndpoint: MakeGetLoginCallbackEndpoint(s),
		PostLoginEndpoint:        MakePostLoginEndpoint(s),
		GetUserViewEndpoint:      MakeGetUserViewEndpoint(s),
		GetUserEditViewEndpoint:  MakeGetUserEditViewEndpoint(s),
		GetUsersViewEndpoint:     MakeGetUsersViewEndpoint(s),
	}
}

var (
	errTypeAssertion = errors.New("Unhandled type assertion")
	errUserNotFound  = errors.New("No user with the email found")
)

// MakeGetUserEndpoint creates an endpoint that makes a call to get a single user
func MakeGetUserEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(getUserRequest)
		if !ok {
			return nil, errTypeAssertion
		}
		_, err := validator.ValidateID(req.ID)
		if err != nil {
			return nil, err
		}
		return s.GetUser(req)
	}
}

// MakeGetUsersEndpoint creates an endpoint to get a list of users
func MakeGetUsersEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(getUsersRequest)
		if !ok {
			return nil, errTypeAssertion
		}
		return s.GetUsers(req)
	}
}

// MakeUpdateUserEndpoint creates the endpoint for updating user
func MakeUpdateUserEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(updateUserRequest)
		if !ok {
			return nil, errTypeAssertion
		}
		return s.UpdateUser(req)
	}
}

// MakeDeleteUserEndpoint creates an endpoint for deleting user
func MakeDeleteUserEndpoint(s Service) Endpoint {
	return func(request interface{}) (res interface{}, err error) {
		req := request.(deleteUserRequest)
		return s.DeleteUser(req)
	}
}

// MakePostRegisterEndpoint handles the business logic for post register.
func MakePostRegisterEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		// log.Println("MakePostRegisterEndpoint event=start")
		req, ok := request.(postRegisterRequest)
		if !ok {
			return nil, errTypeAssertion
		}

		if valid, err := govalidator.ValidateStruct(req); err != nil || !valid {
			return nil, err
		}

		// Check if the email exist.
		user, err := s.CheckUser(req.Email)
		if err != nil {
			return nil, err
		}

		// User already exists, return error message
		if user != nil {
			// log.Println("MakePostRegisterEndpoint event=exit message=user_exist")
			// return postRegisterResponse{
			// 	Ok:    false,
			// 	Error: "User with the email already exists",
			// }, nil
			return nil, errors.New("User with the email already exists")
		}

		response, err := s.CreateUser(createUserRequest(req))
		if err != nil {
			return nil, err
		}

		userID := response.ID

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
			userID,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		})

		tokenString, err := token.SignedString([]byte("$ecret"))
		if err != nil {
			return nil, err
		}

		// log.Println("MakePostRegisterEndpoint event=token_created")
		res := postRegisterResponse{
			Ok:          true,
			UserID:      userID,
			AccessToken: tokenString,
			RedirectURI: "/users/" + userID,
		}

		// log.Printf("MakePostRegisterEndpoint event=exit user_id=%v", userID)
		return res, nil
	}
}

// MakeGetLoginCallbackEndpoint is an empty service
func MakeGetLoginCallbackEndpoint(s Service) Endpoint {
	return func(request interface{}) (res interface{}, err error) {
		req := request.(getLoginCallbackRequest)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
			req.UserID,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		})

		// TODO: Move the secret to a config file
		tokenString, err := token.SignedString([]byte("$ecret"))
		if err != nil {
			return nil, err
		}
		return getLoginCallbackResponse{
			AccessToken: tokenString,
			ExpiresIn:   3600,
		}, nil
	}
}

// MakePostLoginEndpoint contains the service orchestration to carry out login.
func MakePostLoginEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(postLoginRequest)
		if !ok {
			return nil, errTypeAssertion
		}

		// Check if the email exist.
		user, err := s.CheckUser(req.Email)

		log.Printf("MakePostLoginEndpoint type=make_endpoint user=%v\n", user)
		if err != nil {
			return nil, err
		}

		userID := ""
		// User does not exist yet
		if user == nil {
			// No account found error
			return user, errUserNotFound
			// newUserID, err := s.CreateUser(createUserRequest(req))
			// if err != nil {
			// 	return nil, err
			// }
			// userID = newUserID
		}
		userID = user.ID.Hex()

		// token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		// 	userID,
		// 	jwt.StandardClaims{
		// 		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		// 	},
		// })

		// // TODO: Move the secret to a config file
		// tokenString, err := token.SignedString([]byte("$ecret"))
		// if err != nil {
		// 	return nil, err
		// }

		res := postLoginResponse{
			// AccessToken: tokenString,
			Ok:          true,
			UserID:      userID,
			RedirectURI: "/login/callback?user_id=" + userID,
		}

		return res, nil
	}
}

// MakeGetUserViewEndpoint returns the params need to render the user page.
func MakeGetUserViewEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(getUserViewRequest)
		if !ok {
			return nil, errTypeAssertion
		}

		response, err := s.GetUser(getUserRequest(req))
		if err != nil {
			return nil, err
		}

		res := getUserViewResponse{response.Data}
		res.Password = ""
		return res, nil
	}
}

func MakeGetUserEditViewEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(getUserEditViewRequest)
		if !ok {
			return nil, errTypeAssertion
		}

		response, err := s.GetUser(getUserRequest(req))
		if err != nil {
			return nil, err
		}
		res := getUserEditViewResponse{response.Data}
		res.Password = ""
		return res, nil
	}
}
func MakeGetUsersViewEndpoint(s Service) Endpoint {
	return func(request interface{}) (interface{}, error) {
		req, ok := request.(getUsersViewRequest)
		if !ok {
			return nil, errTypeAssertion
		}

		return s.GetUsers(getUsersRequest(req))
	}
}

// METHODS

// GetUser returns a user by id.
func (e Endpoints) GetUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res, err := e.GetUserEndpoint(getUserRequest{ID: ps.ByName("id")})
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		encoder.JSON(w, res, http.StatusOK)
	}
}

// GetUsers returns a list of users.
func (e Endpoints) GetUsers() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res, err := e.GetUsersEndpoint(getUsersRequest{})
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		encoder.JSON(w, res, http.StatusOK)
	}
}

// UpdateUser update a user by id
func (e Endpoints) UpdateUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		log.Printf("UpdateUser type=endpoint request=%#v", u)

		u.ID = bson.ObjectIdHex(ps.ByName("id"))
		res, err := e.UpdateUserEndpoint(updateUserRequest{
			User: u,
		})
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		encoder.JSON(w, res, http.StatusOK)
	}
}

// GetRegister renders the register page.
func (e Endpoints) GetRegister(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t.Render(w, "register", nil)
	}
}

// PostRegister handles the post form submission.
func (e Endpoints) PostRegister(tracer *app.Tracer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Avoid cookies, probably use localStorage
		// Need cookie now for get pages, since we can't inject the cookie inside

		var req postRegisterRequest
		span := tracer.Ctx.StartSpan("register")
		span.SetOperationName("register")
		defer span.Finish()

		span.LogEvent("marshall_body")
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			span.LogFields(tracelog.String("marshall_body_error", err.Error()))
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		// span.LogEvent("validate_request")
		// if valid, err := govalidator.ValidateStruct(req); err != nil || !valid {
		// 	span.LogFields(tracelog.String("validate_request_error", err.Error()))
		// 	encoder.ErrorJSON(w, err, http.StatusBadRequest)
		// 	return
		// }

		span.LogEvent("register_service")
		response, err := e.PostRegisterEndpoint(req)
		if err != nil {
			span.LogFields(tracelog.String("register_service_error", err.Error()))
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		// Create a cookie to store the jwt token.
		res, ok := response.(postRegisterResponse)
		if !ok {
			encoder.ErrorJSON(w, errTypeAssertion, http.StatusBadRequest)
			return
		}

		span.LogEvent("set_cookie")
		cookie := cookie.New()
		cookie.Set(w, "auth", res.AccessToken, time.Now().Add(time.Hour))
		// expiration := time.Now().Add(time.Hour)
		// cookie := http.Cookie{
		// 	Name:     "Auth",
		// 	Value:    res.AccessToken,
		// 	Expires:  expiration,
		// 	HttpOnly: true,
		// 	Path:     "/",
		// 	MaxAge:   50000,
		// 	Secure:   true,
		// }
		// http.SetCookie(w, &cookie)

		span.LogEvent("register_done")
		encoder.JSON(w, res, http.StatusOK)
	}
}

// GetLogin renders the login page.
func (e Endpoints) GetLogin(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t.Render(w, "login", nil)
	}
}

func (e Endpoints) GetLoginCallback() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		userID := r.URL.Query().Get("user_id")
		req := getLoginCallbackRequest{
			UserID: userID,
		}
		response, err := e.GetLoginCallbackEndpoint(req)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		res := response.(getLoginCallbackResponse)

		// expiration := time.Now().Add(time.Hour)
		// // TODO: Set the expiration date
		// cookie := http.Cookie{
		// 	Name:     "access_token",
		// 	Value:    res.AccessToken,
		// 	Expires:  expiration,
		// 	HttpOnly: true,
		// 	MaxAge:   50000,
		// 	Path:     "/",
		// }
		// log.Printf("GetLoginCallback type=endpoint event=setting_cookie cookie=%#v \n", cookie)
		// http.SetCookie(w, &cookie)

		cookie := cookie.New()
		cookie.Set(w, "auth", res.AccessToken, time.Now().Add(time.Hour))
		log.Printf("GetLoginCallback type=endpoint event=exiting")
		http.Redirect(w, r, "/users/"+userID, http.StatusFound)
	}
}

// PostLogin handles the login submission
func (e Endpoints) PostLogin() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req postLoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		response, err := e.PostLoginEndpoint(req)
		log.Printf("PostLogin type=endpoint err=%v", err)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		res, ok := response.(postLoginResponse)
		if !ok {
			encoder.ErrorJSON(w, errTypeAssertion, http.StatusBadRequest)
			return
		}

		// expiration := time.Now().Add(time.Hour)
		// cookie := http.Cookie{
		// 	Name:     "Auth",
		// 	Value:    res.AccessToken,
		// 	Expires:  expiration,
		// 	HttpOnly: true,
		// 	MaxAge:   50000,
		// 	Path:     "/",
		// }
		// log.Printf("PostLogin type=endpoint event=setting_cookie cookie=%#v \n", cookie)
		// http.SetCookie(w, &cookie)
		// log.Printf("PostLogin type=endpoint event=exiting")
		cookie := cookie.New()
		cookie.Set(w, "auth", res.AccessToken, time.Now().Add(time.Hour))
		log.Println("PostLogin set_cookie")
		encoder.JSON(w, res, http.StatusOK)
	}
}
func (e Endpoints) DeleteUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req deleteUserRequest
		req.ID = ps.ByName("id")
		res, err := e.DeleteUserEndpoint(req)

		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		encoder.JSON(w, res, http.StatusOK)
	}
}

// GetUserView loads the user data and populate the view.
func (e Endpoints) GetUserView(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req getUserViewRequest
		req.ID = ps.ByName("id")

		res, err := e.GetUserViewEndpoint(req)
		log.Printf("GetUser type=endpoint user=%#v", res)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		t.Render(w, "user", res)
	}
}

// GetUserEditView displays the edit user view
func (e Endpoints) GetUserEditView(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req getUserEditViewRequest
		req.ID = ps.ByName("id")
		res, err := e.GetUserEditViewEndpoint(req)

		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		t.Render(w, "user-edit", res)
	}
}

// GetUsersView loads the users data and populate the view.
func (e Endpoints) GetUsersView(t *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req getUsersViewRequest
		res, err := e.GetUsersViewEndpoint(req)
		if err != nil {
			encoder.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		// res := response.(getUsersViewResponse)
		log.Printf("GetUsersView event=render_template users=%#v", res)

		t.Render(w, "users", res)
	}
}

// PostLogout clears the user's login cookie and redirect them to the login page
func (e Endpoints) PostLogout() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cookie := cookie.New()
		cookie.Clear(w, "auth")
		http.Redirect(w, r, "/login", 302)
	}
}
