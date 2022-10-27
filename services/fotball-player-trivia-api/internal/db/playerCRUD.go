package db

import (
	"context"
	"fmt"

	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/logger"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/pkg/utils"
)

func GetPlayerInfo(playerName string) (PlayerInfo, error) {
	logger.InfoLogger.Printf("Getting player info for %s", playerName)
	var playerInfo PlayerInfo
	row := db.QueryRowContext(context.Background(), `SELECT name, country, height, weight FROM player_info WHERE name ILIKE CONCAT('%', $1::varchar, '%')`, playerName)

	if err := row.Err(); err != nil {
		logger.ErrorLogger.Println("Could not query database: %w", err)
		return PlayerInfo{}, fmt.Errorf("Could not get player from database: %w", err)
	}

	logger.InfoLogger.Println("Successfully queried database, parsing response...")
	if err := row.Scan(&playerInfo.Name, &playerInfo.Country, &playerInfo.Height, &playerInfo.Weight); err != nil {
		logger.WarningLogger.Println("Could not parse response, possible no match: %w", err)
		return PlayerInfo{}, fmt.Errorf("Could not parse response from db, possible no match: %w", err)
	}
	logger.InfoLogger.Printf("Successfully queried database for player %s", playerName)
	return playerInfo, nil
}

func LoadPlayerData(filePath, schema, table string) error {
	// Create table.
	createTableSql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s( name TEXT, height INT, weight INT, country TEXT)", schema, table)
	_, err := db.Exec(createTableSql)
	if err != nil {
		return fmt.Errorf("Could not create table: %w", err)
	}

	// Start batch processing.
	utils.BatchProcessFile(filePath, func(batch []string) { loadBatch(batch, schema, table) }, 1, 500)
	return nil
}

func loadBatch(batch []string, schema, table string) {
	/*// Sql used for insert.
	insertSql := fmt.Sprintf("INSERT INTO %s.%s (name, height, weight, country) VALUES ($1, $2, $3, $4)", schema, table)

	// Extract the value.
	valuesList := [][]string{}
	for _, line := range batch {
		valuesList = append(valuesList, utils.ParseLine(line, ","))
	}*/
	panic("Implemen this")
}
