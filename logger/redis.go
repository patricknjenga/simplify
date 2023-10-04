package logger

import (
	"context"
	"encoding/json"
	"log"
	"runtime/debug"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLogger struct {
	Client  *redis.Client
	Service string
}

func NewRedisLogger(client *redis.Client, service string) Logger {
	return &RedisLogger{client, service}
}

func (r *RedisLogger) Error(e error) {
	data, err := json.Marshal(Error{
		Message: e.Error(),
		Service: r.Service,
		Time:    time.Now(),
		Trace:   string(debug.Stack()),
	})
	if err != nil {
		log.Println(err)
	}
	r.Client.Publish(context.Background(), "error", string(data))
}
