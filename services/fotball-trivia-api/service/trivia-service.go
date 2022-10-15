package service

import (
	"fmt"
	"math"
	"time"

	"github.com/Deathfireofdoom/fotball-player-trivia/database"
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
	playerTrivia, err = FetchTrivia(playerName)
	FetchTrivia(playerName)
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
	// 1. Fetch country and height for the player from DB.
	playerInfoDB, err := database.DbService.GetPlayerInfo(playerName)
	if err != nil {
		panic("Failed fetching from db.")
	}
	tmpCountryName := playerInfoDB.Country
	fmt.Println(tmpCountryName)
	countryName := "sweden"
	var playerHeight float64 = 150
	var playerWeight float64 = 70

	// Here we concurrently starts two processes.
	// Is this really needed? Maybe not, but since we are calling two different apis we probably
	// get some performance increase. This would be even more important if the underlying
	// processes was heavy IO bound, like delegating the math to another service.
	chanFunFact := make(chan entity.CalculatedFunFact, 1)
	chanOfficialName := make(chan string, 1)

	go calculateFunFacts(countryName, playerHeight, playerWeight, chanFunFact)
	go getOfficialName(countryName, chanOfficialName)

	// Waits for go-rountines to finish.
	funFact := <-chanFunFact
	officialName := <-chanOfficialName

	// Creates object and send it back.
	playerTrivia := entity.PlayerTrivia{
		Name:                playerName,
		Country:             countryName,
		CountryOfficialName: officialName,
		SkinAreaCoveragePPM: funFact.SkinCoverageOfCountry,
		PopulationSharePPM:  funFact.ShareOfPopulation,
	}
	fmt.Println(playerTrivia)
	return playerTrivia, nil

}

// getOfficialName makes a api call to restcountries-api and get the official name
// of the country. Ex. Kingdom of Sweden.
func getOfficialName(countryName string, out chan<- string) {
	restcountriesService := NewRestcountriesService()
	officialName, err := restcountriesService.GetOfficialName(countryName)
	if err != nil {
		panic("ERROR")
	}
	out <- officialName
}

func calculateFunFacts(countryName string, height, weight float64, out chan<- entity.CalculatedFunFact) {
	// First we need to get the area and population of the country, this info we get from ApiNinja
	countryApiNinjaService := NewCountryApiNinja()
	countryApiNinjaInfo, err := countryApiNinjaService.GetCountryInfo(countryName)
	if err != nil {
		panic("ERROR")
	}

	var calculatedFunFact entity.CalculatedFunFact
	// Calculating share of population
	calculatedFunFact.ShareOfPopulation = calculateShareOfPopulation(countryApiNinjaInfo.Population)

	// Calculating area of coverage
	calculatedFunFact.SkinCoverageOfCountry = calculateSkinAreaCoverage(height, weight, countryApiNinjaInfo.SurfaceArea)

	// Sends back facts into channel.
	out <- calculatedFunFact
}

// calculatingSkinAreaCoverage returns the PPM of area coverage.
func calculateSkinAreaCoverage(height, weight float64, countryArea float64) float64 {
	// Calucating skin area coverage of the country. This is done by using Davis, F.A. 1993 method on how to calculate
	// the area of the skin of a human body.
	//
	// "The surface area may be calculated by multiplying 0.007184 times the weight in kilograms raised to the 0.425 power
	// and the height in centimeters raised to the 0.725 power."
	//

	skinArea := math.Pow(0.007184*weight, .425) + math.Pow(height, 0.725)

	// Converting countryArea from km2 to m2
	countryArea = countryArea * 1000000

	// Calculating the share in PPM
	skinAreaCoverage := (skinArea / countryArea) * 1000000

	// Rounding to 5 decimals.
	skinAreaCoverage = toFixed(skinAreaCoverage, 5)
	return skinAreaCoverage

}

// calculateShareOfPopulation calculates the share of population 1 person is,
// the returning float is "ppm", so 1 / 1 000 000 with 5 percision.
func calculateShareOfPopulation(population float64) float64 {
	// Convert population from million.
	population = population / .001

	// Calculate share in PPM.
	shareOfPopulation := (1 / population) * 1000000

	// Rounding to fixed percision.
	shareOfPopulation = toFixed(shareOfPopulation, 5)
	return shareOfPopulation
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
