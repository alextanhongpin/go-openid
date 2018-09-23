package main

import "time"

// TODO: Don't use global variable, scope it at the initialization in a Config
// struct.
var (
	defaultIssuer        = "go-openid"
	defaultJWTSigningKey = "secret"
	defaultDuration      = time.Hour
)

type Config struct {
	AccessTokenKey []byte
}

func NewConfig() *Config {
	return &Config{}
}
