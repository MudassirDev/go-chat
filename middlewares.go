package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MudassirDev/go-chat/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AUTH_KEY)
		if err != nil {
			log.Printf("cookie error: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("auth cookie not found"))
			return
		}

		id, err := auth.VerifyJWT(JWT_SECRET, cookie.Value)
		if err != nil {
			log.Printf("failed to verify jwt %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid jwt token"))
			return
		}
		user, err := DB.GetUserWithID(context.Background(), id)
		if err != nil {
			log.Printf("cookie error: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user doesn't exist"))
			return
		}
		ctx := context.WithValue(r.Context(), AUTH_KEY, user)

		request := r.WithContext(ctx)
		next.ServeHTTP(w, request)
	})
}
