package web

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/MudassirDev/go-chat/db/database"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	TEXT_MESSAGE  string = "TEXT"
	AUDIO_MESSAGE string = "AUDIO"
)

var connections map[uuid.UUID]*websocket.Conn = make(map[uuid.UUID]*websocket.Conn)
var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Recipient   uuid.UUID `json:"recipient_id"`
	MessageType string    `json:"message_type"`
	Content     string    `json:"content,omitempty"`
	ContentData []byte    `json:"content_data,omitempty"`
	Time        time.Time `json:"time"`
}

func (c *APIConfig) handlerWS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromContext(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to create a connection", err)
			return
		}

		var msg Message
		connections[user.ID] = conn
		defer conn.Close()
		for {
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				delete(connections, user.ID)
				return
			}

			if msg.MessageType != AUDIO_MESSAGE && msg.MessageType != TEXT_MESSAGE {
				continue
			}
			if msg.MessageType == AUDIO_MESSAGE {
				filePath, err := c.saveAudio(msg.ContentData)
				if err != nil {
					log.Println(err)
					continue
				}

				msg.Content = filePath
			}

			message, err := c.DB.CreateMessage(context.Background(), database.CreateMessageParams{
				SenderID:    user.ID,
				RecipientID: msg.Recipient,
				Time:        msg.Time,
				Content:     msg.Content,
				MessageType: msg.MessageType,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
			if err != nil {
				conn.WriteJSON(struct {
					Content string `json:"content"`
				}{
					Content: "failed to send message",
				})
				log.Println(err)
				return
			}

			err = conn.WriteJSON(message)
			if err != nil {
				log.Printf("error while writing: %v", err)
			}

			receiverConn, ok := connections[msg.Recipient]
			if !ok {
				continue
			}

			receiverConn.WriteJSON(message)
		}
	})
}

func (c *APIConfig) handlerMessagesTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Recipient database.GetUserWithIDRow `json:"recipient"`
			Messages  []database.Message        `json:"messages"`
		}

		user, err := getUserFromContext(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
			return
		}

		raw_recipient_id := r.PathValue("userid")
		recipient_id, err := uuid.Parse(raw_recipient_id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid id", err)
			return
		}

		if user.ID == recipient_id {
			msg := "cannot initiate chat with yourself"
			respondWithError(
				w,
				http.StatusBadRequest,
				msg,
				errors.New(msg),
			)
			return
		}

		recipient, err := c.DB.GetUserWithID(context.Background(), recipient_id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "no such user", err)
			return
		}

		messages, err := c.DB.GetChatMessages(context.Background(), database.GetChatMessagesParams{
			RecipientID: recipient.ID,
			SenderID:    user.ID,
		})
		if err != nil {
			log.Println(err)
			c.Templates.ExecuteTemplate(w, "messages", Response{
				Recipient: recipient,
			})
			return
		}
		c.Templates.ExecuteTemplate(w, "messages", Response{
			Recipient: recipient,
			Messages:  messages,
		})
	})
}

func (c *APIConfig) handlerFiles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromContext(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
			return
		}

		filename := r.PathValue("filename")
		filepath := path.Join("files", filename)
		_, err = c.DB.GetMessageWithFileName(context.Background(), database.GetMessageWithFileNameParams{
			Content:     filepath,
			RecipientID: user.ID,
		})

		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, filepath)
	})
}
