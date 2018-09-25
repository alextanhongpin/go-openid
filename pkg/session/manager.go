package session

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

// Manager represents the manager that handles the session creation and cleanup.
type Manager struct {
	sync.Once
	sync.RWMutex
	batch int // The size of the keys to gather during cleanup.
	quit  chan struct{}
	data  map[string]*Session
}

// NewManager returns a new session manager.
func NewManager() *Manager {
	return &Manager{
		Once:  sync.Once{},
		batch: 50,
		quit:  make(chan struct{}),
		data:  make(map[string]*Session),
	}
}

// Get the session data by the given session id.
func (m *Manager) Get(sid string) (*Session, error) {
	m.Lock()
	defer m.Unlock()

	sess, exist := m.data[sid]
	if !exist {
		return nil, errors.New("does not exist")
	}

	// Current time is greater than the expire at time (>= 0).
	if time.Since(sess.ExpireAt) >= 0 {
		// Passive session deletion.
		m.Delete(sid)
		return nil, errors.New("session expired")
	}
	return sess, nil
}

func (m *Manager) Delete(sid string) {
	delete(m.data, sid)
}

// Put stores a new session id and the session object in the storage.
func (m *Manager) Put(sid string, sess *Session) {
	m.Lock()
	defer m.Unlock()
	m.data[sid] = sess
}

func (m *Manager) clean() {
	log.Println("mgr: checking for expired sessions...")
	var expiredSIDs []string
	i := 0
	m.RLock()
	for k, v := range m.data {
		if i >= m.batch {
			break
		}
		if time.Since(v.ExpireAt) >= 0 {
			expiredSIDs = append(expiredSIDs, k)
		}
		i++
	}
	m.RUnlock()

	log.Printf("mgr: found %d expires sessions\n", len(expiredSIDs))
	m.Lock()
	for _, v := range expiredSIDs {
		m.Delete(v)
	}
	m.Unlock()
}

// Stop terminates the running goroutine that is responsible for cleaning up
// the expired sessions.
func (m *Manager) Stop() {
	// Ensure the close is only called once.
	m.Once.Do(func() {
		close(m.quit)
	})
}

// Start will run a goroutine to clear the sessions every minute.
func (m *Manager) Start() {
	go func() {
		c := time.Tick(1 * time.Minute)
		for range c {
			select {
			case <-m.quit:
				return
			default:
				m.clean()
			}
		}
	}()
}

// GetSession retrieves the user session from the request.
func (m *Manager) GetSession(r *http.Request) (*Session, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}
	// If the session does not exist, an error will be thrown.
	return m.Get(c.Value)
}

// SetSession sets a new session in the response.
func (m *Manager) SetSession(w http.ResponseWriter, userID string) {
	s := NewSession(userID)
	m.Put(s.SID, s)
	cookie := NewCookie(s.SID)
	http.SetCookie(w, cookie)
}

// HasSession checks if a session exist.
func (m *Manager) HasSession(r *http.Request) bool {
	sess, err := m.GetSession(r)
	return err == nil && sess != nil
}
