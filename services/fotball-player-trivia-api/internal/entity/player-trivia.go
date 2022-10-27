package entity

type PlayerTrivia struct {
	Name                string  `json:"name"`
	Country             string  `json:"country"`
	CountryOfficialName string  `json:"country_official_name"`
	SkinAreaCoveragePPM float64 `json:"skin_area_coverage_ppm"`
	PopulationSharePPM  float64 `json:"population_share_ppm"`
}
