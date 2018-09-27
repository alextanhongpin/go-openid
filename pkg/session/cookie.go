package session

import (
	"net/http"
	"time"
)

// NewCookie returns a new cookie with default settings.
func NewCookie(id string, now time.Time) *http.Cookie {
	duration := 20 * time.Minute
	return &http.Cookie{
		Expires:  now.Add(duration),
		HttpOnly: true,
		MaxAge:   int((duration).Seconds()),
		Name:     Key, // Use a generic name to represent the session id.
		Path:     "/",
		Value:    id,
		// Domain:   "",
		// Secure:   true,
	}
}
