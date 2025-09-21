package web

import (
	"net/http"
	"text/template"
	"time"

	"github.com/MudassirDev/go-chat/db/database"
)

const (
	EXPIRY_TIME time.Duration = time.Hour * 1
	AUTH_KEY    string        = "auth_key"
)

func CreateMux(jwtSecret string, DB *database.Queries, templates *template.Template) *http.ServeMux {
	mux := http.NewServeMux()

	_ = apiConfig{
		db:        DB,
		jwtSecret: jwtSecret,
		templates: templates,
	}

	return mux
}
