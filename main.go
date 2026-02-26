package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/MudassirDev/go-chat/db/database"
	"github.com/MudassirDev/go-chat/internal/web"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	PORT    string
	DB_CONN *sql.DB
	HANDLER *http.ServeMux
	//go:embed db/schema/*.sql
	embedMigrations embed.FS
)

func init() {
	godotenv.Load()
	envs := map[string]string{
		"PORT":       "",
		"JWT_SECRET": "",
		"DB_URL":     "",
	}

	for env := range envs {
		envs[env] = os.Getenv(env)
		validateEnv(envs[env], env)
	}

	conn, err := sql.Open("postgres", envs["DB_URL"])
	if err != nil {
		log.Fatalf("failed to make a connection with DB: %v", err)
	}

	tmpls := template.New("")
	queries := database.New(conn)
	apiCfg := web.APIConfig{
		DB:        queries,
		JwtSecret: envs["JWT_SECRET"],
		Templates: tmpls,
	}

	DB_CONN = conn
	PORT = envs["PORT"]
	HANDLER = web.CreateMux(&apiCfg)
}

func init() {
	goose.SetDialect("postgres")
	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(DB_CONN, "db/schema"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func main() {
	defer DB_CONN.Close()
	srv := http.Server{
		Addr:    ":" + PORT,
		Handler: HANDLER,
	}

	log.Printf("Server is listerning at http://localhost:%v\n", PORT)
	log.Fatal(srv.ListenAndServe())
}

func validateEnv(env, envName string) {
	if env == "" {
		log.Fatal("missing env variable: ", envName)
	}
}
