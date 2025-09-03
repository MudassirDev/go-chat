package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MudassirDev/go-chat/db/database"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Recipient string    `json:"recipient_id"`
	Sender    string    `json:"sender_id"`
	Content   string    `json:"content"`
	Time      time.Time `json:"time"`
}

func handleWS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("failed to make a connection %v\n", err)
			w.Write([]byte("failed"))
			return
		}

		var msg Message
		defer conn.Close()

		for {
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	})
}

func handlerChat() http.Handler {
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
		raw_recipient_id := r.PathValue("userid")
		recipient_id, err := strconv.Atoi(raw_recipient_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to check id"))
			return
		}
		if user.ID == int64(recipient_id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("cannot initiate a chat with yourself"))
			return
		}
		recipient, err := DB.GetUserWithID(context.Background(), int64(recipient_id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no such user!"))
			return
		}
		templates.ExecuteTemplate(w, "chat", recipient)
	})
}
