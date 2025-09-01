package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MudassirDev/go-chat/db/database"
	"github.com/MudassirDev/go-chat/internal/auth"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleCreateUsers(w http.ResponseWriter, r *http.Request) {
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

	password, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("failed to hash the password: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to hash password!"))
		return
	}

	user, err := DB.CreateUser(context.Background(), database.CreateUserParams{
		Username:  req.Username,
		Password:  password,
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

func handleLogin(w http.ResponseWriter, r *http.Request) {
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

	user, err := DB.GetUserWithUsername(context.Background(), req.Username)
	if err != nil {
		log.Printf("no user found: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("a user with this username doesn't exist!"))
		return
	}

	err = auth.VerifyPassword(req.Password, user.Password)
	if err != nil {
		log.Printf("password doesn't match: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong password!"))
		return
	}

	jwtToken, err := auth.CreateJWT(user.ID, JWT_SECRET, EXPIRY_TIME)
	if err != nil {
		log.Printf("failed to create token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("token creation failed!"))
		return
	}

	cookie := http.Cookie{
		Name:     AUTH_KEY,
		Value:    jwtToken,
		Path:     "/",
		Expires:  time.Now().Add(EXPIRY_TIME),
		MaxAge:   int(EXPIRY_TIME),
		Secure:   false,
		HttpOnly: false,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logged in"))
}

func handlerUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawUser := r.Context().Value(AUTH_KEY)
		if rawUser == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		user, ok := rawUser.(database.User)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		users, err := DB.GetAllUsersExceptCurrent(context.Background(), user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}
		templates.ExecuteTemplate(w, "users.html", users)
	})
}
