package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tathagat/10minutechat/room"
)

func SetupRouter(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/room", room.CreateRoom)
		r.Get("/room/{id}/join", room.JoinRoom)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
}
