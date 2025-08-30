package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Recipient string    `json:"recipient_username"`
	Sender    string    `json:"sender_username"`
	Content   string    `json:"content"`
	Time      time.Time `json:"time"`
}

func handleWS(w http.ResponseWriter, r *http.Request) {
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
}
