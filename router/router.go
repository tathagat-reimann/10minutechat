package router

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/gorilla/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/tathagat/10minutechat/room"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func SetupRouter(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/rooms", CreateRoom)
		r.Get("/rooms/{id}", CheckRoom)
		r.Get("/rooms/{id}/join", JoinRoom)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Serve the entire static directory
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	// Create a new room and return the room ID
	roomID := room.CreateRoom()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"room_id": roomID}

	render.JSON(w, r, response)
}

func CheckRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	exists, full, room := room.CheckRoom(roomID)

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if full {
		http.Error(w, "Room is full", http.StatusForbidden)
		return
	}

	response := map[string]string{"room_id": room.ID}
	render.JSON(w, r, response)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	exists, full, _ := room.CheckRoom(roomID)

	if !exists || full {
		// return error
		http.Error(w, "Failed to join room", http.StatusInternalServerError)
		return
	}

	// conn, err := websocketX.Upgrader.Upgrade(w, r, nil)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading WebSocket connection: %v", err)
		http.Error(w, "Failed to establish WebSocket connection", http.StatusInternalServerError)
		return
	}

	room.JoinRoom(roomID, conn)
}
