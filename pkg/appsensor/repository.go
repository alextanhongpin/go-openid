package appsensor

import (
	"sync"
	"time"
)

// Attempt represent the attempts made by the attacker.
type Attempt struct {
	ID string
	// Indicates the curent attempt of the id.
	Count int64
	// Indicates how frequent has this id been logged before.
	Frequency   int64
	CreatedAt   time.Time
	AttemptedAt time.Time
	UpdatedAt   time.Time
}

type repoInMemoryImpl struct {
	sync.RWMutex
	data map[string]*Attempt
}

func NewInMemoryImpl() *repoInMemoryImpl {
	return &repoInMemoryImpl{
		data: make(map[string]*Attempt),
	}
}

func (r *repoInMemoryImpl) Put(id string) {
	r.Lock()
	defer r.Unlock()

	now := time.Now().UTC()
	attempt := &Attempt{
		ID:          id,
		Count:       1,
		CreatedAt:   now,
		AttemptedAt: now,
		UpdatedAt:   now,
	}
	r.data[id] = attempt
}

func (r *repoInMemoryImpl) Get(id string) (*Attempt, bool) {
	r.RLock()
	attempt, exist := r.data[id]
	r.RUnlock()
	return attempt, exist
}

func (r *repoInMemoryImpl) Reset(id string) {
	r.Lock()
	defer r.Unlock()

	attempt, exist := r.data[id]
	if !exist {
		return
	}
	// Reset the count.
	attempt.Count = 0

	// But increase the frequency.
	attempt.Frequency++

	// Update the last updated date.
	attempt.UpdatedAt = time.Now().UTC()
}
