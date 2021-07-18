package cache

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
}

var redisClient *redis.Client

func NewCacheDB(url, pass string) *Redis {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pass, // no password set
		DB:       0,    // use default DB
	})

	return &Redis{Client: redisClient}
}

func Client() *redis.Client {
	return redisClient
}
