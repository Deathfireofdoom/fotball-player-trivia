package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
	"github.com/Deathfireofdoom/fotball-player-trivia/utils"
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
	LoadPlayerData() error
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

func (ds *databaseService) LoadPlayerData() error {
	// Create table.
	createTableSql := "CREATE TABLE IF NOT EXISTS player_info( name TEXT, height INT, weight INT, country TEXT) "
	ds.db.Exec(createTableSql)

	// Start batch processing.
	filePath := "/Users/oskarelvkull/Documents/big-corp/fotball-player-trivia/services/fotball-trivia-api/database/data/player-data-set.csv"
	utils.BatchProcessFile(filePath, ds.loadBatch, 1, 500)
	return nil
}

func (ds *databaseService) loadBatch(batch []string) {
	// Sql used for insert.
	insertSql := "INSERT INTO player_info (name, height, weight, country) VALUES ($1, $2, $3, $4)"

	// Extract the value.
	valuesList := [][]string{}
	for _, line := range batch {
		valuesList = append(valuesList, parseLine(line))
	}
	fmt.Println(valuesList)
	for _, values := range valuesList {
		//ds.db.Exec(insertSql, values)
		fmt.Println(values[0], values[1], values[2], values[3], insertSql)
		height, err := strconv.Atoi(values[1])
		if err != nil {
			panic("ERROR")
		}
		weight, err := strconv.Atoi(values[2])
		if err != nil {
			panic("ERROR")
		}

		ds.db.Exec(insertSql, values[0], height, weight, values[3])

	}

}

func parseLine(line string) []string {
	return strings.Split(line, ",")
}

//https://codewithyury.com/golang-wait-for-all-goroutines-to-finish/
