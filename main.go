package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/MudassirDev/go-chat/db/database"
	"github.com/MudassirDev/go-chat/internal/web"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	PORT    string
	DB_CONN *sql.DB
	//go:embed db/schema/*.sql
	embedMigrations embed.FS
	HANDLER         *http.ServeMux
)

func init() {
	godotenv.Load()

	log.Println("loading env variables")

	port := os.Getenv("PORT")
	validateEnv(port, "PORT")
	PORT = port

	jwtSecret := os.Getenv("JWT_SECRET")
	validateEnv(jwtSecret, "JWT_SECRET")

	dbURL := os.Getenv("DB_URL")
	validateEnv(dbURL, "DB_URL")

	log.Println("env variables loaded")

	log.Println("making a connection with DB")

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to make a connection with DB: %v", err)
	}
	DB_CONN = conn

	log.Println("DB connection formed!")

	queries := database.New(conn)
	handler := web.CreateMux(jwtSecret, queries)
	HANDLER = handler
}

func init() {
	log.Println("running migrations")

	goose.SetDialect("postgres")
	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(DB_CONN, "db/schema"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("migrations ran successfully")
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
