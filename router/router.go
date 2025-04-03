package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/tathagat/10minutechat/room"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/room", room.CreateRoom)
	r.Get("/room/{id}/join", room.JoinRoom)
	return r
}
