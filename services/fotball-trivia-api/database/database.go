package database

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"sync"

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

func StartBatchProcessExample() {
	filePath := "/Users/oskarelvkull/Documents/big-corp/fotball-player-trivia/services/fotball-trivia-api/database/data/player-data-set.csv"
	batchProcessFile(filePath, 1, 10)

}

func batchProcessFile(filePath string, concurrency int, batchSize int) {
	// Cancel chanel, used to communicate that the file has been read successfuly.
	// The channel needs to have 1 buffer since last signal will not be read.
	cancelCh := make(chan bool, 1)
	// Chancel to send the batches to other go-routines. The channel will have 5 messages in the buffer.
	batchCh := make(chan []string, concurrency+5)

	// Using a wait group so we can gracefully exit without loosing any batches, or atleast not loosing any batch due to premature exit heh.
	wg := new(sync.WaitGroup)

	// Starts the concurrency.
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		go exmapleConsumer(wg, batchCh, cancelCh)
	}

	// Creates a file reader that later will be buff read.
	file, err := os.Open(filePath)
	if err != nil {
		panic("Could not read file.")
	}
	defer file.Close()

	// Buffering of batch processing.
	scanner := bufio.NewScanner(file)
	batch := []string{}

	for scanner.Scan() {
		// Checks if batch-size is met. If so, batch it publish to chanel where the go-routines listens to.
		if len(batch) >= batchSize {
			batchCh <- batch

			// Resets batch.
			batch = []string{}
		}
		// Add line to batch.
		batch = append(batch, scanner.Text())
	}
	// Publish last batch even though max-size is not met.
	batchCh <- batch
	// Publishing a nil to communicate that the last message has been read.
	batchCh <- nil

	// Gotta figure out why we need to call wg.Done() one extra time..
	fmt.Println("LOG - Waiting for processes to finish.")
	//wg.Done()
	wg.Wait()
	// Maybe should delete the cancelCancel? Not sure if that is needed, a true value will be stuck there.
	fmt.Println("LOG - All processes finished, successfully processed the full file-.")
}

func exmapleConsumer(wg *sync.WaitGroup, batchCh chan []string, cancelCh chan bool) {
	defer wg.Done()

	for {
		select {
		case batch := <-batchCh:
			if batch != nil {
				exampleProcessing(batch)
			} else {
				fmt.Println("Publishing cancel signal")
				cancelCh <- true
				return
			}
		case signal := <-cancelCh:
			// If a message is publish in the cancelCh all lines has been read.
			// Not sure if this is the best solution since we need to keep publish
			// a new cancel message. But I like the solution.

			fmt.Println("Publishing cancel signal")
			cancelCh <- signal
			return
		}

	}
}

func exampleProcessing(batch []string) {
	//fmt.Println("Processing file.")
	//fmt.Println(batch)
}
