package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT string
)

func init() {
	godotenv.Load()

	port := os.Getenv("PORT")
	validateEnv(port, "PORT")
	PORT = port
}

func main() {

	mux := http.NewServeMux()

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
