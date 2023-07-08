package rest

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcusadriano/go-sse-example/internal/log"

	"github.com/redis/go-redis/v9"
)

type ChatHandler struct {
	rdb *redis.Client
}

func NewChatHandler(rdb *redis.Client) *ChatHandler {

	return &ChatHandler{
		rdb: rdb,
	}
}

func (c ChatHandler) RegisterChatRoutes(r *chi.Mux) {

	r.Route("/api/chat/v1", func(r chi.Router) {
		r.Post("/send/{userId}", c.send)
		r.Get("/stream/{userId}", c.subscribe)
	})
}

func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (c *ChatHandler) send(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	cors(w, r)

	userId := chi.URLParam(r, "userId")
	topic := "chat:" + userId

	log.WithContext(ctx).Debug().Str("topic", topic).Msg("Publishing message")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithContext(ctx).Error().Err(err).Msg("Error reading request body")
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	err = c.rdb.Publish(ctx, topic, body).Err()

	if err != nil {
		log.WithContext(ctx).Err(err).Msg("Error publishing message")
	}
}

func (c *ChatHandler) subscribe(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	cors(w, r)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.WithContext(ctx).Error().Msg("Streaming unsupported!")
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	flusher.Flush()

	userId := chi.URLParam(r, "userId")

	subscriber := newSseRedisSubscriber(
		ctx,
		w,
		c.rdb,
		func() string {
			return "chat:" + userId
		},
	)

	go subscriber.Listen(ctx)

	<-r.Context().Done()

	log.WithContext(ctx).Debug().Str("user_id", userId).Msg("Client closed connection")
	subscriber.Close(ctx)
}