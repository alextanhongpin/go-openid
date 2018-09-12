package main

import (
	"fmt"
	"html/template"

	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

func initEndpoints(key string) *Endpoints {
	mem := database.NewInMem()
	cry := crypto.New(key)
	svc := NewService(mem, cry)
	return NewEndpoints(svc)
}

func initTemplates(t HTMLs, files ...string) func() {
	load := func(f string) string {
		return fmt.Sprintf("templates/%s.tmpl", f)
	}
	return func() {
		layout := template.Must(template.New("base").ParseFiles(load("base")))
		for _, f := range files {
			clone := template.Must(layout.Clone())
			t[f] = template.Must(clone.ParseFiles(load(f)))
		}
	}
}
