package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"realtime-server/internal/chat"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// gán từ main.go một lần lúc khởi động, để không mở kết nối Mongo mới mỗi lần có client connect
var ChatRepo *chat.Repository

func SetChatRepo(repo *chat.Repository) {
	ChatRepo = repo
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	userID := r.URL.Query().Get("user")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	client := &Client{
		Conn:   conn,
		UserID: userID,
		RoomID: roomID,
		Send:   make(chan []byte),
	}

	WS.Register <- client

	go client.readPump()
	go client.writePump()
}

func (c *Client) readPump() {
	defer func() {
		WS.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var m chat.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		if m.Content == "" {
			continue
		}

		m.SenderID, _ = bson.ObjectIDFromHex(c.UserID)
		m.RoomID = c.RoomID
		ChatRepo.SaveMessage(context.TODO(), &m)

		WS.Broadcast <- MessagePayload{
			RoomID:  c.RoomID,
			Message: msg,
		}
	}
}

func (c *Client) writePump() {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
