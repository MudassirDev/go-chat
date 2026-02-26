package web

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	EXPIRY_TIME time.Duration = time.Hour * 1
	AUTH_KEY    string        = "auth_key"
)

func CreateMux(c *APIConfig) *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))

	c.parseTemplates()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			c.Templates.ExecuteTemplate(w, "index.html", nil)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/") {
			splitResult := strings.Split(r.URL.Path, "")
			newPath := strings.Join(splitResult[:len(splitResult)-1], "")
			http.Redirect(w, r, newPath, http.StatusSeeOther)
			return
		}
		http.NotFound(w, r)
	})
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.Handle("GET /chat", c.authMiddleware(
		c.handlerChatTemplate(),
	))
	mux.Handle("GET /users/{userid}", c.authMiddleware(
		c.handlerMessagesTemplate(),
	))
	mux.Handle("GET /files/{filename}", c.authMiddleware(
		c.handlerFiles(),
	))
	mux.Handle("/chat/ws", c.authMiddleware(
		c.handlerWS(),
	))
	mux.Handle("POST /api/users/create", c.postMiddleware(
		c.handlerRegister(),
	))
	mux.Handle("POST /api/users/login", c.postMiddleware(
		c.handlerLogin(),
	))
	mux.HandleFunc("GET /login", c.handlerLoginTemplate)
	mux.HandleFunc("GET /register", c.handlerRegisterTemplate)

	return mux
}

func (c *APIConfig) parseTemplates() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}
	path := path.Join(cwd, "static/templates")
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			_, err := c.Templates.ParseFiles(path)
			if err != nil {
				log.Fatalf("failed to parse templates: %v", err)
			}
		}
		return nil
	})
}
