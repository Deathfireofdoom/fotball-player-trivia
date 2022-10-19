package db

import (
	"fmt"
)

var db *sql.DB

type Config struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DB       string `env:"DB_NAME"`
	SSLMode  string `env:"DB_SSLMODE"`
}

func (c *Config) getConnstring() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
		c.SSLMode,
	)
}

func InitDB(c Config) error {
	var err error

	db, err = sql.Open("pgx", c.getConnstring())
	if err != nil {
		return fmt.Errorf("Could not connect to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("Could not ping db: %w", err)
	}
	return nil
}
