package restCountries

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetOfficalName(countryName string) (string, error) {
	URL := fmt.Sprintf("https://restcountries.com/v3.1/name/%s?fullText=true", countryName)
	resp, err := http.Get(URL)

	if err != nil {
		return "", fmt.Errorf("Could not get answer from Restcountries-api: %w", err)
	}
	defer resp.Body.Close()

	var response RestcountriesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("Could not parse response from Restcountries-api: %w", err)
	}

	return response[0].Name.OfficialName, nil
}
