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

var connections map[int64]*websocket.Conn = make(map[int64]*websocket.Conn)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Recipient   int64     `json:"recipient_id"`
	MessageType string    `json:"message_type"`
	Content     string    `json:"content,omitempty"`
	ContentData []byte    `json:"content_data,omitempty"`
	Time        time.Time `json:"time"`
}

func handleWS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawUser := r.Context().Value(AUTH_KEY)
		if rawUser == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		user, ok := rawUser.(database.GetUserWithIDRow)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("failed to make a connection %v\n", err)
			w.Write([]byte("failed"))
			return
		}

		var msg Message
		connections[user.ID] = conn
		defer conn.Close()
		log.Println(connections)

		for {
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				delete(connections, user.ID)
				return
			}
			log.Println(msg.MessageType)

			// message, err := DB.CreateMessage(context.Background(), database.CreateMessageParams{
			// 	SenderID:    user.ID,
			// 	RecipientID: msg.Recipient,
			// 	Time:        msg.Time,
			// 	Content:     msg.Content,
			// 	MessageType: "TEXT",
			// 	CreatedAt:   time.Now(),
			// 	UpdatedAt:   time.Now(),
			// })
			// if err != nil {
			// 	conn.WriteJSON(struct {
			// 		Content string `json:"content"`
			// 	}{
			// 		Content: "failed to send message",
			// 	})
			// 	log.Println(err)
			// 	return
			// }
			// err = conn.WriteJSON(message)
			// if err != nil {
			// 	log.Printf("error while writing: %v", err)
			// }
			//
			// receiverConn, ok := connections[msg.Recipient]
			// if !ok {
			// 	continue
			// }
			// receiverConn.WriteJSON(msg)
		}
	})
}

func handlerChat() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawUser := r.Context().Value(AUTH_KEY)
		if rawUser == nil {
			log.Println("no user in the request context")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		user, ok := rawUser.(database.GetUserWithIDRow)
		if !ok {
			log.Println("user is not valid")
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
		messages, err := DB.GetChatMessages(context.Background(), database.GetChatMessagesParams{
			RecipientID:   recipient.ID,
			SenderID:      user.ID,
			SenderID_2:    recipient.ID,
			RecipientID_2: user.ID,
		})
		if err != nil {
			log.Println(err)
			templates.ExecuteTemplate(w, "chat", struct {
				Recipient database.GetUserWithIDRow `json:"recipient"`
				Messages  []database.Message        `json:"messages"`
			}{
				Recipient: recipient,
			})
			return
		}
		templates.ExecuteTemplate(w, "chat", struct {
			Recipient database.GetUserWithIDRow `json:"recipient"`
			Messages  []database.Message        `json:"messages"`
		}{
			Recipient: recipient,
			Messages:  messages,
		})
	})
}
