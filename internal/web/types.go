package web

import (
	"text/template"

	"github.com/MudassirDev/go-chat/db/database"
)

type apiConfig struct {
	db        *database.Queries
	jwtSecret string
	cwd       string
	templates *template.Template
}
