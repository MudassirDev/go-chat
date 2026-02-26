package web

import (
	"text/template"

	"github.com/MudassirDev/go-chat/db/database"
)

type APIConfig struct {
	DB        *database.Queries
	JwtSecret string
	Templates *template.Template
}
