package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Deathfireofdoom/fotball-player-trivia/controller"
	"github.com/Deathfireofdoom/fotball-player-trivia/database"
	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
	redisClient "github.com/Deathfireofdoom/fotball-player-trivia/redis"
	"github.com/Deathfireofdoom/fotball-player-trivia/service"
	"github.com/gin-gonic/gin"
)

var (
	triviaService    service.TriviaService       = service.NewTriviaService()
	triviaController controller.TriviaController = controller.NewTriviaController(triviaService)
)

func main() {
	server := gin.Default()
	redisClient.InitializeRedis()

	database.StartBatchProcessExample()

	// Get player trivia
	server.GET("/player-trivia", func(ctx *gin.Context) {
		status, playerTriva := triviaController.GetPlayerTrivia(ctx)
		ctx.JSON(status, playerTriva)
	})

	// Used to ping api to check if it is up.
	server.GET("/test", ConnectionTest)

	// Function used to test functionallity in a lazy way.
	server.GET("/dummy", DummyFunction)

	server.Run()
}

// ConnectionTest is used to ping the server.
func ConnectionTest(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}

func DummyFunction(ctx *gin.Context) {
	resp, err := http.Get("https://restcountries.com/v3.1/name/denmark?fullText=true")
	if err != nil {
		panic("OHNO")
	}
	defer resp.Body.Close()

	var response entity.RestcountriesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic("OHNO could not decode.")
	}
	fmt.Println(response[0].Name.OfficialName)

	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}
