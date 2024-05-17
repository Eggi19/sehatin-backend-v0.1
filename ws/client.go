package ws

import (
	"encoding/json"
	"log"

	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	Id       string `json:"id"`
	UserId   int64  `json:"user_id,omitempty"`
	UserRole string `json:"user_role"`
	RoomId   string `json:"roomId"`
}

type Message struct {
	Content  string `json:"content"`
	RoomId   string `json:"roomId"`
	Id       string `json:"id"`
	UserId   int64  `json:"user_id,omitempty"`
	UserRole string `json:"user_role,omitempty"`
	Type     string `json:"type"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var data dtos.SentMessage

		err = json.Unmarshal([]byte(string(m)), &data)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		msg := &Message{
			Content:  data.Content,
			RoomId:   c.RoomId,
			Id:       c.Id,
			UserId:   c.UserId,
			UserRole: c.UserRole,
			Type:     data.Type,
		}

		hub.Broadcast <- msg
	}
}
