package main

import "flag"

type Config struct {
	Port         int
	TemplateDir  string
	ClientID     string
	ClientSecret string
}

func NewConfig() Config {
	var cfg Config
	flag.IntVar(&cfg.Port, "port", 4000, "the port of the application")
	flag.StringVar(&cfg.TemplateDir, "tpldir", "templates", "the directory the html templates resides")
	flag.StringVar(&cfg.ClientID, "client_id", "", "the openid connect client id")
	flag.Parse()
	return cfg
}
