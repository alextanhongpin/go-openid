package main

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

var (
	ContextKeyClientID     = ContextKey("client_id")
	ContextKeyClientSecret = ContextKey("client_secret")
	ContextKeyTimestamp    = ContextKey("timestamp")
	ContextKeySubject      = ContextKey("subject")
)
