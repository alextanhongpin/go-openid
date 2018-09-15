package cookie

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// CustomCookie implements a cookie with default values
type CustomCookie struct{}

var errCookieNotFound = errors.New("No cookie with the name found")

// Set will create a new cookie
func (c CustomCookie) Set(w http.ResponseWriter, name, value string, duration time.Time) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  duration,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   604800, // 1 week
		// Secure:   true,
	}
	http.SetCookie(w, &cookie)
	log.Println("successfully set cookie %#v", cookie)
}

func (c CustomCookie) Get(r *http.Request, name string) (value string, err error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", errCookieNotFound
	}
	return cookie.Value, nil
}

// Clear removes the cookie by name
func (c CustomCookie) Clear(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		MaxAge:   0,
		Path:     "/",
		// Secure:   true,
	}
	http.SetCookie(w, &cookie)
}

// Getting a list of cookies
// for k, c := range r.Cookies() {
// 	fmt.Println("cookie", k, c, c.Name)
// }

// Getting a cookie by name
// cookie, err := r.Cookie("access_token")

// New creates a new cookie
func New() *CustomCookie {
	return &CustomCookie{}
}
