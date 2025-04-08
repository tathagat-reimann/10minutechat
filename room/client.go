package room

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func sendNewMessageToAllClients(room *Room) {
	for {
		message := <-room.Broadcast
		room.Mutex.Lock()
		for clientWebSocketConnection, clientName := range room.Clients {
			err := clientWebSocketConnection.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending message to %s: %v", clientName, err)
				clientWebSocketConnection.Close()
				delete(room.Clients, clientWebSocketConnection)
			}
		}
		room.Mutex.Unlock()
	}
}

func handleNewMessageFromClient(room *Room, conn *websocket.Conn, clientName string) {
	defer func() {
		room.Mutex.Lock()
		delete(room.Clients, conn)
		room.Mutex.Unlock()
		conn.Close()
		message := Message{
			Type:    "info",
			Sender:  "System",
			Content: clientName + " has left the room.",
		}
		room.Broadcast <- message
		log.Printf("Client %s disconnected from room: %s", clientName, room.ID)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", clientName, err)
			break
		}

		// Parse the incoming message
		var incomingMessage string
		if err := json.Unmarshal(msg, &incomingMessage); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Handle regular chat messages
		message := Message{
			//Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Type:    "chat",
			Sender:  clientName,
			Content: incomingMessage,
		}

		room.Broadcast <- message
	}
}

func sendClientNameToItself(room *Room, conn *websocket.Conn, clientName string) {
	room.Mutex.Lock()
	message := Message{
		Type:    "clientName",
		Sender:  "System",
		Content: clientName,
	}

	if err := conn.WriteJSON(message); err != nil {
		log.Printf("Error sending client name to client: %v", err)
		conn.Close()
		delete(room.Clients, conn)
	}
	room.Mutex.Unlock()
}

func sendClientNameToOtherClients(room *Room, conn *websocket.Conn, clientName string) {
	message := Message{
		Type:    "info",
		Sender:  "System",
		Content: clientName + " has joined the room.",
	}

	room.Mutex.Lock()
	for clientWebSocketConnection, clientName := range room.Clients {
		if clientWebSocketConnection == conn {
			continue // Skip sending the message to the client itself
		}
		err := clientWebSocketConnection.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to %s: %v", clientName, err)
			clientWebSocketConnection.Close()
			delete(room.Clients, clientWebSocketConnection)
		}
	}
	room.Mutex.Unlock()
}
