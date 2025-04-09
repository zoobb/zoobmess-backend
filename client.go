package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type ClientMessage struct {
	Type      string `json:"type"`
	ID        int    `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte),
	}
}

func (c *Client) ReadMessages(hub *Hub) {
	defer func() {
		hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var incoming ClientMessage
		err = json.Unmarshal(msg, &incoming)
		if err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		switch incoming.Type {
		case "message":
			dbMsg := Message{
				Username:  incoming.Username,
				Message:   incoming.Message,
				Timestamp: incoming.Timestamp,
			}
			_ = SaveMessage(dbMsg)
			out, _ := json.Marshal(incoming)
			hub.broadcast <- out

		case "update":
			if incoming.ID > 0 && incoming.Message != "" {
				err := UpdateMessage(incoming.ID, incoming.Message)
				if err == nil {
					log.Printf("Updated message ID %d\n", incoming.ID)
				}
			}

		case "delete":
			if incoming.ID > 0 {
				err := DeleteMessage(incoming.ID)
				if err == nil {
					log.Printf("Deleted message ID %d\n", incoming.ID)
				}
			}
		}

		dbMsg := Message{
			Username:  incoming.Username,
			Message:   incoming.Message,
			Timestamp: incoming.Timestamp,
		}
		_ = SaveMessage(dbMsg)

		out, err := json.Marshal(incoming)
		if err != nil {
			continue
		}

		hub.broadcast <- out
		log.Println("Message broadcasted")
	}
}

func (c *Client) WriteMessages() {
	go func() {
		history, err := LoadLastMessages(50)
		if err == nil {
			for _, m := range history {
				out, err := json.Marshal(ClientMessage{
					Username:  m.Username,
					Message:   m.Message,
					Timestamp: m.Timestamp,
				})
				if err == nil {
					c.send <- out
				}
			}
		}
	}()

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(c.conn)

	for msg := range c.send {
		log.Printf("Sending to client: %s\n", string(msg))
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
