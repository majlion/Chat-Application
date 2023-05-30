package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register new client
	clients[ws] = true

	for {
		var message Message
		// Read message from client
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Broadcast the received message to all connected clients
		broadcast <- message
	}
}

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

var (
	clients   = make(map[*websocket.Conn]bool) // Connected clients
	broadcast = make(chan Message)             // Message channel
)

func handleMessages() {
	for {
		// Receive message from broadcast channel
		message := <-broadcast

		// Send message to all connected clients
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
			fmt.Printf("Message sent: %+v\n", message)
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("Server started on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
