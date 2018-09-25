package session

import (
	"errors"
	"log"
	"sync"
	"time"
)

// Repository represents the interface for the session storage operations.
type Repository interface {
	Close()
	Delete(id string) error
	Get(id string) (*Session, error)
	Open()
	Put(id string, s *Session) error
}

type repositoryInMemoryImpl struct {
	batch int // The size of the keys to gather during cleanup.
	data  map[string]*Session
	quit  chan struct{}
	sync.Once
	sync.RWMutex
}

// NewInMemoryRepository returns a new in-memory storage that conforms to the
// repository interface.
func NewInMemoryRepository() *repositoryInMemoryImpl {
	return &repositoryInMemoryImpl{
		Once:  sync.Once{},
		batch: 50,
		data:  make(map[string]*Session),
		quit:  make(chan struct{}),
	}
}

func (r *repositoryInMemoryImpl) clean() {
	log.Println("mgr: checking for expired sessions...")
	var expiredSIDs []string
	i := 0
	r.RLock()
	for k, v := range r.data {
		if i >= r.batch {
			break
		}
		if time.Since(v.ExpireAt) >= 0 {
			expiredSIDs = append(expiredSIDs, k)
		}
		i++
	}
	r.RUnlock()

	log.Printf("mgr: found %d expires sessions\n", len(expiredSIDs))
	r.Lock()
	for _, v := range expiredSIDs {
		r.Delete(v)
	}
	r.Unlock()
}

func (r *repositoryInMemoryImpl) worker() {
	c := time.Tick(1 * time.Minute)
	for range c {
		select {
		case <-r.quit:
			return
		default:
			r.clean()
		}
	}
}

// Open initialize the storage.
func (r *repositoryInMemoryImpl) Open() {
	go r.worker()
}

// Close terminates the storage safely.
func (r *repositoryInMemoryImpl) Close() {
	r.Once.Do(func() {
		close(r.quit)
	})
}

// Get returns the session data by the given session id.
func (r *repositoryInMemoryImpl) Get(id string) (*Session, error) {
	r.RLock()
	sess, exist := r.data[id]
	r.RUnlock()
	if !exist {
		return nil, errors.New("does not exist")
	}

	// Current time is greater than the expire at time (>= 0).
	if time.Since(sess.ExpireAt) >= 0 {
		// Passive session deletion.
		r.Delete(id)
		return nil, errors.New("session expired")
	}
	return sess, nil
}

// Delete remove the session from the storage.
func (r *repositoryInMemoryImpl) Delete(id string) error {
	r.Lock()
	delete(r.data, id)
	r.Unlock()

	return nil
}

// Put stores a new session id and the session object in the storage.
func (r *repositoryInMemoryImpl) Put(id string, sess *Session) error {
	r.Lock()
	r.data[id] = sess
	r.Unlock()

	return nil
}
