package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// New returns a configured Redis client. The caller is responsible for closing it.
func New(host, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})
}

// DSN returns a human-readable connection string for logging.
func DSN(host string, db int) string {
	return fmt.Sprintf("redis://%s/%d", host, db)
}
