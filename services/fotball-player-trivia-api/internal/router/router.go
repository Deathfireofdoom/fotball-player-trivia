package router

import (
	"github.com/gin-gonic/gin"
)

func Configure() {
	router := gin.Default()

	router.GET("/player-trivia")
	router.GET("/test", connectionTest)

}
