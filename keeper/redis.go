package keeper

import (
	"context"
	"encoding/base64"

	"github.com/patricknjenga/simplify/logger"
	"github.com/redis/go-redis/v9"
)

type RedisKeeper struct {
	Client *redis.Client
	Logger logger.Logger
}

func NewRedisKeeper(client *redis.Client, logger logger.Logger, service string) Keeper {
	return &RedisKeeper{client, logger}
}

func (r RedisKeeper) Get(key string) string {
	val, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		r.Logger.Error(err)
		return ""
	}
	decoded, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		r.Logger.Error(err)
		return ""
	}
	return string(decoded)
}

func (r RedisKeeper) Set(key string, value string) {
	value = base64.StdEncoding.EncodeToString([]byte(value))
	r.Client.Set(context.Background(), key, value, 0)
}
