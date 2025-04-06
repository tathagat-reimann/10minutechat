package conf

import (
	"log"
	"os"
	"strconv"
)

var (
	AllowedHost     string
	MaxRoomCapacity int
)

func init() {
	// Initialize AllowedHost
	AllowedHost = os.Getenv("ALLOWED_HOST")
	if AllowedHost == "" {
		AllowedHost = "localhost:8080" // Default to localhost if ALLOWED_HOST is not set
	}
	log.Printf("AllowedHost for Websocket: %s", AllowedHost)

	// Initialize MaxRoomCapacity
	MaxRoomCapacity = 2 // Default value
	capacity := os.Getenv("MAX_ROOM_CAPACITY")
	if capacity != "" {
		value, err := strconv.Atoi(capacity)
		if err == nil {
			MaxRoomCapacity = value
		}
	}
	log.Printf("MaxRoomCapacity: %d", MaxRoomCapacity)
}
