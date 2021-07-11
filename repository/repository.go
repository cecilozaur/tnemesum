package repository

import "github.com/cecilozaur/tnemesum/domain"

// Repository defines the interface for a repository
type Repository interface {
	// GetItems will return all items inside the repository
	GetItems() []domain.City
	// Get return an associated item, if one is found, otherwise an error is returned
	Get(key uint64) (domain.City, error)
	// GetForecast will return the forecast for the specified city or an error if the info is not found
	GetForecast(key uint64) (domain.Forecast, error)
	// Store saves the item with key inside the repository
	Store(key uint64, item domain.City)
	// UpdateForecast updates the forecast for the city specified by key
	UpdateForecast(key uint64, forecast domain.Forecast) bool
}
