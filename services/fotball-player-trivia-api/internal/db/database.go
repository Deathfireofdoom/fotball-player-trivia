package db

import (
	"database/sql"
	"fmt"

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
	db, err = sql.Open("pgx", c.getConnstring())
	fmt.Println(c.getConnstring())
	if err != nil {
		fmt.Println("Could not open to db.")
		fmt.Println(err)
		return fmt.Errorf("Could not connect to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Could not ping to db.")
		return fmt.Errorf("Could not ping db: %w", err)
	}
	fmt.Println("Reached db")
	return nil
}
