package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Is this allowed? The plan is to use database.DbService.saveCustomer() in another
// file.
var DbService DatabaseService

func Connect(host, user, password, dbName string) {
	DbService = NewDatabaseService(host, user, password, dbName)
	DbService.InitDB()
}

func Close() {
	DbService.CloseDB()
}

type DatabaseService interface {
	InitDB() error
	CloseDB() error
	PingDB() error
	migrateDB() error
	GetPlayerInfo(string) (entity.PlayerInfoDB, error)
}

type databaseService struct {
	dns url.URL
	db  *sql.DB
}

func NewDatabaseService(host, user, password, dbName string) DatabaseService {
	dns := url.URL{
		Scheme: "postgres",
		Host:   host,
		User:   url.UserPassword(user, password),
		Path:   dbName,
	}
	q := dns.Query()
	q.Add("sslmode", "disable")

	dns.RawQuery = q.Encode()

	return &databaseService{dns: dns}
}

// initDB connects to db.
func (ds *databaseService) InitDB() error {
	db, err := sql.Open("pgx", ds.dns.String())
	if err != nil {
		fmt.Println("sql.Open", err)
	}

	ds.db = db
	return nil
}

// CloseDB closes DB.
func (ds *databaseService) CloseDB() error {
	fmt.Println("db closed")
	return ds.db.Close()
}

// pingDB pings db to see if connected.
func (ds *databaseService) PingDB() error {
	err := ds.db.PingContext(context.Background())
	if err != nil {
		fmt.Println("db.ping", err)
		return err
	}
	fmt.Println("Sucessfully pinged db.")
	return nil
}

func (ds *databaseService) migrateDB() error {
	panic("IMPLEMEN THIS")
}

func (ds *databaseService) GetPlayerInfo(playerName string) (entity.PlayerInfoDB, error) {
	var playerInfoDB entity.PlayerInfoDB
	row := ds.db.QueryRowContext(context.Background(), `SELECT name, country FROM player_info WHERE name like %lower($1)%`, playerName)
	if err := row.Err(); err != nil {
		panic("Could not get name from db.") // TODO send back empty and raise "Player was not found error."
	}

	if err := row.Scan(&playerInfoDB.Name, &playerInfoDB.Country); err != nil {
		panic("Could not get name from db,") // ToDO swen back empty and raise player was not found error.
	}

	return playerInfoDB, nil
}

//https://codewithyury.com/golang-wait-for-all-goroutines-to-finish/
