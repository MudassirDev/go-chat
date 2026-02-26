package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/MudassirDev/go-chat/db/database"
	"github.com/MudassirDev/go-chat/internal/auth"
	"github.com/google/uuid"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *APIConfig) handlerRegister() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Request
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid payload", err)
			return
		}

		password, err := auth.HashPassword(req.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to create user", err)
			return
		}

		user, err := c.DB.CreateUser(context.Background(), database.CreateUserParams{
			Username:  req.Username,
			Password:  password,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				respondWithError(w, http.StatusBadRequest, "username already taken", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, "failed to create user", err)
			return
		}

		respondWithJSON(w, http.StatusCreated, user)
	})
}

func (c *APIConfig) handlerLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Request
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid payload", err)
			return
		}

		user, err := c.DB.GetUserWithUsername(context.Background(), req.Username)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "user doesnt exist", err)
			return
		}

		err = auth.VerifyPassword(req.Password, user.Password)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "wrong password", err)
			return
		}

		jwtToken, err := auth.CreateJWT(user.ID, c.JwtSecret, EXPIRY_TIME)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to create token", err)
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
		respondWithJSON(w, http.StatusOK, "logged in")
	})
}

func (c *APIConfig) handlerChatTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromContext(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
			return
		}

		users, err := c.DB.GetAllUsersExceptCurrent(context.Background(), user.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to get users", err)
			return
		}

		c.Templates.ExecuteTemplate(w, "chat", struct {
			Users    []database.GetAllUsersExceptCurrentRow
			SenderID uuid.UUID
			Username string
		}{
			Users:    users,
			SenderID: user.ID,
			Username: user.Username,
		})
	})
}
