package main

import (
	"github.com/alextanhongpin/go-openid/pkg/html5"
	"github.com/alextanhongpin/go-openid/pkg/session"
)

// Endpoints represent the endpoints for the OpenIDConnect.
type Endpoints struct {
	service        *serviceImpl
	sessionManager *session.Manager
	template       *html5.Template
}
