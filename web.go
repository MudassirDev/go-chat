package main

import (
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	templates = template.New("")
)

func CreateMux() *http.ServeMux {
	setupTemplate()
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		templates.ExecuteTemplate(w, "index.html", nil)
	})
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.Handle("/users", AuthMiddleware(handlerUsers()))
	mux.Handle("GET /users/{userid}", AuthMiddleware(handlerChat()))
	mux.Handle("GET /files/{filename}", AuthMiddleware(handleFiles()))

	mux.Handle("/chat/{userid}", AuthMiddleware(handleWS()))
	mux.HandleFunc("POST /api/users/create", handleCreateUsers)
	mux.HandleFunc("POST /api/users/login", handleLogin)

	return mux
}

func setupTemplate() {
	filepath.WalkDir("static/templates", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			_, err := templates.ParseFiles(path)
			if err != nil {
				log.Fatalf("failed to parse templates: %v", err)
			}
		}
		return nil
	})
}
