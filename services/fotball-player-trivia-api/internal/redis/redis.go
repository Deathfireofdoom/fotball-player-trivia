package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/entity"
)

var Client *RedisClient

func InitializeRedis(host, port string) {

	Client = NewRedisClient(host, port)
}

type RedisClient struct {
	client *redis.Client
}

// NewRedisClient returns RedisClient
func NewRedisClient(host, port string) *RedisClient {
	client := redis.NewClient(
		&redis.Options{
			Addr:        fmt.Sprintf("%s:%s", host, port),
			DB:          0,
			DialTimeout: 100 * time.Millisecond,
			ReadTimeout: 100 * time.Millisecond,
		})

	// Checks if redis is reachable.
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		fmt.Println("Redis is not reachable")
		return nil
	}
	fmt.Println("Reached redis")

	return &RedisClient{
		client: client,
	}

}

// GetTrivia query the cache to see if player trivia is already present.
func (c *RedisClient) GetTrivia(ctx context.Context, playerName string) (entity.PlayerTrivia, error) {
	cmd := c.client.Get(ctx, playerName)

	cmdb, err := cmd.Bytes()
	if err != nil {
		return entity.PlayerTrivia{}, err
	}

	b := bytes.NewReader(cmdb)

	var playerTrivia entity.PlayerTrivia

	if err := gob.NewDecoder(b).Decode(&playerTrivia); err != nil {
		return entity.PlayerTrivia{}, err
	}

	return playerTrivia, nil
}

// SaveTrivia saves trivia with playername as lookup key.
func (c *RedisClient) SaveTrivia(ctx context.Context, playerTrivia entity.PlayerTrivia, durationSeconds float64) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(playerTrivia); err != nil {
		return err
	}

	return c.client.Set(ctx, playerTrivia.Name, b.Bytes(), time.Duration(durationSeconds)*time.Second).Err()
}
