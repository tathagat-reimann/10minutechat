package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/tathagat/10minutechat/router"
)

func main() {
	r := chi.NewRouter()

	// use middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	router.SetupRouter(r)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
