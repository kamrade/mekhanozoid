package web

import (
	"html/template"
	"net/http"
)

type Handlers struct {
	templates *template.Template
}

func NewHandlers() *Handlers {
	t := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/game.html",
	))

	return &Handlers{
		templates: t,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Mekhanozoid web server is running"))
}

func (h *Handlers) Game(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Game",
	}

	if err := h.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
