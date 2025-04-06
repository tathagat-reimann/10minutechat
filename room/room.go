package room

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tathagat/10minutechat/conf"
	"github.com/tathagat/10minutechat/websocketX"
)

type Room struct {
	ID        string
	Clients   map[*websocket.Conn]string // Map WebSocket connections to client names
	Broadcast chan Message               // Channel for broadcasting messages
	Mutex     sync.Mutex
	CreatedAt string // Timestamp for room creation
}

type Message struct {
	Timestamp string `json:"timestamp"` // Timestamp of the message
	Sender    string `json:"sender"`    // Name of the sender
	Content   string `json:"content"`   // Message content
}

var (
	rooms       = make(map[string]*Room)
	roomsMu     sync.Mutex
	randomNames = []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank"}
)

func getRandomName() string {
	return randomNames[rand.Intn(len(randomNames))]
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := uuid.New().String()
	roomID = strings.ReplaceAll(roomID, "-", "") // Remove dashes from the roomID
	roomsMu.Lock()
	defer roomsMu.Unlock()

	rooms[roomID] = &Room{
		ID:        roomID,
		Clients:   make(map[*websocket.Conn]string),
		Broadcast: make(chan Message),
		// TODO: time using client timezone (then also need to save server time)
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"), // Format timestamp
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

	log.Printf("Room created: %s at %s", roomID, rooms[roomID].CreatedAt)
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

	if len(room.Clients) >= conf.MaxRoomCapacity {
		http.Error(w, "Room is full", http.StatusForbidden)
		return
	}

	conn, err := websocketX.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading WebSocket connection: %v", err)
		http.Error(w, "Failed to establish WebSocket connection", http.StatusInternalServerError)
		return
	}

	clientName := getRandomName() // Assign a random name to the client
	// make sure the name is unique in the room
	// BEWARE: if list of random names is smaller than allowed clients per roon,
	// this can lead to infinite loop
	for {
		exists := false
		for _, name := range room.Clients {
			if name == clientName {
				exists = true
				break
			}
		}
		if !exists {
			break
		}
		clientName = getRandomName()
	}
	// Add the client to the room
	room.Mutex.Lock()
	room.Clients[conn] = clientName
	room.Mutex.Unlock()

	log.Printf("Client %s joined room: %s", clientName, roomID)
	go handleClientMessages(room, conn, clientName)
}
