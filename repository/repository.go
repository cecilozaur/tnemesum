package repository

import "github.com/cecilozaur/tnemesum/domain"

type Repository interface {
	GetItems() []domain.City
	Get(key uint64) (domain.City, error)
	Store(key uint64, item domain.City)
	UpdateForecast(key uint64, forecast domain.Forecast) bool
	GetForecast(key uint64, day int) domain.Forecast
}
