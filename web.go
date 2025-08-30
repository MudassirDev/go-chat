package main

import (
	"net/http"
)

func CreateMux() *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/templates/index.html")
	})
	mux.HandleFunc("/ws", handleWS)
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/api/users/create", handleCreateUsers)

	return mux
}
