package domain

type City struct {
	Id       uint64  `json:"id"`
	UUID     string  `json:"uuid"`
	Name     string  `json:"name"`
	Lat      float64 `json:"latitude"`
	Lng      float64 `json:"longitude"`
	Forecast Forecast
}
