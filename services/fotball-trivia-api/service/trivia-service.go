package service

import (
	"time"

	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
	"github.com/Deathfireofdoom/fotball-player-trivia/redis"
	"github.com/gin-gonic/gin"
)

type TriviaService interface {
	GetPlayerTrivia(string, *gin.Context) (entity.PlayerTrivia, error)
}

type triviaService struct{}

func NewTriviaService() TriviaService {
	return &triviaService{}
}

// GetPlayerTrivia gets trivia about player by first checking if the playerTrivia already reside in Redis-cache
// if not it starts the full process to fetch the trivia.
func (ts *triviaService) GetPlayerTrivia(playerName string, ctx *gin.Context) (entity.PlayerTrivia, error) {
	// Tries to get value from Redis.
	playerTrivia, err := redis.Client.GetTrivia(ctx, playerName)
	// This means the PlayerTrivia was found in Redis, returning it.
	if err == nil {
		return playerTrivia, nil
	}

	// Calls the full process to fetch trivia since it was not cached.
	playerTrivia, err = MockGetTrivia(playerName)

	if err != nil {
		panic("Could not get PlayerTrivia") // TODO make this more elegant.
	}

	// Saves playerTrivia into Cache.
	err = redis.Client.SaveTrivia(ctx, playerTrivia, 25)
	if err != nil {
		panic("Could not save playerTrivia to Redis.") // TODO make this more elegant.
	}

	// Returns trivia
	return playerTrivia, nil
}

// MockGetTrivia just a mock function to get trivia, should be replaced soon.
func MockGetTrivia(playerName string) (entity.PlayerTrivia, error) {
	time.Sleep(time.Duration(2) * time.Second) // Simulate some kind of delay.
	return entity.PlayerTrivia{Name: playerName, Country: "Test"}, nil
}

func FetchTrivia(playerName string) (entity.PlayerTrivia, error) {
	panic("Implement this.")
	// 1. Fetch country for the player from DB.

	// This part should be concurrent.
	// 2.1.1 Make API call to get Land area and Population for the players country.
	// 2.1.2 Do stupid calculation on Area of land and Percentage of population.

	// 2.2.1 Make api call to get official name of country from restcountriesapi, ex. Kingdom of Sweden

	// 3. Wait for all to finish, then make playerTrivia and return it.

}
