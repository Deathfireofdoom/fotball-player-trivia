package db

import (
	"context"
	"fmt"

	"github.com/Deathfireofdoom/fotball-player-trivia-api/pkg/utils"
)

func GetPlayerInfo(playerName string) (PlayerInfo, error) {
	var playerInfo PlayerInfo
	row := db.QueryRowContext(context.Background(), `SELECT name, country, height, weight FROM player_info WHERE name ILIKE CONCAT('%', $1::varchar, '%')`, playerName)

	if err := row.Err(); err != nil {
		fmt.Println(err)
		fmt.Println(playerName)
		return PlayerInfo{}, fmt.Errorf("Could not get player from database: %w", err)
	}

	if err := row.Scan(&playerInfo.Name, &playerInfo.Country, &playerInfo.Height, &playerInfo.Weight); err != nil {
		fmt.Println("Could not parse")
		// sql: no rows in result set
		fmt.Println(err)
		return PlayerInfo{}, fmt.Errorf("Could not parse response from db: %w", err)
	}
	fmt.Println(playerInfo)
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
