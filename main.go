package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var hub = NewHub()
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(w, "Could not upgrade connection", err)
		return
	}

	client := NewClient(conn)
	hub.register <- client

	go client.ReadMessages(hub)
	go client.WriteMessages()
}

func main() {
	InitDB()
	go hub.Run()

	http.HandleFunc("/chat", handleWebSocket)

	fmt.Println("Server started on :8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}
