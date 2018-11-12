package session

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// Key represents a general name for the cookie.
const Key = "id"

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
