package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	PORT     string
	upgrader websocket.Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func init() {
	godotenv.Load()

	port := os.Getenv("PORT")
	validateEnv(port, "PORT")
	PORT = port
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("failed to make a connection", err)
		w.Write([]byte("failed"))
	}
	defer conn.Close()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/templates/index.html")
	})
	mux.HandleFunc("/ws", handleWS)
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

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
