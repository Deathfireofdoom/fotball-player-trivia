package trivia

import (
	"fmt"
	"math"

	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/db"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/entity"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/internal/redis"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/pkg/countryApi"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/pkg/restCountries"
	"github.com/Deathfireofdoom/fotball-player-trivia-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GetPlayerTriva(playerName string, ctx *gin.Context) (entity.PlayerTrivia, error) {
	// Checks if playerName's fun fact is already in redis-cache.
	// GetTrivia will return err = nil if playerTrivia was found, otherwise it will return error.
	var err error
	var playerTrivia entity.PlayerTrivia

	// Checks if redis client is initialized, if not, this mean that redis was not reachable.
	// Still works but without the caching mechanism.
	if redis.Client != nil {
		playerTrivia, err := redis.Client.GetTrivia(ctx, playerName)

		// Hit in cache
		if err == nil {
			fmt.Println("HITING CACHE")
			return playerTrivia, nil
		}
	}

	// No hit in cache
	playerTrivia, err = generatePlayerTrivia(playerName)
	if err != nil {
		fmt.Println("ERROR - generatePlayerTrivia")
		fmt.Println(err)
		return entity.PlayerTrivia{}, fmt.Errorf("could not generate trivia: %w", err)
	}

	if redis.Client != nil {
		// Saves generated playerTrivia in cache for next time.
		err = redis.Client.SaveTrivia(ctx, playerTrivia, 25)
		if err != nil {
			return playerTrivia, fmt.Errorf("could not save playerInfo into cache: %w", err)
		}
	}
	fmt.Println("About to return")
	return playerTrivia, nil
}

func generatePlayerTrivia(playerName string) (entity.PlayerTrivia, error) {
	// Fetch country, height and weight from DB.
	playerInfo, err := db.GetPlayerInfo(playerName)
	if err != nil {
		return entity.PlayerTrivia{}, fmt.Errorf("could not get player-info from DB: %w", &err)
	}

	// Makes channel to get results from go-routines.
	chanFunFacts := make(chan CalculatedFunFact)
	chanOfficialName := make(chan string)

	// Starts go routines.
	go getOfficialName(playerInfo.Country, chanOfficialName)
	go calculateFunFacts(playerInfo.Country, playerInfo.Height, playerInfo.Weight, chanFunFacts)

	// Waits for routines to finnish.
	officialName := <-chanOfficialName
	funFact := <-chanFunFacts

	return entity.PlayerTrivia{
		Name:                playerName,
		Country:             playerInfo.Country,
		CountryOfficialName: officialName,
		SkinAreaCoveragePPM: funFact.SkinCoverageOfCountry,
		PopulationSharePPM:  funFact.ShareOfPopulation,
	}, nil

}

// getOfficialName makes a api call to restcountries-api and get the official name
// of the country. Ex. Kingdom of Sweden.
func getOfficialName(countryName string, out chan<- string) {
	officalName, err := restCountries.GetOfficalName(countryName)
	if err != nil {
		panic("Implement this.")
	}
	out <- officalName
}

func calculateFunFacts(countryName string, height, weight float64, out chan<- CalculatedFunFact) {
	// First we need to get the area and population of the country, this info we get from ApiNinja

	countryInfo, err := countryApi.GetCountryInfo(countryName)
	if err != nil {
		panic(err)
	}

	var calculatedFunFact CalculatedFunFact
	// Calculating share of population
	calculatedFunFact.ShareOfPopulation = calculateShareOfPopulation(countryInfo.Population)

	// Calculating area of coverage
	calculatedFunFact.SkinCoverageOfCountry = calculateSkinAreaCoverage(height, weight, countryInfo.SurfaceArea)

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
	skinAreaCoverage = utils.ToFixed(skinAreaCoverage, 5)
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
	shareOfPopulation = utils.ToFixed(shareOfPopulation, 5)
	return shareOfPopulation
}
