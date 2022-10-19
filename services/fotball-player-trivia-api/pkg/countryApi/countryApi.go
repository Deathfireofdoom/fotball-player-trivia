package countryApi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	apiKey string = "Xbjl4nonbcC1Cnbzurq2aA==w7iNEXQkbHlEhz7P"
)

func GetCountryInfo(countryName string) (CountryInfo, error) {
	// Setting up request.
	URL := fmt.Sprintf("https://api-ninjas.com/v1/country?name=%s", countryName)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header = http.Header{"X-Api-key": {apiKey}}

	resp, err := client.Do(req)
	if err != nil {
		return CountryInfo{}, fmt.Errorf("Could not get answer from Ninja api: %w", err)
	}
	defer resp.Body.Close()

	var response ResponseApiNinja
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return CountryInfo{}, fmt.Errorf("Could not parse resposne from Ninja api: %w", err)
	}

	return response[0], nil
}
