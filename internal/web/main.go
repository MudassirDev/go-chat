package web

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/MudassirDev/go-chat/db/database"
)

const (
	EXPIRY_TIME time.Duration = time.Hour * 1
	AUTH_KEY    string        = "auth_key"
)

var (
	IS_DEVELOPMENT bool
)

func CreateMux(jwtSecret string, DB *database.Queries, isDevelopment bool) *http.ServeMux {
	IS_DEVELOPMENT = isDevelopment
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))

	tmpls := template.New("")
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory %v", err)
	}

	apiCfg := apiConfig{
		db:        DB,
		jwtSecret: jwtSecret,
		templates: tmpls,
		cwd:       cwd,
	}
	apiCfg.setupTemplate()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			apiCfg.templates.ExecuteTemplate(w, "index.html", nil)
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
	mux.Handle("GET /chat", apiCfg.AuthMiddleware(apiCfg.handlerUsers()))
	mux.Handle("GET /users/{userid}", apiCfg.AuthMiddleware(apiCfg.handlerChat()))
	mux.Handle("GET /files/{filename}", apiCfg.AuthMiddleware(apiCfg.handleFiles()))

	mux.Handle("/chat/{userid}", apiCfg.AuthMiddleware(apiCfg.handleWS()))
	mux.HandleFunc("POST /api/users/create", apiCfg.handleCreateUsers)
	mux.HandleFunc("POST /api/users/login", apiCfg.handleLogin)
	mux.HandleFunc("GET /login", apiCfg.handleLoginTemplate)
	mux.HandleFunc("GET /register", apiCfg.handleRegisterTemplate)

	return mux
}

func (c *apiConfig) setupTemplate() {
	path := path.Join(c.cwd, "static/templates")
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			_, err := c.templates.ParseFiles(path)
			if err != nil {
				log.Fatalf("failed to parse templates: %v", err)
			}
		}
		return nil
	})
}
