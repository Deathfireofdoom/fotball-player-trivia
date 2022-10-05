package entity

type CountryInfo struct {
	Name struct {
		OfficialName string `json:"official"`
	} `json:"name"`
}

type RestcountriesResponse []CountryInfo
