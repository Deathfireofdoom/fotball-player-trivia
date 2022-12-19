package router

import (
	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/db"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/logger"
	"github.com/gin-gonic/gin"
)

func migrateDatabase(ctx *gin.Context) {
	logger.InfoLogger.Printf("Starting migrating process.") // TODO add to see if database is migrated
	db.LoadPlayerData("./internal/assets/data/player-data-set.csv", "player_info")
	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}
