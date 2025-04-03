package room

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tathagat/10minutechat/websocketX"
)

type Room struct {
	ID        string
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mutex     sync.Mutex
}

var (
	rooms   = make(map[string]*Room)
	roomsMu sync.Mutex
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := uuid.New().String()
	roomsMu.Lock()
	defer roomsMu.Unlock()

	rooms[roomID] = &Room{
		ID:        roomID,
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}

	go handleMessages(rooms[roomID])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"room_id": roomID}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	log.Printf("Room created: %s", roomID)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	log.Printf("Request to join room: %s", roomID)
	roomsMu.Lock()
	room, exists := rooms[roomID]
	roomsMu.Unlock()

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if len(room.Clients) >= websocketX.MaxRoomCapacity() {
		http.Error(w, "Room is full", http.StatusForbidden)
		return
	}

	conn, err := websocketX.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading WebSocket connection: %v", err)
		http.Error(w, "Failed to establish WebSocket connection", http.StatusInternalServerError)
		return
	}

	room.Mutex.Lock()
	room.Clients[conn] = true
	room.Mutex.Unlock()

	log.Printf("Client joined room: %s", roomID)
	go handleClientMessages(room, conn)
}
