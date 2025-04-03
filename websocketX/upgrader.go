package websocketX

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

// TODO config

// MaxRoomCapacity dynamically fetches the room capacity from the environment variable
func MaxRoomCapacity() int {
	capacity := os.Getenv("MAX_ROOM_CAPACITY")
	value, err := strconv.Atoi(capacity)
	if err == nil {
		return value
	}
	log.Printf("Invalid MAX_ROOM_CAPACITY value, defaulting to 2: %v", err)
	return 2 // Default value
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		allowedHost := os.Getenv("ALLOWED_HOST")
		if allowedHost == "" {
			allowedHost = "localhost:8080" // Default to localhost if ALLOWED_HOST is not set
		}
		return r.Host == allowedHost
	},
}
