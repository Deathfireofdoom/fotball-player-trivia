package controller

import (
	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
	"github.com/Deathfireofdoom/fotball-player-trivia/service"
	"github.com/gin-gonic/gin"
)

type TriviaController interface {
	GetPlayerTrivia(ctx *gin.Context) (int, entity.PlayerTrivia)
}

type triviaController struct {
	service service.TriviaService
}

func NewTriviaController(service service.TriviaService) TriviaController {
	return triviaController{
		service: service,
	}
}

func (tc triviaController) GetPlayerTrivia(ctx *gin.Context) (int, entity.PlayerTrivia) {
	playerParam, ok := ctx.GetQuery("name")
	if ok {
		playerTrivia, err := tc.service.GetPlayerTrivia(playerParam)
		if err != nil {
			panic("Failed at getting player trivia.")
		}
		return 200, playerTrivia
	}
	panic("No Param was given.") //TODO fix this, so it returns something better than panic.
}
