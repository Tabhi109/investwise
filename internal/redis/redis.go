package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis wraps the go-redis client connection
type Redis struct {
	Client *redis.Client
}

// Connect parses the Redis URL and establishes connectivity
func Connect(ctx context.Context, url string) (*Redis, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	// Tweak client options if needed (e.g. read/write timeouts)
	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return &Redis{Client: client}, nil
}

// Close closes the underlying Redis client connection pool
func (r *Redis) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}
