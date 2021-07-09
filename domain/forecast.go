package domain

type Forecast []ForecastDay

type ForecastDay struct {
	Day       string `json:"day"`
	Condition string `json:"condition"`
}
