package room

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func handleMessages(room *Room) {
	for {
		message := <-room.Broadcast
		room.Mutex.Lock()
		for client, name := range room.Clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending message to %s: %v", name, err)
				client.Close()
				delete(room.Clients, client)
			}
		}
		room.Mutex.Unlock()
	}
}

func handleClientMessages(room *Room, conn *websocket.Conn, clientName string) {
	defer func() {
		room.Mutex.Lock()
		delete(room.Clients, conn)
		room.Mutex.Unlock()
		conn.Close()
		log.Printf("Client %s disconnected from room: %s", clientName, room.ID)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", clientName, err)
			break
		}

		// Parse the incoming message
		var incomingMessage map[string]string
		if err := json.Unmarshal(msg, &incomingMessage); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Handle "getName" message type
		if incomingMessage["type"] == "getName" {
			response := map[string]string{
				"type": "name",
				"name": clientName,
			}
			conn.WriteJSON(response)
			continue
		}

		// Handle regular chat messages
		message := Message{
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Sender:    clientName,
			Content:   incomingMessage["content"],
		}

		room.Broadcast <- message
	}
}
