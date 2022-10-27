package main

import (
	"os"

	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/db"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/redis"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/router"
)

func main() {
	router := router.Configure()
	router.Run()
}

func init() {
	// Initialize
	initializeDB()
	initializeRedis()
}

func initializeRedis() {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	redis.InitializeRedis(host, port)
}

func initializeDB() {
	db_host := os.Getenv("DB_HOST")
	if db_host == "" {
		db_host = "0.0.0.0"
	}

	db_port := os.Getenv("DB_PORT")
	if db_port == "" {
		db_port = "5432"
	}

	db_user := os.Getenv("DB_USER")
	if db_user == "" {
		panic("no username for database was supplied in env.")
	}

	db_password := os.Getenv("DB_PASSWORD")
	if db_password == "" {
		panic("no password for database was supplied in env.")
	}

	db_default_db := os.Getenv("DB_DEFAULT_DB")
	if db_default_db == "" {
		db_default_db = "dev"
	}

	db.InitDB(db.Config{
		Host:     db_host,
		Port:     db_port,
		User:     db_user,
		Password: db_password,
		DB:       db_default_db,
	})
}
