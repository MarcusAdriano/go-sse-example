package rest

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcusadriano/go-sse-example/internal/log"
	"github.com/redis/go-redis/v9"
)

const (
	topicsPrefix = "chat:instances:"
	connectionsKey = "chat:%s:connections"

)

type RedisInstanceManager struct {
	id  string
	ctx context.Context
	rdb *redis.Client
	rps *redis.PubSub
	messageChannel <-chan string
}

func NewManager(ctx context.Context, rdb *redis.Client) *RedisInstanceManager {

	id := uuid.New().String()
	rps := rdb.Subscribe(ctx, topicsPrefix+id)

}

func (s *RedisInstanceManager) Close(ctx context.Context) {
	err := s.rps.Close()
	if err != nil {
		log.WithContext(ctx).Error().Msg("Error closing redis subscriber")
	}
}

func (s *RedisInstanceManager) Listen(ctx context.Context) {
	msgs := s.subscribe(ctx)

	for msg := range msgs {
		s.sendSseEvent(ctx, msg.Payload)
	}}

func (s *RedisInstanceManager) subscribe(ctx context.Context) <-chan *redis.Message {
	topicStr := s.topic()
	log.WithContext(ctx).Debug().Str("topic", topicStr).Msg("Subscribing to topic")
	s.rps = s.rdb.Subscribe(ctx, topicStr)
	return s.rps.Channel()
}