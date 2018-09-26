package appsensor

import (
	"log"
	"time"
)

type LoginDetector interface {
	Stat(id string) *Attempt
	IsLocked(id string) bool
	Increment(id string)
}

type loginDetector struct {
	repository   *repoInMemoryImpl
	threshold    int64
	lockDuration time.Duration
}

func NewLoginDetector() *loginDetector {
	return &loginDetector{
		repository:   NewInMemoryImpl(),
		threshold:    3,
		lockDuration: 15 * time.Minute,
	}
}

func (l *loginDetector) Stat(id string) *Attempt {
	attempt, _ := l.repository.Get(id)
	return attempt
}
func (l *loginDetector) IsLocked(id string) bool {
	attempt, exist := l.repository.Get(id)
	if !exist {
		return false
	}
	if attempt.Count < l.threshold {
		return false
	}

	// Check if it's possible to unlock the account in the next turn.
	if time.Since(attempt.AttemptedAt) > l.lockDuration {
		l.repository.Reset(id)
		return false
	}

	elapsed := l.lockDuration - time.Since(attempt.AttemptedAt)
	log.Printf("loginDetector: %s will be unlocked after %v\n", id, elapsed)

	return true
}

func (l *loginDetector) Increment(id string) {
	attempt, exist := l.repository.Get(id)
	if !exist {
		l.repository.Put(id)
		return
	}
	attempt.Count++
	attempt.AttemptedAt = time.Now().UTC()
	attempt.UpdatedAt = time.Now().UTC()
}
