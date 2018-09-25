package session

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"
)

// Key represents a general name for the cookie.
const Key = "id"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Session represents the information that is tracked by the session.
type Session struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  time.Time
	IP        string
	SessionID string
	UserAgent string
	UserID    string
}

// NewSession returns a new session and cookie.
func NewSession(userID string) *Session {
	now := time.Now().UTC()
	sess := &Session{
		CreatedAt: now,
		ExpireAt:  now.Add(20 * time.Minute),
		IP:        "",
		SessionID: NewSessionID(now),
		UserAgent: "",
		UserID:    userID,
	}
	return sess
}

// NewSessionID creates a new session id from the given time.
func NewSessionID(t time.Time) string {
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

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
