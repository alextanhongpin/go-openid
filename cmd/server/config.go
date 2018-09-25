package main

import "time"

// TODO: Don't use global variable, scope it at the initialization in a Config
// struct.
var (
	defaultIssuer        = "go-openid"
	defaultJWTSigningKey = "secret"
	defaultDuration      = time.Hour
)

// Config represents the global app config.
type Config struct {
	AccessTokenKey []byte
}

// NewConfig returns a new config.
func NewConfig() *Config {
	return &Config{}
}
