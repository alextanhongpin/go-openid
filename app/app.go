package app

import (
	"github.com/julienschmidt/httprouter"
)

// Env is the environment for the application and holds references to the database, router and template context.
// Env needs to be initialized only once in the application lifecyle, and is passed down through dependency injection.
type Env struct {
	// The database we are connecting to
	DB *Database
	// The router for the application
	Router *httprouter.Router
	// The cached templates
	Tmpl *Template
	// The cache mechanism that we are using
	Cache *Cache
	// Opentracing implementation
	Tracer *Tracer
}
