package router

import (
	"github.com/gin-gonic/gin"
)

func Configure() *gin.Engine {
	router := gin.Default()

	router.GET("/player-trivia", getPlayerTrivia)
	router.GET("/test", connectionTest)

	return router
}
