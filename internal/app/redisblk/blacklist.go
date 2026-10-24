package redisblk

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Blacklist struct {
	client *redis.Client
}

func NewFromEnv() (*Blacklist, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6380"
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Blacklist{client: client}, nil
}

func key(token string) string {
	sum := sha256.Sum256([]byte(token))
	return "dendro:jwt:blacklist:" + hex.EncodeToString(sum[:])
}

func (b *Blacklist) Add(ctx context.Context, token string, ttl time.Duration) error {
	return b.client.Set(ctx, key(token), "1", ttl).Err()
}

func (b *Blacklist) Has(ctx context.Context, token string) (bool, error) {
	count, err := b.client.Exists(ctx, key(token)).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
