package web

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/MudassirDev/go-chat/internal/auth"
)

func (c *APIConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AUTH_KEY)
		if err != nil {
			log.Printf("cookie error: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		id, err := auth.VerifyJWT(c.JwtSecret, cookie.Value)
		if err != nil {
			log.Printf("failed to verify jwt %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user, err := c.DB.GetUserWithID(context.Background(), id)
		if err != nil {
			log.Printf("cookie error: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), AUTH_KEY, user)

		request := r.WithContext(ctx)
		next.ServeHTTP(w, request)
	})
}

func (c *APIConfig) postMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			respondWithError(
				w,
				http.StatusBadRequest,
				"invalid header content type",
				errors.New("invalid header content type"),
			)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("wrong Content-Type"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
