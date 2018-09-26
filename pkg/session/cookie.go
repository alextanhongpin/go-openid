package session

import (
	"net/http"
	"time"
)

// NewCookie returns a new cookie with default settings.
func NewCookie(id string) *http.Cookie {
	return &http.Cookie{
		Name:  Key, // Use a generic name to represent the session id.
		Value: id,
		// Path:     "/",
		// Domain:   "",
		MaxAge: int((20 * time.Minute).Seconds()),
		// Secure:   true,
		HttpOnly: true,
	}
}
