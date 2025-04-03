package room

import (
    "log"

    "github.com/gorilla/websocket"
)

func handleMessages(room *Room) {
    for {
        msg := <-room.Broadcast
        room.Mutex.Lock()
        for client := range room.Clients {
            client.WriteMessage(websocket.TextMessage, msg)
        }
        room.Mutex.Unlock()
    }
}

func handleClientMessages(room *Room, conn *websocket.Conn) {
    defer func() {
        room.Mutex.Lock()
        delete(room.Clients, conn)
        room.Mutex.Unlock()
        conn.Close()
        log.Printf("Client disconnected from room: %s", room.ID)
    }()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }
        room.Broadcast <- msg
    }
}