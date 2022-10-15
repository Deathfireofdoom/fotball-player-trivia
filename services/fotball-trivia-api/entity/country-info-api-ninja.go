package entity

type CountryInfoApiNinja struct {
	Population  float64 `json:"population"`
	SurfaceArea float64 `json:"surface_area"`
}

type ResponseApiNinja []CountryInfoApiNinja
