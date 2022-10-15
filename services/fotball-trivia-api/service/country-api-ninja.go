package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
)

var (
	apiKey string = "Xbjl4nonbcC1Cnbzurq2aA==w7iNEXQkbHlEhz7P"
)

type CountryApiNinjaService interface {
	GetCountryInfo(string) (entity.CountryInfoApiNinja, error)
}

type countryApiNinjaService struct{}

func NewCountryApiNinja() CountryApiNinjaService {
	return &countryApiNinjaService{}
}

func (cs *countryApiNinjaService) GetCountryInfo(country string) (entity.CountryInfoApiNinja, error) {
	URL := fmt.Sprintf("https://api.api-ninjas.com/v1/country?name=%s", country)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header = http.Header{"X-Api-Key": {apiKey}}
	resp, err := client.Do(req)
	if err != nil {
		panic("Could not get answer from ninja api.")
	}
	defer resp.Body.Close()

	var response entity.ResponseApiNinja

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		panic("Could not decode response from ninja api.")
	}
	fmt.Println(response)
	return response[0], nil
}
