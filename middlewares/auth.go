package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"

	"github.com/alextanhongpin/go-openid/models"
	"github.com/alextanhongpin/go-openid/utils/cookie"
	"github.com/alextanhongpin/go-openid/utils/encoder"
)

// Key is the type of the context
type Key string

var userID Key

const cookieName string = "auth"

var errUnauthorizedUser = errors.New("An API key is required to access this api")

// ValidateAuth checks if the user has the access token cookie that is required to login
func ValidateAuth(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Println("validateAuth type=middleware event=check_cookie")
		c := cookie.New()
		auth, err := c.Get(r, cookieName)
		log.Println("getting cookie", auth, err)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		token, err := jwt.ParseWithClaims(auth, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("$ecret"), nil
		})
		if err != nil {
			log.Printf("validateAuth type=middleware event=validate_token err=%v \n", err)
			c.Clear(w, cookieName)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), userID, *claims)
			next(w, r.WithContext(ctx), ps)
		} else {
			// Token expired
			c.Clear(w, cookieName)
			log.Printf("validateAuth type=middleware event=validate_token ok=%v claims=%v \n", ok, claims)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	})
}

// ProtectRoute prevents authenticated users to access
func ProtectRoute(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Printf("Protect type=middleware event=read_cookie url=%v", r.URL)
		c := cookie.New()
		auth, err := c.Get(r, cookieName)

		// The token is not found
		if err != nil {
			// log.Println("Protect type=middleware event=check_cookie message=cookie not found")
			// http.Error(w, err.Error(), http.StatusBadRequest)
			// next(w, r, ps)
			log.Printf("Protect type=middleware event=check_url url=%v host=%v", r.URL, r.Host)

			// Display the Login or Register page, else redirect them to the Login page
			if r.URL.String() == "/login" || r.URL.String() == "/register" || r.URL.String() == "/" {
				next(w, r, ps)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		token, err := jwt.ParseWithClaims(auth, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("$ecret"), nil
		})
		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			// ctx := context.WithValue(r.Context(), userID, *claims)
			// next(w, r.WithContext(ctx), ps)
			// User already have a valid token
			userID := claims.UserID
			http.Redirect(w, r, "/users/"+userID, http.StatusFound)
			return
		}
		// c.Clear(w, cookieName)
		next(w, r, ps)
	})
}

// ProtectAPI prevents unauthorized users from accessing the api endpoint
func ProtectAPI(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		log.Println("ProtectAPI type=middleware event=check_cookie")
		c := cookie.New()
		auth, err := c.Get(r, cookieName)

		if err != nil {
			log.Printf("ProtectAPI type=middleware event=unauthorized message=access token not found error=%v", err)
			encoder.ErrorJSON(w, errUnauthorizedUser, http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(auth, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("$ecret"), nil
		})
		if err != nil {
			log.Printf("ProtectAPI type=middleware event=validate_token err=%v \n", err)
			c.Clear(w, cookieName)
			encoder.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			log.Printf("ProtectAPI type=middleware event=access_token_present ok=%v claims=%v", ok, *claims)
			ctx := context.WithValue(r.Context(), userID, *claims)
			next(w, r.WithContext(ctx), ps)
			return
		}
		c.Clear(w, cookieName)
		encoder.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	})
}
