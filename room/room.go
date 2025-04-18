package room

import (
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tathagat/10minutechat/conf"
)

type Room struct {
	ID        string
	Clients   map[*websocket.Conn]string // Map WebSocket connections to client names
	Broadcast chan Message               // Channel for broadcasting messages
	Mutex     sync.Mutex
	CreatedAt string // Timestamp for room creation
}

type Message struct {
	//Timestamp string `json:"timestamp"` // Timestamp of the message
	Type    string `json:"type"`    // Type of the message (e.g., "chat", "clientName", "info")
	Sender  string `json:"sender"`  // Name of the sender
	Content string `json:"content"` // Message content
}

var (
	rooms    = make(map[string]*Room)
	roomsMu  sync.Mutex
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func getRandomName() string {
	return conf.RandomNames[rand.Intn(len(conf.RandomNames))]
}

func getNewClientName(usedNames []string) string {
	clientName := getRandomName() // Assign a random name to the client

	if slices.Contains(usedNames, clientName) {
		// If the name is already used, generate a new one
		clientName = getNewClientName(usedNames)
	}

	return clientName
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
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"), // Format timestamp
	}

	go sendNewMessageToAllClients(rooms[roomID])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"room_id": roomID}

	render.JSON(w, r, response)

	log.Printf("Room created: %s at %s", roomID, rooms[roomID].CreatedAt)
}

func CheckRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	log.Printf("Request to check room: %s", roomID)
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

	response := map[string]string{"room_id": roomID}
	render.JSON(w, r, response)
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

	// conn, err := websocketX.Upgrader.Upgrade(w, r, nil)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading WebSocket connection: %v", err)
		http.Error(w, "Failed to establish WebSocket connection", http.StatusInternalServerError)
		return
	}

	// Extract the list of used names from the room's Clients map
	room.Mutex.Lock()
	usedNames := make([]string, 0, len(room.Clients))
	for _, name := range room.Clients {
		usedNames = append(usedNames, name)
	}
	room.Mutex.Unlock()

	// Assign a unique name to the client
	clientName := getNewClientName(usedNames)

	// Add the client to the room
	room.Mutex.Lock()
	room.Clients[conn] = clientName
	room.Mutex.Unlock()

	log.Printf("Client %s joined room: %s", clientName, roomID)
	go sendClientNameToItself(room, conn, clientName)
	go sendClientNameToOtherClients(room, conn, clientName)
	go handleNewMessageFromClient(room, conn, clientName)
}
