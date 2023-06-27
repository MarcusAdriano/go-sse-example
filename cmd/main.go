package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/marcusadriano/go-redis/internal/log"
	"github.com/marcusadriano/go-redis/internal/rest"
	"github.com/redis/go-redis/v9"
)

func main() {
	log.Default().Info().Msg("Starting server")

	mux := chi.NewMux()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(log.RestLogger)
	mux.Use(middleware.Recoverer)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	chatRest := rest.NewChatHandler(rdb)
	chatRest.RegisterChatRoutes(mux)

	err := http.ListenAndServe(":8080", mux)
	log.Default().Err(err).Msg("Error starting server")
}
