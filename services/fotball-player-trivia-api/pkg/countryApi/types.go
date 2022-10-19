package countryApi

type CountryInfo struct {
	Population  float64 `json:"population"`
	SurfaceArea float64 `json:"surface_area"`
}

type ResponseApiNinja []CountryInfo
