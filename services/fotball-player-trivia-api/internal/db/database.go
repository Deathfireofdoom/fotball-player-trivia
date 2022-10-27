package db

import (
	"database/sql"
	"fmt"

	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

type Config struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DB       string `env:"DB_NAME"`
}

func (c *Config) getConnstring() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
	)
}

func InitDB(c Config) error {
	var err error
	logger.InfoLogger.Println("Initializing postgres db.")
	logger.InfoLogger.Println("Opening connection to postgres db.")
	db, err = sql.Open("pgx", c.getConnstring())
	if err != nil {
		logger.ErrorLogger.Println("Could not open connection to db: %w", err)
		return fmt.Errorf("could not connect to postgres-db: %w", err)
	}

	logger.InfoLogger.Println("Pinging postgres...")
	err = db.Ping()
	if err != nil {
		logger.ErrorLogger.Println("Could not ping postgres db: %w", err)
		return fmt.Errorf("Could not ping db: %w", err)
	}
	logger.InfoLogger.Println("Succesfully connected to postgres DB.")
	return nil
}
