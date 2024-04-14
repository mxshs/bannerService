package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
    rdb *redis.Client
    ttl time.Duration
}

func NewRedisCache(address string, password string, db int) *RedisCache {
    rdb := redis.NewClient(&redis.Options{
        Addr: address,
        Password: password,
        DB: db,
    })

    return &RedisCache{
        rdb: rdb,
        ttl: 5 * time.Minute,
    }
}
