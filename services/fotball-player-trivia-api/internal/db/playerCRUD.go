package db

func GetPlayerInfo(playerName string) (PlayerInfo, error) {
	var playerInfo PlayerInfo
	row := db.QueryRowContext(context.Background(), `SELECT name, country, height, weight FROM player_info WHERE name like %lower($1)%`, playerName)

	if err := row.Err(); err != nil {
		return PlayerInfo{}, fmt.Errorf("Could not get player from database: %w", err)
	}

	if err := row.Scan(&playerInfo.Name, &playerInfo.Country, &playerInfo.Height, &playerInfo.Weight); err != nil {
		return PlayerInfo{}, fmt.Errorf("Could not parse response from db: %w", err)
	}

	return playerInfo, nil
}

func LoadPlayerData(filePath, schema, table string) error {
	// Create table.
	createTableSql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s( name TEXT, height INT, weight INT, country TEXT)", schema, table)
	err := db.Exec(createTableSql) // What does this return, maybe we should handle this.
	if err != nil {
		return fmt.Errorf("Could not create table: %w", err)
	}

	// Start batch processing.
	utils.BatchProcessFile(filePath)

}

func loadBatch(batch []string, schema, table string) {
	// Sql used for insert.
	insertSql := fmt.Sprintf("INSERT INTO %s.%s (name, height, weight, country) VALUES ($1, $2, $3, $4)", schema, table)

	// Extract the value.
	valuesList := [][]string{}
	for _, line := range batch {
		valuesList = append(valuesList, utils.ParseLine(line, ","))
	}

}
