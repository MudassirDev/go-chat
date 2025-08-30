package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MudassirDev/go-chat/db/database"
)

func handleCreateUsers(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong Content-Type"))
		return
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		log.Printf("failed to decode res: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to decode"))
		return
	}

	user, err := DB.CreateUser(context.Background(), database.CreateUserParams{
		Username:  req.Username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Printf("failed to create user: %v", err)
		if strings.Contains(err.Error(), "UNIQUE") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("username already taken!"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to create user"))
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to create json payload: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user created! but failed to respond"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}
