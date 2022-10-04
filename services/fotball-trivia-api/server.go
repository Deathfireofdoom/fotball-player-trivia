package main

import (
	"github.com/Deathfireofdoom/fotball-player-trivia/controller"
	"github.com/Deathfireofdoom/fotball-player-trivia/service"
	"github.com/gin-gonic/gin"
)

var (
	triviaService    service.TriviaService       = service.NewTriviaService()
	triviaController controller.TriviaController = controller.NewTriviaController(triviaService)
)

func main() {
	server := gin.Default()

	// Get player trivia
	server.GET("/player-trivia", func(ctx *gin.Context) {
		status, playerTriva := triviaController.GetPlayerTrivia(ctx)
		ctx.JSON(status, playerTriva)
	})

	// Used to ping api to check if it is up.
	server.GET("/test", ConnectionTest)

	server.Run()
}

// ConnectionTest is used to ping the server.
func ConnectionTest(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}
