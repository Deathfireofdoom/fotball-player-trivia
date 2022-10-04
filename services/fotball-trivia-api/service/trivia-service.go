package service

import (
	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
)

type TriviaService interface {
	GetPlayerTrivia(string) (entity.PlayerTrivia, error)
}

type triviaService struct{}

func NewTriviaService() TriviaService {
	return &triviaService{}
}

func (ts *triviaService) GetPlayerTrivia(playerName string) (entity.PlayerTrivia, error) {
	return entity.PlayerTrivia{Name: playerName, Country: "Test"}, nil
}
