package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Deathfireofdoom/fotball-player-trivia/entity"
)

type RestcountriesService interface {
	GetOfficialName(string) (string, error)
}

type restcountriesService struct{}

func NewRestcountriesService() RestcountriesService {
	return &restcountriesService{}
}

func (rs *restcountriesService) GetOfficialName(country string) (string, error) {
	url := fmt.Sprintf("https://restcountries.com/v3.1/name/%s?fullText=true", country)
	resp, err := http.Get(url)
	if err != nil {
		panic("Could not get get answer from restcountries api.")
	}
	defer resp.Body.Close()

	var response entity.RestcountriesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic("Could not decode response from restcountries api.")
	}

	return response[0].Name.OfficialName, nil

}
