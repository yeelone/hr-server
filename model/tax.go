package model

type TaxConf struct {
	Threshold          float64   `json:"threshold"`
	OldAgeRating       float64   `json:"old_age_rating"`
	MedicalRating      float64   `json:"medical_rating"`
	UnemploymentRating float64   `json:"unemployment_rating"`
	HousingFundRating  float64   `json:"housing_fund_rating"`
	Level              []float64 `json:"level"`
	Rating             []float64 `json:"rating"`
	Deduction          []float64 `json:"deduction"`
}
