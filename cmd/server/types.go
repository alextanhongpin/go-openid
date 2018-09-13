package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

// HTMLs represent the html templates stored as dictionary.
type HTMLs struct {
	datadir   string
	templates map[string]*template.Template
	sync.Once
}

// NewHTMLs returns a new HTMLs struct.
func NewHTMLs(datadir string) *HTMLs {
	return &HTMLs{
		datadir:   datadir,
		templates: make(map[string]*template.Template),
	}
}

// Render renders the html output with the given data.
func (h HTMLs) Render(w http.ResponseWriter, name string, data interface{}) {
	if t, ok := h.templates[name]; !ok {
		err := fmt.Sprintf("template with the name %s does not exist", name)
		http.Error(w, err, http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t.Execute(w, data)
	}
}

func (h *HTMLs) path(f string) string {
	return fmt.Sprintf("%s/%s.tmpl", h.datadir, f)
}

// Load loads the file from the base path.
func (h *HTMLs) Load(files ...string) {
	h.Do(func() {
		layout := template.Must(template.New("base").ParseFiles(h.path("base")))
		for _, f := range files {
			clone := template.Must(layout.Clone())
			h.templates[f] = template.Must(clone.ParseFiles(h.path(f)))
		}
	})
}
