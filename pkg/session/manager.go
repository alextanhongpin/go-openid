package session

import (
	"net/http"
	"time"
)

// Manager represents the manager that handles the session creation and
// cleanup.
type Manager struct {
	repo Repository
}

// NewManager returns a new session manager.
func NewManager() *Manager {
	return &Manager{
		repo: NewInMemoryRepository(),
	}
}

// Stop terminates the running goroutine that is responsible for cleaning up
// the expired sessions.
func (m *Manager) Stop() {
	m.repo.Close()
}

// Start will run a goroutine to clear the sessions every minute.
func (m *Manager) Start() {
	m.repo.Open()
}

// GetSession retrieves the user session from the request.
func (m *Manager) GetSession(r *http.Request) (*Session, error) {
	c, err := r.Cookie(Key)
	if err != nil {
		return nil, err
	}
	// If the session does not exist, an error will be thrown.
	return m.repo.Get(c.Value)
}

// SetSession sets a new session in the response.
func (m *Manager) SetSession(w http.ResponseWriter, userID string) {
	s := NewSession(userID)
	c := NewCookie(s.SessionID, time.Now().UTC())

	m.repo.Put(s.SessionID, s)

	http.SetCookie(w, c)
}

// Delete removes a session from the session store.
func (m *Manager) Delete(sessionID string) error {
	return m.repo.Delete(sessionID)
}

// HasSession checks if a session exist.
func (m *Manager) HasSession(r *http.Request) bool {
	sess, err := m.GetSession(r)
	return err == nil && sess != nil
}
