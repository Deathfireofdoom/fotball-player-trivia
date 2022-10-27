package router

import (
	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/trivia"
	"github.com/gin-gonic/gin"
)

func getPlayerTrivia(ctx *gin.Context) {
	playerName, ok := ctx.GetQuery("name")
	if ok {
		playerTrivia, err := trivia.GetPlayerTriva(playerName, ctx)
		if err != nil {
			ctx.JSON(500, map[string]string{"message": playerTrivia.Name})
			return
		}
		ctx.JSON(200, playerTrivia)
		return
	}
	ctx.JSON(400, map[string]string{"message": "Missing parameter name."})
}
