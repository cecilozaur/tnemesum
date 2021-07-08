package domain

type ForecastResult struct {
	Forecast Forecast `json:"forecast"`
}

type Forecast struct {
	ForecastDay []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date      string `json:"date"`
	DateEpoch int    `json:"date_epoch"`
	Day       Day    `json:"day"`
}

type Day struct {
	Condition Condition `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
	Code int    `json:"code"`
}
