package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
	"github.com/go-redis/redis/v8"
)

var (
	RedisAddress string = "localhost:6379" // TODO make this env variable
)

var Client *RedisClient

func InitializeRedis() {
	Client = NewRedisClient()
}

type RedisClient struct {
	client *redis.Client
}

// NewRedisClient returns RedisClient
func NewRedisClient() *RedisClient {
	client := redis.NewClient(
		&redis.Options{
			Addr:        RedisAddress,
			DB:          0,
			DialTimeout: 100 * time.Millisecond,
			ReadTimeout: 100 * time.Millisecond,
		})

	// Checks if redis is reachable.
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil
	}

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
