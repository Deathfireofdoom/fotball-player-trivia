package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/entity"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/logger"
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
	logger.InfoLogger.Printf("Initializing redis at %s:%s", host, port)
	client := redis.NewClient(
		&redis.Options{
			Addr:        fmt.Sprintf("%s:%s", host, port),
			DB:          0,
			DialTimeout: 100 * time.Millisecond,
			ReadTimeout: 100 * time.Millisecond,
		})

	// Checks if redis is reachable.
	logger.InfoLogger.Println("Pinging redis.")
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		logger.WarningLogger.Println("Could not ping redis, will continue without cachce.")
		return nil
	}
	logger.InfoLogger.Println("Sucessfully reached redis.")
	return &RedisClient{
		client: client,
	}

}

// GetTrivia query the cache to see if player trivia is already present.
func (c *RedisClient) GetTrivia(ctx context.Context, playerName string) (entity.PlayerTrivia, error) {
	logger.InfoLogger.Printf("Search cache for player: %s", playerName)
	cmd := c.client.Get(ctx, playerName)

	cmdb, err := cmd.Bytes()
	if err != nil {
		logger.InfoLogger.Println("No hit in cache.")
		return entity.PlayerTrivia{}, err
	}

	b := bytes.NewReader(cmdb)

	var playerTrivia entity.PlayerTrivia
	if err := gob.NewDecoder(b).Decode(&playerTrivia); err != nil {
		logger.WarningLogger.Println("Could not decode response from cache, will fetch results again: %w", err)
		return entity.PlayerTrivia{}, err
	}

	logger.InfoLogger.Println("Hit in cache.")
	return playerTrivia, nil
}

// SaveTrivia saves trivia with playername as lookup key.
func (c *RedisClient) SaveTrivia(ctx context.Context, playerTrivia entity.PlayerTrivia, durationSeconds float64) error {
	logger.InfoLogger.Println("Saving trivia to cachce.")
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(playerTrivia); err != nil {
		logger.ErrorLogger.Println("Could not save trivia to cachce: %w", err)
		return err
	}

	err := c.client.Set(ctx, playerTrivia.Name, b.Bytes(), time.Duration(durationSeconds)*time.Second).Err()
	if err != nil {
		logger.ErrorLogger.Println("Could not save trivia to cache: %w", err)
		return err
	}
	logger.InfoLogger.Println("Successfully saved trivia to cachce.")
	return nil
}
