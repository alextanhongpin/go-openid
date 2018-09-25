package session

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"
)

var cookieName = "id"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Session represents the information that is tracked by the session.
type Session struct {
	UserID    string
	SID       string
	CreatedAt time.Time
	ExpireAt  time.Time
	IP        string
	UserAgent string
}

// NewSession returns a new session and cookie.
func NewSession(userID string) *Session {
	now := time.Now().UTC()
	sess := &Session{
		UserID:    userID,
		SID:       NewSessionID(now),
		CreatedAt: now,
		ExpireAt:  now.Add(20 * time.Minute),
		IP:        "",
		UserAgent: "",
	}
	return sess
}

// NewSessionID creates a new session id from the given time.
func NewSessionID(t time.Time) string {
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

// NewCookie returns a new cookie with default settings.
func NewCookie(sid string) *http.Cookie {
	return &http.Cookie{
		Name:  cookieName, // Use a generic name to represent the session id.
		Value: sid,
		// Path:     "/",
		// Domain:   "",
		MaxAge: int((20 * time.Minute).Seconds()),
		// Secure:   true,
		HttpOnly: true,
	}
}
