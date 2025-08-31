package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MudassirDev/go-chat/db/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

var (
	PORT    string
	DB_CONN *sql.DB
	//go:embed db/schema/*.sql
	embedMigrations embed.FS
	DB              *database.Queries
	JWT_SECRET      string
)

const (
	DB_PATH     string        = "app.db"
	EXPIRY_TIME time.Duration = time.Hour * 1
	AUTH_KEY    string        = "auth_key"
)

func init() {
	godotenv.Load()

	log.Println("loading env variables")

	port := os.Getenv("PORT")
	validateEnv(port, "PORT")
	PORT = port

	jwtSecret := os.Getenv("JWT_SECRET")
	validateEnv(jwtSecret, "JWT_SECRET")
	JWT_SECRET = jwtSecret

	log.Println("env variables loaded")

	log.Println("making a connection with DB")

	conn, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Fatalf("failed to make a connection with DB: %v", err)
	}
	DB_CONN = conn

	log.Println("DB connection formed!")

	log.Println("creating new DB Queries")

	queries := database.New(conn)
	DB = queries

	log.Println("new DB Query created!")
}

func init() {
	log.Println("running migrations")

	goose.SetDialect("sqlite3")
	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(DB_CONN, "db/schema"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("migrations ran successfully")
}

func main() {
	defer DB_CONN.Close()
	mux := CreateMux()

	srv := http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Server is listerning at http://localhost:%v\n", PORT)
	log.Fatal(srv.ListenAndServe())
}

func validateEnv(env, envName string) {
	if env == "" {
		log.Fatal("missing env variable: ", envName)
	}
}
